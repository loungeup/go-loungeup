package restasks

import "github.com/google/uuid"

type Store interface {
	ReadByID(id uuid.UUID) (*Task, error)
	Write(task *Task) error
}
