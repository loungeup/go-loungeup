package models

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
)

type Entity struct {
	ID                uuid.UUID           `json:"id"`
	LegacyID          int                 `json:"legacyId,omitempty"`
	Type              EntityType          `json:"type"`
	Name              string              `json:"name,omitempty"`
	Slug              string              `json:"slug,omitempty"`
	Image             string              `json:"image,omitempty"`
	Languages         DataValue[[]string] `json:"languages,omitempty"`
	Timezone          string              `json:"timezone,omitempty"`
	Country           string              `json:"country,omitempty"`
	PostalCode        string              `json:"postalCode,omitempty"`
	City              string              `json:"city,omitempty"`
	Address           string              `json:"address,omitempty"`
	Rooms             int                 `json:"rooms,omitempty"`
	Currency          res.SoftRef         `json:"currency,omitempty"`
	ConvertAmounts    bool                `json:"convertAmounts,omitempty"`
	IndexGuestProfile bool                `json:"indexGuestProfile,omitempty"`
	Chain             res.SoftRef         `json:"chain,omitempty"`
	Group             res.SoftRef         `json:"group,omitempty"`
	Reseller          res.SoftRef         `json:"reseller,omitempty"`
	Teams             res.SoftRef         `json:"teams,omitempty"`
	CreatedAt         time.Time           `json:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt"`
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

func (t EntityType) String() string { return string(t) }

type EntitySelector struct {
	EntityID uuid.UUID
}

func (s EntitySelector) RID() string { return "authority.entities." + s.EntityID.String() }

type EntityAccountsSelector struct {
	EntityID      uuid.UUID
	Limit, Offset int
	Name          string
}

func (s EntityAccountsSelector) EncodedQuery() string {
	query := url.Values{}

	sanitizedLimit := 25
	if s.Limit > 0 {
		sanitizedLimit = s.Limit
	}
	query.Set("limit", strconv.Itoa(sanitizedLimit))

	sanitizedOffset := 0
	if s.Offset > 0 {
		sanitizedOffset = s.Offset
	}
	query.Set("offset", strconv.Itoa(sanitizedOffset))

	if s.Name != "" {
		query.Set("name", s.Name)
	}

	return query.Encode()
}

func (s EntityAccountsSelector) RID() string {
	return "authority.entities." + s.EntityID.String() + ".accounts"
}

type EntityUpdates struct {
	ConvertAmountsTaskRID    string `json:"convertAmountsTaskRid"`
	IndexGuestProfile        bool   `json:"indexGuestProfile"`
	IndexGuestProfileTaskRID string `json:"indexGuestProfileTaskRid"`
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

type BuildEntityESQueryParams struct {
	Conditions         *SearchConditions `json:"conditions,omitempty"`
	RawConditions      json.RawMessage   `json:"-"`
	GuestIDs           uuid.UUIDs        `json:"guestIds,omitempty"`
	DisableEntityScope bool              `json:"disableEntityScope,omitempty"`
}

var _ (json.Marshaler) = (*BuildEntityESQueryParams)(nil)

func (p *BuildEntityESQueryParams) MarshalJSON() ([]byte, error) {
	if p.RawConditions != nil {
		return json.Marshal(map[string]any{
			"conditions":         p.RawConditions,
			"guestIds":           p.GuestIDs.Strings(),
			"disableEntityScope": p.DisableEntityScope,
		})
	}

	type Alias BuildEntityESQueryParams

	return json.Marshal((*Alias)(p))
}
