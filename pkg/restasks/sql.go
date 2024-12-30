package restasks

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/errors"
	"github.com/loungeup/go-loungeup/pkg/log"
)

var (
	//go:embed data/create_tasks_table.sql
	createTasksTableSQLQuery string

	//go:embed data/read_task_by_id.sql
	readTaskByIDSQLQuery string

	//go:embed data/write_task.sql
	writeTaskSQLQuery string
)

type sqlStore struct {
	db         *sql.DB
	logger     *log.Logger
	purgedAt   time.Time
	purgeMutex sync.Mutex
	retention  time.Duration
}

type sqlStoreOption func(*sqlStore)

func NewSQLStore(db *sql.DB, options ...sqlStoreOption) (*sqlStore, error) {
	const defaultRetention = 7 * 24 * time.Hour

	result := &sqlStore{
		db:        db,
		logger:    log.Default(),
		retention: defaultRetention,
	}
	for _, option := range options {
		option(result)
	}

	if err := result.init(); err != nil {
		return nil, fmt.Errorf("could not initialize SQL store: %w", err)
	}

	return result, nil
}

func WithSQLStoreLogger(logger *log.Logger) sqlStoreOption {
	return func(s *sqlStore) { s.logger = logger }
}

func WithSQLStoreRetention(retention time.Duration) sqlStoreOption {
	return func(s *sqlStore) { s.retention = retention }
}

var _ (Store) = (*sqlStore)(nil)

func (s *sqlStore) ReadByID(id uuid.UUID) (*Task, error) {
	model := &taskSQLModel{}
	if err := s.db.QueryRow(readTaskByIDSQLQuery, id).Scan(
		&model.ID,
		&model.ErrorMessage,
		&model.Result,
		&model.Progress,
		&model.StartedAt,
		&model.EndedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &errors.Error{Code: errors.CodeNotFound}
		} else {
			return nil, fmt.Errorf("could not read task from DB: %w", err)
		}
	}

	result, err := mapSQLModelToTask(model)
	if err != nil {
		return nil, fmt.Errorf("could not map SQL model to task: %w", err)
	}

	defer s.eventuallyPurge()

	return result, nil
}

func (s *sqlStore) Write(task *Task) error {
	model, err := mapTaskToSQLModel(task)
	if err != nil {
		return fmt.Errorf("could not map task to SQL model: %w", err)
	}

	if _, err := s.db.Exec(
		writeTaskSQLQuery,
		model.ID,
		model.ErrorMessage,
		model.Result,
		model.Progress,
		model.StartedAt,
		model.EndedAt,
	); err != nil {
		return fmt.Errorf("could not write task to DB: %w", err)
	}

	defer s.eventuallyPurge()

	return nil
}

func (s *sqlStore) eventuallyPurge() error {
	if !s.shouldPurge() {
		return nil
	}

	return s.purge()
}

func (s *sqlStore) purge() error {
	if !s.purgeMutex.TryLock() {
		return nil // A purge is already in progress.
	}
	defer s.purgeMutex.Unlock()

	startedAt := time.Now()

	l1 := s.logger.With(slog.String("traceID", uuid.NewString()))
	l1.Debug("Purging DB",
		slog.String("retention", s.retention.String()),
		slog.Time("purgedAt", s.purgedAt),
		slog.Time("startedAt", startedAt),
	)

	result, err := s.db.Exec(
		"DELETE FROM tasks WHERE ended_at < NOW() - INTERVAL ? SECOND",
		s.retention.Seconds(),
	)
	if err != nil {
		return fmt.Errorf("could not purge DB: %w", err)
	}

	totalDeletedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not count deleted rows: %w", err)
	}

	l1.Debug("DB purged",
		slog.Duration("duration", time.Since(startedAt)),
		slog.Int64("totalDeletedRows", totalDeletedRows),
	)

	s.purgedAt = time.Now()

	return nil
}

func (s *sqlStore) init() error {
	if _, err := s.db.Exec(createTasksTableSQLQuery); err != nil {
		return fmt.Errorf("could not create 'tasks' table: %w", err)
	}

	if _, err := s.db.Exec(
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS tasks_ended_at_idx ON tasks (ended_at)`,
	); err != nil {
		return fmt.Errorf("could not create 'tasks_ended_at_idx' index: %w", err)
	}

	return nil
}

func (s *sqlStore) shouldPurge() bool { return time.Since(s.purgedAt) > s.retention }

type taskSQLModel struct {
	ID           uuid.UUID
	Progress     int
	ErrorMessage sql.NullString
	Result       *json.RawMessage
	StartedAt    time.Time
	EndedAt      sql.NullTime
}

func (m *taskSQLModel) decodeResult() (any, error) {
	if m.Result == nil {
		return nil, nil
	}

	var result any
	if err := json.Unmarshal(*m.Result, &result); err != nil {
		return nil, fmt.Errorf("could not decode SQL model result: %w", err)
	}

	return result, nil
}

func mapSQLModelToTask(model *taskSQLModel) (*Task, error) {
	decodedResult, err := model.decodeResult()
	if err != nil {
		return nil, err
	}

	return &Task{
		ID:       model.ID,
		Progress: model.Progress,
		Error: func() error {
			if model.ErrorMessage.Valid {
				return fmt.Errorf("%s", model.ErrorMessage.String)
			}

			return nil
		}(),
		Result:    decodedResult,
		StartedAt: model.StartedAt,
		EndedAt:   model.EndedAt.Time,
	}, nil
}

func mapTaskToSQLModel(task *Task) (*taskSQLModel, error) {
	encodedResult, err := mapValueToNullJSON(task.Result)
	if err != nil {
		return nil, fmt.Errorf("could not map task result to SQL model: %w", err)
	}

	return &taskSQLModel{
		ID:       task.ID,
		Progress: task.Progress,
		ErrorMessage: func() sql.NullString {
			if err := task.Error; err != nil {
				return sql.NullString{String: err.Error(), Valid: true}
			}

			return sql.NullString{}
		}(),
		Result:    encodedResult,
		StartedAt: task.StartedAt,
		EndedAt: func() sql.NullTime {
			if endedAt := task.EndedAt; !endedAt.IsZero() {
				return sql.NullTime{Time: endedAt, Valid: true}
			}

			return sql.NullTime{}
		}(),
	}, nil
}

func mapValueToNullJSON(value any) (*json.RawMessage, error) {
	if value == nil {
		return nil, nil
	}

	encodedResult, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("could not encode value: %w", err)
	}

	return (*json.RawMessage)(&encodedResult), nil
}
