package restasks

import (
	"os"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBadgerStore(t *testing.T) {
	store := NewBadgerStore(openTestBadgerDB(t), WithBadgerStoreRetention(time.Second))

	in := &Task{
		ID:        uuid.New(),
		Progress:  50,
		Error:     assert.AnError,
		Result:    true,
		StartedAt: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		EndedAt:   time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	require.NoError(t, store.Write(in))

	got, err := store.ReadByID(in.ID)
	require.NoError(t, err)
	require.Equal(t, in, got)
}

func openTestBadgerDB(t *testing.T) *badger.DB {
	path, err := os.MkdirTemp("/tmp/", "restasks-badger-store-")
	require.NoError(t, err)

	result, err := badger.Open(badger.DefaultOptions(path))
	require.NoError(t, err)

	return result
}
