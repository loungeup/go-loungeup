//nolint:gomnd,mnd
package testdata

import (
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/client/models"
)

var (
	EntityChain = &models.Entity{
		ID:   uuid.New(),
		Type: models.EntityTypeChain,
	}

	EntityGroup = &models.Entity{
		ID:   uuid.New(),
		Type: models.EntityTypeGroup,
	}

	Entity = &models.Entity{
		ID:             uuid.New(),
		LegacyID:       1,
		Type:           models.EntityTypeAccount,
		Name:           "Test Account",
		Slug:           "testaccount",
		Image:          "https://example.com/image.jpg",
		Languages:      models.NewDataValue([]string{"en"}),
		Timezone:       "Europe/Paris",
		Country:        "FR",
		PostalCode:     "31520",
		City:           "Ramonville-Saint-Agne",
		Address:        "12 avenue de l'Europe",
		Rooms:          100,
		Currency:       res.SoftRef("authority.currencies.eur"),
		ConvertAmounts: true,
		Chain:          res.SoftRef(AccountChainSelector.RID()),
		Group:          res.SoftRef(AccountGroupSelector.RID()),
		CreatedAt:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	EntityAccountsSelector = &models.EntityAccountsSelector{
		EntityID: Entity.ID,
	}

	EntityCollection = `[
		{"rid": "` + EntitySelector.RID() + `"}
	]`

	EntityModel = `{
		"id": "` + Entity.ID.String() + `",
		"legacyId": 1,
		"type": "account",
		"name": "Test Account",
		"slug": "testaccount",
		"image": "https://example.com/image.jpg",
		"languages": {
			"data": ["en"]
		},
		"timezone": "Europe/Paris",
		"country": "FR",
		"postalCode": "31520",
		"city": "Ramonville-Saint-Agne",
		"address": "12 avenue de l'Europe",
		"rooms": 100,
		"currency": {
			"rid": "authority.currencies.eur"
		},
		"convertAmounts": true,
		"createdAt": "2020-01-01T00:00:00Z",
		"updatedAt": "2020-01-01T00:00:00Z",
		"chain": {
			"rid": "` + AccountChainSelector.RID() + `"
		},
		"group": {
			"rid": "` + AccountGroupSelector.RID() + `"
		}
	}`

	ChainModel = `{
		"id": "` + EntityChain.ID.String() + `",
		"type": "chain"
	}`

	GroupModel = `{
		"id": "` + EntityGroup.ID.String() + `",
		"type": "group"
	}`

	EntitySelector = &models.EntitySelector{
		EntityID: Entity.ID,
	}

	AccountChainSelector = &models.EntitySelector{
		EntityID: EntityChain.ID,
	}

	AccountGroupSelector = &models.EntitySelector{
		EntityID: EntityGroup.ID,
	}
)
