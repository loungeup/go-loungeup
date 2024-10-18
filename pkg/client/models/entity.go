package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
)

type Entity struct {
	ID             uuid.UUID           `json:"id"`
	LegacyID       int                 `json:"legacyId,omitempty"`
	Type           EntityType          `json:"type"`
	Name           string              `json:"name,omitempty"`
	Slug           string              `json:"slug,omitempty"`
	Image          string              `json:"image,omitempty"`
	Languages      DataValue[[]string] `json:"languages,omitempty"`
	Timezone       string              `json:"timezone,omitempty"`
	Country        string              `json:"country,omitempty"`
	PostalCode     string              `json:"postalCode,omitempty"`
	City           string              `json:"city,omitempty"`
	Address        string              `json:"address,omitempty"`
	Rooms          int                 `json:"rooms,omitempty"`
	Currency       res.SoftRef         `json:"currency,omitempty"`
	ConvertAmounts bool                `json:"convertAmounts,omitempty"`
	Chain          res.SoftRef         `json:"chain,omitempty"`
	Group          res.SoftRef         `json:"group,omitempty"`
	Reseller       res.SoftRef         `json:"reseller,omitempty"`
	Teams          res.SoftRef         `json:"teams,omitempty"`
	CreatedAt      time.Time           `json:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt"`
}

func (e Entity) ChainID() uuid.UUID    { return getEntityIDFromRID(string(e.Chain)) }
func (e Entity) GroupID() uuid.UUID    { return getEntityIDFromRID(string(e.Group)) }
func (e Entity) ResellerID() uuid.UUID { return getEntityIDFromRID(string(e.Reseller)) }

func (e Entity) CurrencyCode() string {
	if e.Currency == "" {
		return ""
	}

	return strings.TrimPrefix(string(e.Currency), "authority.currencies.")
}

type EntityType string

const (
	EntityTypeAccount  EntityType = "account"
	EntityTypeChain    EntityType = "chain"
	EntityTypeGroup    EntityType = "group"
	EntityTypeReseller EntityType = "reseller"
)

type EntitySelector struct {
	EntityID uuid.UUID
}

func (s EntitySelector) RID() string { return "authority.entities." + s.EntityID.String() }

type EntityAccountsSelector struct {
	EntityID      uuid.UUID
	Limit, Offset int
}

func (s EntityAccountsSelector) EncodedQuery() string {
	sanitizedLimit := 25
	if s.Limit > 0 {
		sanitizedLimit = s.Limit
	}

	sanitizedOffset := 0
	if s.Offset > 0 {
		sanitizedOffset = s.Offset
	}

	return "limit=" + strconv.Itoa(sanitizedLimit) + "&offset=" + strconv.Itoa(sanitizedOffset)
}

func (s EntityAccountsSelector) RID() string {
	return "authority.entities." + s.EntityID.String() + ".accounts"
}

type EntityUpdates struct {
	ConvertAmountsTaskRID string `json:"convertAmountsTaskRid"`
}

func getEntityIDFromRID(rid string) uuid.UUID {
	if rid == "" {
		return uuid.Nil
	}

	result, err := uuid.Parse(strings.TrimPrefix(string(rid), "authority.entities."))
	if err != nil {
		return uuid.Nil
	}

	return result
}
