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

		pager := pagination.NewPager(readIDsPage, pagination.WithPagerLimit(1))

		for pager.Next() {
			assert.Len(t, pager.Page(), 1)
			pagesCount++
		}

		assert.NoError(t, pager.Err())
		assert.Equal(t, 3, pagesCount)
	})

	t.Run("first page is shorter than the limit", func(t *testing.T) {
		pagesCount := 0

		pager := pagination.NewPager(func(limit, offset int) (uuid.UUIDs, error) {
			if offset != 0 {
				return nil, nil
			}

			return uuid.UUIDs{uuid.New()}, nil
		}, pagination.WithPagerLimit(2))

		for pager.Next() {
			pagesCount++
		}

		assert.NoError(t, pager.Err())
		assert.Equal(t, 1, pagesCount)
	})

	t.Run("with error", func(t *testing.T) {
		pager := pagination.NewPager(func(limit, offset int) (uuid.UUIDs, error) {
			return nil, assert.AnError
		})

		assert.False(t, pager.Next())
		assert.ErrorIs(t, pager.Err(), assert.AnError)
	})
}

func TestBoundLimit(t *testing.T) {
	tests := map[string]struct {
		in, want int
	}{
		"less than limit":    {in: 5, want: 5},
		"equal to limit":     {in: 10, want: 10},
		"greater than limit": {in: 15, want: 10},
		"less than zero":     {in: -5, want: 10},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, pagination.NewLimit(tt.in).Bound(10))
		})
	}
}

func TestBoundOffset(t *testing.T) {
	tests := map[string]struct {
		in, want int
	}{
		"less than zero":    {in: -5, want: 0},
		"greater than zero": {in: 5, want: 5},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, pagination.NewOffset(tt.in).Bound())
		})
	}
}

func readIDsPage(limit, offset int) (uuid.UUIDs, error) {
	if offset >= 3 {
		return nil, nil
	}

	return uuid.UUIDs{uuid.New()}, nil
}
