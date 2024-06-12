package restasks

import (
	"sync"

	"github.com/loungeup/go-loungeup/pkg/errors"
)

type TaskReadWriter interface {
	TaskReader
	TaskWriter
}

type TaskReader interface {
	ReadTask(rid string) (*task, error)
}

type TaskWriter interface {
	WriteTask(task *task) error
}

type inMemoryStore struct {
	capacity int
	tasks    sync.Map
}

const defaultInMemoryStoreCapacity = 1000

// NewInMemoryStore create a new in-memory tasks store with the given capacity.
// If the capacity is 0, the default capacity will be used.
func NewInMemoryStore(capacity int) *inMemoryStore {
	if capacity == 0 {
		capacity = defaultInMemoryStoreCapacity
	}

	return &inMemoryStore{
		capacity: capacity,
		tasks:    sync.Map{},
	}
}

var _ TaskReadWriter = (*inMemoryStore)(nil)

func (s *inMemoryStore) ReadTask(rid string) (*task, error) {
	taskRecord, ok := s.tasks.Load(rid)
	if !ok {
		return nil, &errors.Error{Code: errors.CodeNotFound, Message: "Task not found"}
	}

	result, ok := taskRecord.(*task)
	if !ok {
		return nil, &errors.Error{Code: errors.CodeInternal, Message: "Task is not a task"}
	}

	return result, nil
}

func (s *inMemoryStore) WriteTask(task *task) error {
	if s.isFull() {
		s.deleteOldestTask()
	}

	s.tasks.Store(task.rid(), task)

	return nil
}

func (s *inMemoryStore) deleteOldestTask() { s.deleteTask(s.readOldestTask().rid()) }

func (s *inMemoryStore) deleteTask(rid string) { s.tasks.Delete(rid) }

func (s *inMemoryStore) isFull() bool { return s.len() >= s.capacity }

func (s *inMemoryStore) len() int {
	result := 0

	s.tasks.Range(func(_, _ any) bool {
		result++

		return true
	})

	return result
}

func (s *inMemoryStore) readOldestTask() *task {
	result := &task{}

	s.tasks.Range(func(_, value any) bool {
		task, ok := value.(*task)
		if !ok {
			return true
		}

		if result.CreatedAt.IsZero() || task.CreatedAt.Before(result.CreatedAt) {
			result = task
		}

		return true
	})

	return result
}
