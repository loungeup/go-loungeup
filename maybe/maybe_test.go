package maybe_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/maybe"
	"github.com/loungeup/go-loungeup/translations"
	"github.com/stretchr/testify/assert"
)

func TestOptionalValues(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, 1, maybe.Int(0).Or(1))
		assert.Equal(t, 1, maybe.Int(1).Or(2))
	})

	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "hello", maybe.String("").Or("hello"))
		assert.Equal(t, "hello", maybe.String("hello").Or("world"))
	})

	t.Run("time", func(t *testing.T) {
		now := time.Now()
		assert.Equal(t, maybe.Time(time.Time{}).Or(now), now)
		assert.Equal(t, maybe.Time(now).Or(time.Time{}), now)
	})

	t.Run("translations", func(t *testing.T) {
		assert.Equal(t,
			translations.Translations{"en": "Hello"}, maybe.Translations(nil).Or(translations.Translations{"en": "Hello"}),
		)
		assert.Equal(t,
			translations.Translations{"en": "Hello"}, maybe.Translations(translations.Translations{"en": "Hello"}).Or(translations.Translations{"en": "Hi"}),
		)
	})

	t.Run("uuid", func(t *testing.T) {
		id := uuid.New()
		assert.Equal(t, maybe.UUID(uuid.Nil).Or(id), id)
		assert.Equal(t, maybe.UUID(id).Or(uuid.Nil), id)
	})
}
