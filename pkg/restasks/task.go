package restasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/pkg/errors"
)

const (
	taskStatusCompleted = "completed"
	taskStatusFailed    = "failed"
	taskStatusStarted   = "started"

	taskMinProgress = 0
	taskMaxProgress = 100
)

type task struct {
	ServiceName string
	ID          uuid.UUID
	Progress    int
	Error       error
	Result      any
	CreatedAt   time.Time
}

// newTask creates a new task for the service with the given name.
func newTask(serviceName string) *task {
	return &task{
		ServiceName: serviceName,
		ID:          uuid.New(),
		Progress:    taskMinProgress,
		CreatedAt:   time.Now(),
	}
}

// rid of the task for the given service.
func (t *task) rid() string {
	return t.ServiceName + ".tasks." + t.ID.String()
}

// setProgress of the task. It must be between taskMinProgress and taskMaxProgress.
func (t *task) setProgress(progress int) error {
	if progress == t.Progress {
		return nil
	}

	if progress < taskMinProgress || progress > taskMaxProgress {
		return &errors.Error{
			Code:    errors.CodeInvalid,
			Message: fmt.Sprintf("Progress must be between %d and %d", taskMinProgress, taskMaxProgress),
		}
	}

	t.Progress = progress

	return nil
}

// status of the task.
func (t *task) status() string {
	switch {
	case t.Error != nil:
		return taskStatusFailed
	case t.Progress == taskMaxProgress:
		return taskStatusCompleted
	default:
		return taskStatusStarted
	}
}

func (t *task) toChangeEventProperties() map[string]any {
	result := map[string]any{
		"progress": t.Progress,
		"status":   t.status(),
	}

	if t.Error != nil {
		result["error"] = errors.ErrorMessage(t.Error)
	}

	if t.Result != nil {
		result["result"] = &res.DataValue[any]{Data: t.Result}
	}

	return result
}

func (t *task) toModel() *taskModel {
	result := &taskModel{
		Progress: t.Progress,
		Status:   t.status(),
	}

	if t.Error != nil {
		result.Error = errors.ErrorMessage(t.Error)
	}

	if t.Result != nil {
		result.Result = &res.DataValue[any]{Data: t.Result}
	}

	return result
}

type taskModel struct {
	Progress int                 `json:"progress"`
	Status   string              `json:"status"`
	Error    string              `json:"error,omitempty"`
	Result   *res.DataValue[any] `json:"result,omitempty"`
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
