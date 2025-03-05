package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRistrettoCache(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cache, err := NewRistretto(0)
		require.NoError(t, err)

		cache.Write("foo", "bar")
		waitForCache()

		require.Equal(t, "bar", cache.Read("foo"))

		cache.Write("foo", nil)
		waitForCache()

		require.Nil(t, cache.Read("foo"))
	})

	t.Run("too large item", func(t *testing.T) {
		cache, err := NewRistretto(tooSmallRistrettoCache)
		require.NoError(t, err)

		cache.Write("baz", "this item is too large")
		waitForCache()

		require.Equal(t, nil, cache.Read("baz"))
		require.Equal(t, int64(7), getRistrettoValueCost("bar"))
	})
}

func TestGetRistrettoValueCost(t *testing.T) {
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

	for i := 0; i < len(tests); i++ {
		if i > 0 {
			require.Greater(t, getRistrettoValueCost(tests[i]), getRistrettoValueCost(tests[i-1]))
		}
	}
}

func waitForCache() { time.Sleep(100 * time.Millisecond) }
