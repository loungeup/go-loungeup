package restasks

import "github.com/google/uuid"

type MockStore struct {
	ReadByIDFunc func(id uuid.UUID) (*Task, error)
	WriteFunc    func(task *Task) error
}

var _ (Store) = (*MockStore)(nil)

func (s *MockStore) ReadByID(id uuid.UUID) (*Task, error) { return s.ReadByIDFunc(id) }
func (s *MockStore) Write(task *Task) error               { return s.WriteFunc(task) }
