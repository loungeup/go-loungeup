package restasks

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/pkg/errors"
)

// TODO: Set progress and status from sub-tasks.

const (
	taskStatusCompleted = "completed"
	taskStatusFailed    = "failed"
	taskStatusStarted   = "started"

	taskMinProgress = 0
	taskMaxProgress = 100
)

type task struct {
	serviceName string
	id          uuid.UUID
	progress    int
	err         error
	result      any
	subTaskRIDs []string
}

func (t *task) addSubTaskRID(rid string) { t.subTaskRIDs = append(t.subTaskRIDs, rid) }

func (t *task) rid() string         { return t.serviceName + ".tasks." + t.id.String() }
func (t *task) subTasksRID() string { return t.rid() + ".sub-tasks" }

// setProgress of the task. It must be between taskMinProgress and taskMaxProgress.
func (t *task) setProgress(progress int) error {
	if progress == t.progress {
		return nil
	}

	if progress < taskMinProgress || progress > taskMaxProgress {
		return &errors.Error{
			Code:    errors.CodeInvalid,
			Message: fmt.Sprintf("Progress must be between %d and %d", taskMinProgress, taskMaxProgress),
		}
	}

	t.progress = progress

	return nil
}

// status of the task.
func (t *task) status() string {
	switch {
	case t.err != nil:
		return taskStatusFailed
	case t.progress == taskMaxProgress:
		return taskStatusCompleted
	default:
		return taskStatusStarted
	}
}

func (t *task) toChangeEventProperties() map[string]any {
	result := map[string]any{
		"progress": t.progress,
		"status":   t.status(),
	}

	if t.err != nil {
		result["error"] = errors.ErrorMessage(t.err)
	}

	if t.result != nil {
		result["result"] = &res.DataValue[any]{Data: t.result}
	}

	if len(t.subTaskRIDs) > 0 {
		result["subTasks"] = res.SoftRef(t.subTasksRID())
	}

	return result
}

func (t *task) toModel() *taskModel {
	result := &taskModel{
		Progress: t.progress,
		Status:   t.status(),
	}

	if t.err != nil {
		result.Error = errors.ErrorMessage(t.err)
	}

	if t.result != nil {
		result.Result = &res.DataValue[any]{Data: t.result}
	}

	if len(t.subTaskRIDs) > 0 {
		result.SubTasks = res.SoftRef(t.subTasksRID())
	}

	return result
}

type taskModel struct {
	Progress int                 `json:"progress"`
	Status   string              `json:"status"`
	Error    string              `json:"error,omitempty"`
	Result   *res.DataValue[any] `json:"result,omitempty"`
	SubTasks res.SoftRef         `json:"subTasks,omitempty"`
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
