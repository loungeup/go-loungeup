package restasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/pkg/errors"
)

type Server struct {
	service *res.Service
	store   Store
}

func NewServer(service *res.Service, store Store) *Server {
	result := &Server{service, store}
	result.addHandlers()

	return result
}

func (s *Server) CreateTask() (string, error) {
	newTask := &Task{
		ID:        uuid.New(),
		Progress:  taskMinProgress,
		StartedAt: time.Now(),
	}
	if err := s.store.Write(newTask); err != nil {
		return "", fmt.Errorf("could not write task: %w", err)
	}

	return s.makeTaskRID(newTask), nil
}

func (s *Server) CompleteTask(rid string, result any) error {
	return s.readAndWriteTaskFromRID(rid, func(task *Task) error {
		task.setResult(result)

		return nil
	})
}

func (s *Server) FailTask(rid string, err error) error {
	return s.readAndWriteTaskFromRID(rid, func(task *Task) error {
		task.setError(err)

		return nil
	})
}

func (s *Server) SetTaskProgress(rid string, progress int) error {
	return s.readAndWriteTaskFromRID(rid, func(task *Task) error {
		return task.setProgress(progress)
	})
}

func (s *Server) addHandlers() {
	s.service.Handle("tasks.$taskID", res.GetModel(func(request res.ModelRequest) {
		id, err := uuid.Parse(request.PathParam("taskID"))
		if err != nil {
			request.Error(&res.Error{
				Code:    res.CodeInvalidParams,
				Message: "Invalid task ID",
				Data:    err.Error(),
			})
		}

		task, err := s.store.ReadByID(id)
		if err != nil {
			if errors.ErrorCode(err) == errors.CodeNotFound {
				request.NotFound()
			} else {
				request.Error(&res.Error{
					Code:    res.CodeInternalError,
					Message: "Could not read task",
					Data: map[string]string{
						"errorMessage": err.Error(),
						"id":           id.String(),
					},
				})
			}

			return
		}

		request.Model(mapTaskToRESModel(task))
	}))
}

func (s *Server) makeTaskRID(task *Task) string {
	return s.service.FullPath() + ".tasks." + task.ID.String()
}

func (s *Server) parseTaskIDFromRID(rid string) (uuid.UUID, error) {
	values, _ := res.Pattern(s.service.FullPath() + ".tasks.$taskID").Values(rid)

	rawID, ok := values["taskID"]
	if !ok {
		return uuid.Nil, fmt.Errorf("could not extract task ID from RID")
	}

	result, err := uuid.Parse(rawID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not parse task ID: %w", err)
	}

	return result, nil
}

func (s *Server) readAndWriteTaskFromRID(rid string, modifyTaskFunc func(task *Task) error) error {
	id, err := s.parseTaskIDFromRID(rid)
	if err != nil {
		return err
	}

	task, err := s.store.ReadByID(id)
	if err != nil {
		return fmt.Errorf("could not read task by ID: %w", err)
	}

	if err := modifyTaskFunc(task); err != nil {
		return err
	}

	if err := s.store.Write(task); err != nil {
		return fmt.Errorf("could not write task: %w", err)
	}

	if err := s.service.With(s.makeTaskRID(task), func(resource res.Resource) {
		resource.ChangeEvent(mapTaskToRESChangeEventProperties(task))
	}); err != nil {
		return fmt.Errorf("could not send task change event: %w", err)
	}

	return nil
}

type taskRESModel struct {
	Progress  int                 `json:"progress"`
	Status    string              `json:"status"`
	Error     string              `json:"error,omitempty"`
	Result    *res.DataValue[any] `json:"result,omitempty"`
	StartedAt string              `json:"startedAt"`
	EndedAt   string              `json:"endedAt,omitempty"`
}

func (m *taskRESModel) decodeResult(value any) error {
	if m.Result == nil {
		return nil
	}

	encodedResult, err := json.Marshal(m.Result.Data)
	if err != nil {
		return fmt.Errorf("could not encode task RES model result: %w", err)
	}

	if err := json.Unmarshal(encodedResult, value); err != nil {
		return fmt.Errorf("could not decode task RES model result: %w", err)
	}

	return nil
}

func (m *taskRESModel) isRunning() bool { return m.Progress < taskMaxProgress }

func mapTaskToRESModel(task *Task) *taskRESModel {
	result := &taskRESModel{
		Progress:  task.Progress,
		Status:    task.status().String(),
		StartedAt: formatTime(task.StartedAt),
	}

	if err := task.Error; err != nil {
		result.Error = err.Error()
	}

	if task.Result != nil {
		result.Result = &res.DataValue[any]{Data: task.Result}
	}

	if endedAt := task.EndedAt; !endedAt.IsZero() {
		result.EndedAt = formatTime(endedAt)
	}

	return result
}

func mapTaskToRESChangeEventProperties(task *Task) map[string]any {
	result := map[string]any{
		"progress":  task.Progress,
		"status":    task.status(),
		"startedAt": formatTime(task.StartedAt),
	}

	if err := task.Error; err != nil {
		result["error"] = err.Error()
	}

	if task.Result != nil {
		result["result"] = &res.DataValue[any]{Data: task.Result}
	}

	if endedAt := task.EndedAt; !endedAt.IsZero() {
		result["endedAt"] = formatTime(endedAt)
	}

	return result
}

func formatTime(t time.Time) string { return t.Format(time.RFC3339) }
