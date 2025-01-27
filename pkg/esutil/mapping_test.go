package esutil_test

import (
	"testing"

	"github.com/loungeup/go-loungeup/pkg/esutil"
	"github.com/stretchr/testify/assert"
)

func TestMappingKeys(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		keys := esutil.GlobalMappingKeys()
		assert.Equal(t, "booking.id", keys.Booking.ID)
		assert.Equal(t, "guest.account.id", keys.Guest.Account.ID)
		assert.Equal(t, "guest.chain.id", keys.Guest.Chain.ID)
		assert.Equal(t, "guest.group.id", keys.Guest.Group.ID)
	})

	t.Run("scoped", func(t *testing.T) {
		keys, err := esutil.NewScopedMappingKeys(esutil.MappingKeysScopeAccount)
		assert.NoError(t, err)
		assert.Equal(t, "booking.id", keys.Booking.ID)
		assert.Equal(t, "guest.account.id", keys.Guest.ID)
	})

	t.Run("invalid scope", func(t *testing.T) {
		_, err := esutil.NewScopedMappingKeys(esutil.MappingKeysScopeUnknown)
		assert.Error(t, err)
	})
}
