package pagination_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/pagination"
	"github.com/stretchr/testify/assert"
)

func TestPager(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		pagesCount := 0

		pager := pagination.NewPager(readIDsPage, pagination.WithLimit(1))

		for pager.Next() {
			assert.Len(t, pager.Page(), 1)
			pagesCount++
		}

		assert.NoError(t, pager.Err())
		assert.Equal(t, 3, pagesCount)
	})

	t.Run("with error", func(t *testing.T) {
		pager := pagination.NewPager(func(limit, offset int) (uuid.UUIDs, error) {
			return nil, assert.AnError
		})

		assert.False(t, pager.Next())
		assert.ErrorIs(t, pager.Err(), assert.AnError)
	})
}

func readIDsPage(limit, offset int) (uuid.UUIDs, error) {
	if offset >= 3 {
		return nil, nil
	}

	return uuid.UUIDs{uuid.New()}, nil
}
