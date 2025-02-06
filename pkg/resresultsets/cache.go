package resresultsets

import (
	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/cache"
	"github.com/loungeup/go-loungeup/pkg/errors"
)

type CacheStore struct {
	cache cache.ReadWriter
}

func NewCacheStore(cache cache.ReadWriter) *CacheStore {
	return &CacheStore{cache: cache}
}

var _ (Store) = (*CacheStore)(nil)

func (store *CacheStore) ReadByID(id uuid.UUID) (*ResultSet, error) {
	if result, ok := store.cache.Read("result-sets." + id.String()).(*ResultSet); ok {
		return result, nil
	}

	return nil, &errors.Error{Code: errors.CodeNotFound}
}

func (store *CacheStore) Write(set *ResultSet) error {
	store.cache.Write("result-sets."+set.ID.String(), set)

	return nil
}
