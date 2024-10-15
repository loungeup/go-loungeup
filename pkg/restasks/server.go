package restasks

import (
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/pkg/cache"
	"github.com/loungeup/go-loungeup/pkg/errors"
	"github.com/loungeup/go-loungeup/pkg/log"
)

type Server struct {
	cache   cache.ReadWriter
	service *res.Service
}

// NewServer creates a new server.
func NewServer(cache cache.ReadWriter, service *res.Service) *Server {
	result := &Server{cache, service}
	result.addRESHandlers()

	return result
}

// CreateTask and returns its RID.
func (s *Server) CreateTask() (string, error) {
	task := newTask(s.service.FullPath())
	cacheTask(s.cache, task)

	return task.rid(), nil
}

// CompleteTask with the given result.
func (s *Server) CompleteTask(rid string, result any) error {
	task, err := readCachedTask(s.cache, rid)
	if err != nil {
		return err
	}

	task.Progress = taskMaxProgress
	task.Error = nil
	task.Result = result

	return s.sendTaskChangeEvent(task)
}

// FailTask with the given error.
func (s *Server) FailTask(rid string, err error) error {
	task, readError := readCachedTask(s.cache, rid)
	if readError != nil {
		return readError
	}

	task.Progress = taskMinProgress
	task.Error = err
	task.Result = nil

	return s.sendTaskChangeEvent(task)
}

func (s *Server) SetTaskProgress(rid string, progress int) error {
	task, err := readCachedTask(s.cache, rid)
	if err != nil {
		return err
	}

	if err := task.setProgress(progress); err != nil {
		return err
	}

	return s.sendTaskChangeEvent(task)
}

// addRESHandlers to the server.
func (s *Server) addRESHandlers() {
	s.service.Handle("tasks.$taskID", res.GetModel(func(request res.ModelRequest) {
		task, err := readCachedTask(s.cache, request.ResourceName())
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
		resource.ChangeEvent(task.toChangeEventProperties())
	})
}

func cacheTask(cache cache.Writer, task *task) { cache.Write(task.rid(), task) }

func readCachedTask(cache cache.Reader, rid string) (*task, error) {
	if task, ok := cache.Read(rid).(*task); ok {
		return task, nil
	}

	return nil, &errors.Error{Code: errors.CodeNotFound, Message: "Task not found"}
}
