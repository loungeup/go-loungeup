package restasks

import (
	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/cache"
	"github.com/loungeup/go-loungeup/errors"
)

type CacheStore struct {
	cache cache.ReadWriter
}

func NewCacheStore(cache cache.ReadWriter) *CacheStore {
	return &CacheStore{cache: cache}
}

var _ (Store) = (*CacheStore)(nil)

func (c *CacheStore) ReadByID(id uuid.UUID) (*Task, error) {
	if result, ok := c.cache.Read("tasks." + id.String()).(*Task); ok {
		return result, nil
	}

	return nil, &errors.Error{Code: errors.CodeNotFound}
}

func (c *CacheStore) Write(task *Task) error {
	c.cache.Write("tasks."+task.ID.String(), task)

	return nil
}
