package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityCurrencyCode(t *testing.T) {
	entity := &Entity{
		Currency: "authority.currencies.usd",
	}

	assert.Equal(t, "usd", entity.CurrencyCode())
}
