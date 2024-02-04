package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/dgraph-io/ristretto"
)

const (
	defaultRistrettoCacheCost = 0
	defaultRistrettoCacheTTL  = 5 * time.Minute
)

// RistrettoCacheSize provides an abstraction to simplify the configuration of the cache. A medium cache can hold
// approximatively 32MB, while a large cache can hold approximatively 64MB.
type RistrettoCacheSize int

const (
	tooSmallRistrettoCache RistrettoCacheSize = iota // Used to test the case where the item is too large.
	MediumRistrettoCache
	LargeRistrettoCache
)

// Config of the ristretto cache for the given size.
func (s RistrettoCacheSize) Config() *ristretto.Config {
	const (
		bufferItems = 64 // The documentation recommends this value.

		largeRistrettoCacheMaxCost     = 64_000_000 // 64MB.
		largeRistrettoCacheNumCounters = 10_000_000 // To hold approximatively 1,000,000 keys.

		mediumRistrettoCacheMaxCost     = 32_000_000 // 32MB.
		mediumRistrettoCacheNumCounters = 5_000_000  // To hold approximatively 500,000 keys.

		tooSmallRistrettoCacheMaxCost     = 10
		tooSmallRistrettoCacheNumCounters = 100
	)

	return &ristretto.Config{
		NumCounters: func() int64 {
			switch s {
			case tooSmallRistrettoCache:
				return tooSmallRistrettoCacheNumCounters
			case LargeRistrettoCache:
				return largeRistrettoCacheNumCounters
			default:
				return mediumRistrettoCacheNumCounters
			}
		}(),
		MaxCost: func() int64 {
			switch s {
			case tooSmallRistrettoCache:
				return tooSmallRistrettoCacheMaxCost
			case LargeRistrettoCache:
				return largeRistrettoCacheMaxCost
			default:
				return mediumRistrettoCacheMaxCost
			}
		}(),
		BufferItems: bufferItems,
	}
}

type Ristretto struct{ baseCache *ristretto.Cache }

// NewRistretto creates a new ristretto cache with the given size.
func NewRistretto(s RistrettoCacheSize) (*Ristretto, error) {
	baseCache, err := ristretto.NewCache(s.Config())
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

	r.baseCache.SetWithTTL(key, v, getRistrettoValueCost(v), defaultRistrettoCacheTTL)
}

func getRistrettoValueCost(v any) int64 {
	var b bytes.Buffer

	_ = gob.NewEncoder(&b).Encode(v)

	return int64(len(b.Bytes()))
}
