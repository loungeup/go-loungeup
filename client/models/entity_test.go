package models

import (
	"testing"

	"github.com/jirenius/go-res"
	"github.com/stretchr/testify/assert"
)

func TestEntityCurrencyCode(t *testing.T) {
	entity := &Entity{
		Currency: "authority.currencies.usd",
	}

	assert.Equal(t, "usd", entity.CurrencyCode())
}

func TestEntityRelationIDs(t *testing.T) {
	entity := &Entity{
		Chain:    res.SoftRef("authority.entities.a3a8e725-cb7f-4c90-a936-23f854b267cf"),
		Group:    res.SoftRef("authority.entities.1d8b8d07-e809-4563-aa18-de5d83830a33"),
		Reseller: res.SoftRef("authority.entities.4d21bca5-0773-4e77-b824-2b751d863663"),
	}

	assert.Equal(t, "a3a8e725-cb7f-4c90-a936-23f854b267cf", entity.ChainID().String())
	assert.Equal(t, "1d8b8d07-e809-4563-aa18-de5d83830a33", entity.GroupID().String())
	assert.Equal(t, "4d21bca5-0773-4e77-b824-2b751d863663", entity.ResellerID().String())
}
