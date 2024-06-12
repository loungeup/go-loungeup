package restasks

import (
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/pkg/errors"
	"github.com/loungeup/go-loungeup/pkg/log"
)

type Server struct {
	service *res.Service
	tasks   TaskReadWriter
}

// Option is a function type that can be used to configure a [Server].
type Option func(*Server)

// NewServer creates a new server with the given service and [TaskReadWriter].
func NewServer(service *res.Service, tasks TaskReadWriter) *Server {
	result := &Server{service, tasks}
	result.addRESHandlers()

	return result
}

// CreateTask and returns its RID.
func (s *Server) CreateTask() (string, error) {
	task := newTask(s.service.FullPath())
	if err := s.tasks.WriteTask(task); err != nil {
		return "", err
	}

	return task.rid(), nil
}

// CompleteTask with the given result.
func (s *Server) CompleteTask(rid string, result any) error {
	task, err := s.tasks.ReadTask(rid)
	if err != nil {
		return err
	}

	task.Error = nil
	task.Result = result

	return s.sendTaskChangeEvent(task)
}

// FailTask with the given error.
func (s *Server) FailTask(rid string, err error) error {
	task, readError := s.tasks.ReadTask(rid)
	if readError != nil {
		return readError
	}

	task.Error = err
	task.Result = nil

	return s.sendTaskChangeEvent(task)
}

// addRESHandlers to the server.
func (s *Server) addRESHandlers() {
	s.service.Handle("tasks.$taskID", res.GetModel(func(request res.ModelRequest) {
		task, err := s.tasks.ReadTask(request.ResourceName())
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)
			return
		}

		request.Model(task.toModel())
	}))
}

// sendTaskChangeEvent for the given task.
func (s *Server) sendTaskChangeEvent(task *task) error {
	return s.service.With(task.rid(), func(resource res.Resource) {
		resource.ChangeEvent(task.toModel())
	})
}
