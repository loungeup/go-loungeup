package restasks

import (
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

// toModel converts the task to a RES model.
func (t *task) toModel() map[string]any {
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
