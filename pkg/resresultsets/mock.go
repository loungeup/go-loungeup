package resresultsets

import "github.com/google/uuid"

type MockStore struct {
	ReadByIDFunc func(id uuid.UUID) (*ResultSet, error)
	WriteFunc    func(set *ResultSet) error
}

var _ (Store) = (*MockStore)(nil)

func (store *MockStore) ReadByID(id uuid.UUID) (*ResultSet, error) { return store.ReadByIDFunc(id) }
func (store *MockStore) Write(set *ResultSet) error                { return store.WriteFunc(set) }
