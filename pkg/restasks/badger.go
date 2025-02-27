package restasks

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/errors"
	"github.com/loungeup/go-loungeup/pkg/log"
)

type badgerStore struct {
	db        *badger.DB
	logger    *log.Logger
	retention time.Duration
}

type badgerStoreOption func(*badgerStore)

func NewBadgerStore(db *badger.DB, options ...badgerStoreOption) *badgerStore {
	const defaultRetention = time.Hour * 24 * 7

	result := &badgerStore{
		db:        db,
		logger:    log.Default(),
		retention: defaultRetention,
	}
	for _, option := range options {
		option(result)
	}

	go result.runGC()

	return result
}

func WithBadgerStoreLogger(logger *log.Logger) badgerStoreOption {
	return func(s *badgerStore) { s.logger = logger }
}

func WithBadgerStoreRetention(retention time.Duration) badgerStoreOption {
	return func(s *badgerStore) { s.retention = retention }
}

var _ (Store) = (*badgerStore)(nil)

func (s *badgerStore) ReadByID(id uuid.UUID) (*Task, error) {
	model := &badgerTaskModel{ID: id}

	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(model.key())
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return &errors.Error{
					Code:            errors.CodeNotFound,
					Message:         "Task not found",
					UnderlyingError: err,
				}
			} else {
				return err
			}
		}

		if err := item.Value(func(encodedModel []byte) error { return json.Unmarshal(encodedModel, model) }); err != nil {
			return fmt.Errorf("could not decode Badger task model: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapBadgerModelToTask(model), nil
}

func (s *badgerStore) Write(task *Task) error {
	model := mapTaskToBadgerModel(task)

	encodedModel, err := model.encode()
	if err != nil {
		return fmt.Errorf("could not encode Badger task model: %w", err)
	}

	if err := s.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(badger.NewEntry(model.key(), encodedModel).WithTTL(s.retention))
	}); err != nil {
		return fmt.Errorf("could not write task model to Badger DB: %w", err)
	}

	return nil
}

// https://dgraph.io/docs/badger/get-started/#garbage-collection
func (s *badgerStore) runGC() {
	const (
		discardRatio = 0.5
		runInterval  = 5 * time.Minute
	)

	ticker := time.NewTicker(runInterval)
	defer ticker.Stop()

	for range ticker.C {
	again:
		l1 := s.logger.With(slog.String("traceId", uuid.NewString()))
		l1.Debug("Running Badger GC")

		if err := s.db.RunValueLogGC(discardRatio); err == nil {
			// One call would only result in removal of at max one log file. As an optimization, immediately re-run it
			// whenever it returns nil error (indicating a successful value log GC)
			goto again
		} else if !errors.Is(err, badger.ErrNoRewrite) {
			l1.Error("Could not run Badger GC", slog.Any("error", err))
		}
	}
}

type badgerTaskModel struct {
	ID           uuid.UUID       `json:"id"`
	Progress     int             `json:"progress"`
	ErrorMessage string          `json:"errorMessage"`
	Result       json.RawMessage `json:"result"`
	StartedAt    string          `json:"startedAt"`
	EndedAt      string          `json:"endedAt"`
}

func (m *badgerTaskModel) encode() ([]byte, error) { return json.Marshal(m) }

func (m *badgerTaskModel) key() []byte { return []byte(m.ID.String()) }

func mapBadgerModelToTask(model *badgerTaskModel) *Task {
	return &Task{
		ID:       model.ID,
		Progress: model.Progress,
		Error: func() error {
			if message := model.ErrorMessage; message != "" {
				return fmt.Errorf("%s", message)
			}

			return nil
		}(),
		Result: func() any {
			var result any
			_ = json.Unmarshal(model.Result, &result)

			return result
		}(),
		StartedAt: func() time.Time {
			result, _ := parseTime(model.StartedAt)

			return result
		}(),
		EndedAt: func() time.Time {
			result, _ := parseTime(model.EndedAt)

			return result
		}(),
	}
}

func mapTaskToBadgerModel(task *Task) *badgerTaskModel {
	return &badgerTaskModel{
		ID:       task.ID,
		Progress: task.Progress,
		ErrorMessage: func() string {
			if task.Error != nil {
				return task.Error.Error()
			}

			return ""
		}(),
		Result: func() json.RawMessage {
			result, _ := json.Marshal(task.Result)

			return result
		}(),
		StartedAt: formatTime(task.StartedAt),
		EndedAt:   formatTime(task.EndedAt),
	}
}
