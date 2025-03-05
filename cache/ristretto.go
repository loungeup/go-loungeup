package cache

import (
	"bytes"
	"encoding/gob"
	"log/slog"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/loungeup/go-loungeup/log"
)

const (
	defaultRistrettoCacheCost     = 0
	defaultRistrettoCacheDuration = 5 * time.Minute
)

// RistrettoCacheSize provides an abstraction to simplify the configuration of the cache. A medium cache can hold
// approximatively 32MB, while a large cache can hold approximatively 64MB.
type RistrettoCacheSize int

const (
	tooSmallRistrettoCache RistrettoCacheSize = iota + 1 // Used to test the case where the item is too large.
	MediumRistrettoCache
	LargeRistrettoCache
	VeryLargeRistrettoCache
)

// Config of the ristretto cache for the given size.
func (s RistrettoCacheSize) Config() *ristretto.Config[string, any] {
	const (
		bufferItems = 64 // The documentation recommends this value.

		veryLargeRistrettoCacheMaxCost     = 128_000_000 // 128MB.
		veryLargeRistrettoCacheNumCounters = 20_000_000  // To hold approximatively 2,000,000 keys.

		largeRistrettoCacheMaxCost     = 64_000_000 // 64MB.
		largeRistrettoCacheNumCounters = 10_000_000 // To hold approximatively 1,000,000 keys.

		mediumRistrettoCacheMaxCost     = 32_000_000 // 32MB.
		mediumRistrettoCacheNumCounters = 5_000_000  // To hold approximatively 500,000 keys.

		tooSmallRistrettoCacheMaxCost     = 10
		tooSmallRistrettoCacheNumCounters = 100
	)

	return &ristretto.Config[string, any]{
		NumCounters: func() int64 {
			switch s {
			case tooSmallRistrettoCache:
				return tooSmallRistrettoCacheNumCounters
			case LargeRistrettoCache:
				return largeRistrettoCacheNumCounters
			case VeryLargeRistrettoCache:
				return veryLargeRistrettoCacheNumCounters
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
			case VeryLargeRistrettoCache:
				return veryLargeRistrettoCacheMaxCost
			default:
				return mediumRistrettoCacheMaxCost
			}
		}(),
		BufferItems: bufferItems,
		Metrics:     true,
	}
}

type Ristretto struct {
	baseCache *ristretto.Cache[string, any]
	logger    *log.Logger
}

type RistrettoOption func(*Ristretto)

// NewRistretto creates a new ristretto cache with the given size.
func NewRistretto(size RistrettoCacheSize, options ...RistrettoOption) (*Ristretto, error) {
	baseCache, err := ristretto.NewCache(size.Config())
	if err != nil {
		return nil, err
	}

	result := &Ristretto{
		baseCache: baseCache,
		logger:    log.Default(),
	}
	for _, option := range options {
		option(result)
	}

	return result, nil
}

func WithRistrettoLogger(logger *log.Logger) RistrettoOption {
	return func(r *Ristretto) { r.logger = logger }
}

var _ ReadWriter = (*Ristretto)(nil)

func (r *Ristretto) Read(key string) any {
	result, _ := r.baseCache.Get(key)

	return result
}

func (r *Ristretto) Size() uint64 {
	return r.baseCache.Metrics.CostAdded() - r.baseCache.Metrics.CostEvicted()
}

func (r *Ristretto) Write(key string, value any) {
	r.WriteWithDuration(key, value, defaultRistrettoCacheDuration)
}

func (r *Ristretto) WriteWithDuration(key string, value any, duration time.Duration) {
	if value == nil {
		r.baseCache.Del(key)

		return
	}

	cost := getRistrettoValueCost(value)

	if !r.baseCache.SetWithTTL(key, value, cost, duration) {
		r.logger.Error("Could not cache value",
			slog.String("key", key),
			slog.Any("value", value),
			slog.Int64("cost", cost),
			slog.Duration("duration", duration),
		)
	}
}

func getRistrettoValueCost(value any) int64 {
	var buffer bytes.Buffer

	_ = gob.NewEncoder(&buffer).Encode(value)

	return int64(len(buffer.Bytes()))
}
