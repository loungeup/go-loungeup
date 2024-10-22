package maybe_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/maybe"
	"github.com/loungeup/go-loungeup/pkg/translations"
	"github.com/stretchr/testify/assert"
)

func TestOptionalValues(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, maybe.Int(0).Or(1), 1)
		assert.Equal(t, maybe.Int(1).Or(2), 1)
	})

	t.Run("string", func(t *testing.T) {
		assert.Equal(t, maybe.String("").Or("hello"), "hello")
		assert.Equal(t, maybe.String("hello").Or("world"), "hello")
	})

	t.Run("time", func(t *testing.T) {
		now := time.Now()
		assert.Equal(t, maybe.Time(time.Time{}).Or(now), now)
		assert.Equal(t, maybe.Time(now).Or(time.Time{}), now)
	})

	t.Run("translations", func(t *testing.T) {
		assert.Equal(t,
			maybe.Translations(nil).Or(translations.Translations{"en": "Hello"}),
			translations.Translations{"en": "Hello"},
		)
		assert.Equal(t,
			maybe.Translations(translations.Translations{"en": "Hello"}).Or(translations.Translations{"en": "Hi"}),
			translations.Translations{"en": "Hello"},
		)
	})

	t.Run("uuid", func(t *testing.T) {
		id := uuid.New()
		assert.Equal(t, maybe.UUID(uuid.Nil).Or(id), id)
		assert.Equal(t, maybe.UUID(id).Or(uuid.Nil), id)
	})
}
