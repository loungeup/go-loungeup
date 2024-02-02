package cache

import (
	"testing"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/stretchr/testify/require"
)

func TestRistrettoCache(t *testing.T) {
	cache, err := NewRistretto(&ristretto.Config{
		BufferItems: 64,
		MaxCost:     10,  // 10 bytes.
		NumCounters: 100, // 10 times the cache capacity.
		Metrics:     true,
	})
	require.NoError(t, err)

	t.Run("simple", func(t *testing.T) {
		cache.Write("foo", "bar")
		waitForCache()

		require.Equal(t, "bar", cache.Read("foo"))

		cache.Write("foo", nil)
		waitForCache()

		require.Nil(t, cache.Read("foo"))
	})

	t.Run("too large item", func(t *testing.T) {
		cache.Write("baz", "this item is too large")
		waitForCache()

		require.Equal(t, nil, cache.Read("baz"))
		require.Equal(t, 7, getValueSize("bar"))
	})
}

func TestGetValueSize(t *testing.T) {
	type User struct {
		FirstName string
		LastName  string
	}

	// Just make sure that bigger values have bigger sizes.
	tests := []any{
		1,
		"foo",
		[]string{"foo"},
		[]string{"foo", "bar", "baz", "qux"},
		&User{FirstName: "John", LastName: "Doe"},
		&User{FirstName: "Johnny", LastName: "Doe"},
		&User{FirstName: "Johnny", LastName: "Depp"},
	}

	previousSize := 0

	for i := 0; i < len(tests); i++ {
		size := getValueSize(tests[i])

		if i > 0 {
			require.Greater(t, size, previousSize)
		}

		previousSize = size
	}
}

func waitForCache() { time.Sleep(100 * time.Millisecond) }
