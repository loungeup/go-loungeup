package restasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/errors"
	"github.com/nats-io/nats.go/jetstream"
)

type jetStreamKeyValueStore struct{ store jetstream.KeyValue }

func NewJetStreamKeyValueStore(store jetstream.KeyValue) *jetStreamKeyValueStore {
	return &jetStreamKeyValueStore{store}
}

var _ (Store) = (*jetStreamKeyValueStore)(nil)

func (s *jetStreamKeyValueStore) ReadByID(id uuid.UUID) (*Task, error) {
	entry, err := s.store.Get(context.Background(), id.String())
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, &errors.Error{Code: errors.CodeNotFound}
		} else {
			return nil, err
		}
	}

	model := &jetStreamTaskModel{}
	if err := json.Unmarshal(entry.Value(), model); err != nil {
		return nil, fmt.Errorf("could not decode JetStream task model: %w", err)
	}

	return mapJetStreamModelToTask(model), nil
}

func (s *jetStreamKeyValueStore) Write(task *Task) error {
	encodedModel, err := json.Marshal(mapTaskToJetStreamModel(task))
	if err != nil {
		return fmt.Errorf("could not encode JetStream task model: %w", err)
	}

	if _, err := s.store.Put(context.Background(), task.ID.String(), encodedModel); err != nil {
		return fmt.Errorf("could not write task model to JetStream: %w", err)
	}

	return nil
}

type jetStreamTaskModel struct {
	ID           uuid.UUID       `json:"id"`
	Progress     int             `json:"progress"`
	ErrorMessage string          `json:"errorMessage"`
	Result       json.RawMessage `json:"result"`
	StartedAt    string          `json:"startedAt"`
	EndedAt      string          `json:"endedAt"`
}

func mapJetStreamModelToTask(model *jetStreamTaskModel) *Task {
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

func mapTaskToJetStreamModel(task *Task) *jetStreamTaskModel {
	return &jetStreamTaskModel{
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
