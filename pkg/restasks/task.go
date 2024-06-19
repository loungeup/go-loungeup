package restasks

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/pkg/errors"
)

const (
	taskStatusCompleted = "completed"
	taskStatusFailed    = "failed"
	taskStatusStarted   = "started"
)

type task struct {
	ServiceName string
	ID          uuid.UUID
	Error       error
	Result      any
	CreatedAt   time.Time
}

// newTask creates a new task for the service with the given name.
func newTask(serviceName string) *task {
	return &task{
		ServiceName: serviceName,
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
	}
}

func (t *task) toChangeEventProperties() map[string]any {
	result := map[string]any{
		"status": t.status(),
	}

	if t.Error != nil {
		result["error"] = errors.ErrorMessage(t.Error)
	}

	if t.Result != nil {
		result["result"] = &res.DataValue{Data: t.Result}
	}

	return result
}

func (t *task) toModel() *taskModel {
	result := &taskModel{
		Status: t.status(),
	}

	if t.Error != nil {
		result.Error = errors.ErrorMessage(t.Error)
	}

	if t.Result != nil {
		result.Result = &res.DataValue{Data: t.Result}
	}

	return result
}

// rid of the task for the given service.
func (t *task) rid() string {
	return t.ServiceName + ".tasks." + t.ID.String()
}

// status of the task.
func (t *task) status() string {
	switch {
	case t.Error != nil:
		return taskStatusFailed
	case t.Result != nil:
		return taskStatusCompleted
	default:
		return taskStatusStarted
	}
}

type taskModel struct {
	Status string         `json:"status"`
	Error  string         `json:"error,omitempty"`
	Result *res.DataValue `json:"result,omitempty"`
}

func (m *taskModel) decodeResult(value any) error {
	if m.Result == nil {
		return nil
	}

	encodedResult, err := json.Marshal(m.Result.Data)
	if err != nil {
		return err
	}

	return json.Unmarshal(encodedResult, value)
}

func (m *taskModel) err() error {
	if m.Error == "" {
		return nil
	}

	return &errors.Error{Code: errors.CodeInternal, Message: m.Error}
}

func (m *taskModel) isRunning() bool { return m.Status == taskStatusStarted }
