package restasks

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	taskMinProgress = 0
	taskMaxProgress = 100
)

type Task struct {
	ID                 uuid.UUID
	Progress           int
	Error              error
	Result             any
	StartedAt, EndedAt time.Time
}

func (t *Task) setProgress(progress int) error {
	if progress == t.Progress {
		return nil
	}

	if progress < taskMinProgress || progress > taskMaxProgress {
		return fmt.Errorf("progress must be between %d and %d", taskMinProgress, taskMaxProgress)
	}

	t.Progress = progress

	return nil
}

func (t *Task) setError(err error) {
	t.Progress = taskMaxProgress
	t.Error = err
	t.Result = nil
	t.EndedAt = time.Now()
}

func (t *Task) setResult(result any) {
	t.Progress = taskMaxProgress
	t.Error = nil
	t.Result = result
	t.EndedAt = time.Now()
}

func (t *Task) status() taskStatus {
	switch {
	case t.Error != nil:
		return taskStatusFailed
	case t.Result != nil:
		return taskStatusCompleted
	default:
		return taskStatusStarted
	}
}

type taskStatus string

const (
	taskStatusCompleted taskStatus = "completed"
	taskStatusFailed    taskStatus = "failed"
	taskStatusStarted   taskStatus = "started"
)

func (s taskStatus) String() string { return string(s) }
