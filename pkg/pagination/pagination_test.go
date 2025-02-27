package pagination

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPager(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		pagesCount := 0

		pageReader := NewOffsetPageReader(readIDsPage)
		pager := NewPager(pageReader, WithPageSize(1))

		for pager.Next() {
			assert.Len(t, pager.Page(), 1)

			pagesCount++
		}

		assert.NoError(t, pager.Err())
		assert.Equal(t, 3, pagesCount)

		pager.Reset()
		assert.Equal(t, 0, pageReader.offset)
		assert.Nil(t, pager.lastErr)
		assert.Nil(t, pager.lastPage)
	})

	t.Run("first page is shorter than the size", func(t *testing.T) {
		pagesCount := 0

		pager := NewPager(NewOffsetPageReader(func(size, offset int) (uuid.UUIDs, error) {
			if offset != 0 {
				return nil, nil
			}

			return uuid.UUIDs{uuid.New()}, nil
		}), WithPageSize(2))

		for pager.Next() {
			pagesCount++
		}

		assert.NoError(t, pager.Err())
		assert.Equal(t, 1, pagesCount)
	})

	t.Run("with error", func(t *testing.T) {
		pager := NewPager(NewOffsetPageReader(func(size, offset int) (uuid.UUIDs, error) {
			return nil, assert.AnError
		}))

		assert.False(t, pager.Next())
		assert.ErrorIs(t, pager.Err(), assert.AnError)
	})

	t.Run("keyset", func(t *testing.T) {
		pageReader := NewKeysetPageReader(func(size int, lastKey string) (uuid.UUIDs, string, error) {
			if lastKey != "" {
				return nil, "", nil
			}

			return uuid.UUIDs{uuid.New()}, "foo", nil
		})
		pager := NewPager(pageReader)

		for pager.Next() {
		}

		assert.NoError(t, pager.Err())

		pager.Reset()
		assert.Equal(t, "", pageReader.lastKey)
		assert.Nil(t, pager.lastErr)
		assert.Nil(t, pager.lastPage)
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
			assert.Equal(t, tt.want, NewLimit(tt.in).Bound(10))
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
			assert.Equal(t, tt.want, NewOffset(tt.in).Bound())
		})
	}
}

func readIDsPage(_, offset int) (uuid.UUIDs, error) {
	if offset >= 3 {
		return nil, nil
	}

	return uuid.UUIDs{uuid.New()}, nil
}
