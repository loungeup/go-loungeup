package restasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
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
func (s *Server) CreateTask() string {
	task := &task{
		serviceName: s.service.FullPath(),
		id:          uuid.New(),
		progress:    taskMinProgress,
	}
	cacheTask(s.cache, task)

	return task.rid()
}

func (s *Server) CreateSubTask(parentRID string) (string, error) {
	subTaskRID := s.CreateTask()

	if response := resprot.SendRequest(s.service.Conn(), "call."+parentRID+".sub-tasks.new", &resprot.Request{
		Params: res.Ref(subTaskRID),
	}, time.Second); response.HasError() {
		return "", fmt.Errorf("could not create sub-task: %w", response.Error)
	}

	return subTaskRID, nil
}

// CompleteTask with the given result.
func (s *Server) CompleteTask(rid string, result any) error {
	task, err := readCachedTask(s.cache, rid)
	if err != nil {
		return err
	}

	task.progress = taskMaxProgress
	task.err = nil
	task.result = result

	return s.sendTaskChangeEvent(task)
}

// FailTask with the given error.
func (s *Server) FailTask(rid string, err error) error {
	task, readError := readCachedTask(s.cache, rid)
	if readError != nil {
		return readError
	}

	task.progress = taskMinProgress
	task.err = err
	task.result = nil

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

	s.service.Handle("tasks.$taskID.sub-tasks",
		res.Call("new", func(request res.CallRequest) {
			parent, err := readCachedTask(s.cache, strings.TrimSuffix(request.ResourceName(), ".sub-tasks"))
			if err != nil {
				errors.LogAndWriteRESError(log.Default(), request, err)
				return
			}

			subTaskRID := res.Ref("")
			request.ParseParams(&subTaskRID)

			parent.addSubTaskRID(string(subTaskRID))

			request.Resource(string(subTaskRID))
			request.Service().Reset([]string{request.ResourceName()}, nil)
		}),
		res.GetCollection(func(request res.CollectionRequest) {
			task, err := readCachedTask(s.cache, strings.TrimSuffix(request.ResourceName(), ".sub-tasks"))
			if err != nil {
				errors.LogAndWriteRESError(log.Default(), request, err)
				return
			}

			collection := []res.Ref{}
			for _, subTaskRID := range task.subTaskRIDs {
				collection = append(collection, res.Ref(subTaskRID))
			}

			request.Collection(collection)
		}),
	)
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
