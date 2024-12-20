package restasks

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/errors"
)

var (
	//go:embed data/create_tasks_table.sql
	createTasksTableSQLQuery string

	//go:embed data/read_task_by_id.sql
	readTaskByIDSQLQuery string

	//go:embed data/write_task.sql
	writeTaskSQLQuery string
)

type sqlStore struct{ db *sql.DB }

func NewSQLStore(db *sql.DB) (*sqlStore, error) {
	if _, err := db.Exec(createTasksTableSQLQuery); err != nil {
		return nil, fmt.Errorf("could not create 'tasks' SQL table: %w", err)
	}

	return &sqlStore{db}, nil
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
			return nil, fmt.Errorf("could not read task from SQL DB: %w", err)
		}
	}

	result, err := mapSQLModelToTask(model)
	if err != nil {
		return nil, fmt.Errorf("could not map SQL model to task: %w", err)
	}

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
		return fmt.Errorf("could not write task to SQL DB: %w", err)
	}

	return nil
}

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
				return fmt.Errorf(model.ErrorMessage.String)
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
