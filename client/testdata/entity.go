package testdata

import (
	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/resmodels"
)

var (
	// Définir d'abord les ID qui sont utilisés dans plusieurs endroits
	EntityChainID = uuid.New()
	EntityGroupID = uuid.New()
	EntityID      = uuid.New()

	// Ensuite définir les sélecteurs qui utilisent ces ID
	EntitySelector = &resmodels.EntitySelector{
		EntityID: EntityID,
	}

	AccountChainSelector = &resmodels.EntitySelector{
		EntityID: EntityChainID,
	}

	AccountGroupSelector = &resmodels.EntitySelector{
		EntityID: EntityGroupID,
	}

	// Maintenant définir les entités complètes
	EntityChain = &resmodels.Entity{
		ID:   EntityChainID.String(),
		Type: resmodels.EntityTypeChain,
	}

	EntityGroup = &resmodels.Entity{
		ID:   EntityGroupID.String(),
		Type: resmodels.EntityTypeGroup,
	}

	Entity = &resmodels.Entity{
		ID:       EntityID.String(),
		LegacyID: 1,
		Type:     resmodels.EntityTypeAccount,
		Name:     "Test Account",
		Slug:     "testaccount",
		Image:    "https://example.com/image.jpg",
		// Languages:      models.NewDataValue([]string{"en"}),
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
		// CreatedAt:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		// UpdatedAt:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Définir les autres variables qui dépendent des précédentes
	EntityAccountsSelector = &resmodels.EntityAccountsSelector{
		EntityID: EntityID,
	}

	// Chaînes JSON qui dépendent des variables définies plus haut
	EntityCollection = `[
		{"rid": "` + EntitySelector.RID() + `"}
	]`

	EntityModel = `{
		"id": "` + EntityID.String() + `",
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
		"id": "` + EntityChainID.String() + `",
		"type": "chain"
	}`

	GroupModel = `{
		"id": "` + EntityGroupID.String() + `",
		"type": "group"
	}`
)
