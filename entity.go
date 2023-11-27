package loungeup

import (
	"time"

	"github.com/google/uuid"
)

// Entity represents a LoungeUp entity. It is provided by the Authority server.
type Entity struct {
	// ID of the entity.
	ID uuid.UUID `json:"id"`

	// LegacyID of the entity. Used to match the entity with the Legacy system.
	LegacyID int `json:"legacyId"`

	// Type of the entity ("account", "chain", "group" or "reseller").
	Type string `json:"type"`

	// Name of the entity.
	Name string `json:"name"`

	// Slug generated from the name of the entity.
	Slug string `json:"slug"`

	// Image URL of the entity.
	Image string `json:"image"`

	// Languages used by the entity as ISO 639-1 codes.
	Languages []string `json:"languages"`

	// Timezone of the entity (e.g. "Europe/Paris").
	Timezone string `json:"timezone"`

	// Country of the entity as ISO 3166-1 alpha-2 code (e.g. "FR").
	Country string `json:"country"`

	// PostalCode of the entity (e.g. "31520").
	PostalCode string `json:"postalCode"`

	// City of the entity (e.g. "Ramonville-Saint-Agne").
	City string `json:"city"`

	// Address of the entity (e.g. "12 avenue de l'Europe").
	Address string `json:"address"`

	// Rooms count of the entity.
	Rooms int `json:"rooms"`

	// CurrencyCode of the entity (e.g. "eur").
	CurrencyCode string `json:"currencyCode"`

	// ChainID of the entity. If the entity is not part of a chain, this field is equal to uuid.Nil.
	ChainID uuid.UUID `json:"chainId"`

	// GroupID of the entity. If the entity is not part of a group, this field is equal to uuid.Nil.
	GroupID uuid.UUID `json:"groupId"`

	// ResellerID of the entity. If the entity is not part of a reseller, this field is equal to uuid.Nil.
	ResellerID uuid.UUID `json:"resellerId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RelatedIDs returns the IDs of the entity and its related entities.
func (e Entity) RelatedIDs() []uuid.UUID {
	result := []uuid.UUID{e.ID}

	if e.ChainID != uuid.Nil {
		result = append(result, e.ChainID)
	}

	if e.GroupID != uuid.Nil {
		result = append(result, e.GroupID)
	}

	if e.ResellerID != uuid.Nil {
		result = append(result, e.ResellerID)
	}

	return result
}

// EntitySelector used to select an entity.
type EntitySelector struct {
	ID uuid.UUID `json:"id"`
}

type EntitiesReadWriter interface {
	EntitiesReader
	EntitiesWriter
}

type EntitiesReader interface {
	// ReadEntity returns the entity matching the given selector.
	ReadEntity(givenSelector EntitySelector) (*Entity, error)
}

type EntitiesWriter interface{}
