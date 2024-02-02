package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/ristretto"
)

const defaultRistrettoCacheCost = 0

type Ristretto struct{ baseCache *ristretto.Cache }

func NewRistretto(config *ristretto.Config) (*Ristretto, error) {
	baseCache, err := ristretto.NewCache(config)
	if err != nil {
		return nil, err
	}

	return &Ristretto{baseCache}, nil
}

var _ ReadWriter = (*Ristretto)(nil)

func (r *Ristretto) Read(key string) any {
	result, _ := r.baseCache.Get(key)

	return result
}

func (r *Ristretto) Write(key string, v any) {
	if v == nil {
		r.baseCache.Del(key)
		return
	}

	r.baseCache.Set(key, v, int64(getValueSize(v)))
}

func getValueSize(v any) int {
	var b bytes.Buffer

	_ = gob.NewEncoder(&b).Encode(v)

	return int(len(b.Bytes()))
}
