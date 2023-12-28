package cache

import (
	"github.com/dgraph-io/ristretto"
)

const defaultRistrettoCacheCost = 0

type Ristretto struct {
	underlyingCache *ristretto.Cache
}

func NewRistretto(config *ristretto.Config) (*Ristretto, error) {
	underlyingCache, err := ristretto.NewCache(config)
	if err != nil {
		return nil, err
	}

	return &Ristretto{underlyingCache}, nil
}

var _ ReadWriter = (*Ristretto)(nil)

func (r *Ristretto) Read(key string) any {
	result, _ := r.underlyingCache.Get(key)

	return result
}

func (r *Ristretto) Write(key string, value any) {
	r.underlyingCache.Set(key, value, defaultRistrettoCacheCost)
}
