package resmodels

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/client/models"
)

type EntityType string

const (
	EntityTypeAccount  EntityType = "account"
	EntityTypeChain    EntityType = "chain"
	EntityTypeGroup    EntityType = "group"
	EntityTypeReseller EntityType = "reseller"
)

type Entity struct {
	ID                    string                   `json:"id,omitempty"`
	LegacyID              int                      `json:"legacyId,omitempty"`
	Type                  EntityType               `json:"type,omitempty"`
	Name                  string                   `json:"name,omitempty"`
	Slug                  string                   `json:"slug,omitempty"`
	Image                 string                   `json:"image,omitempty"`
	Languages             *res.DataValue[[]string] `json:"languages,omitempty"`
	Timezone              string                   `json:"timezone,omitempty"`
	Country               string                   `json:"country,omitempty"`
	PostalCode            string                   `json:"postalCode,omitempty"`
	City                  string                   `json:"city,omitempty"`
	Address               string                   `json:"address,omitempty"`
	Rooms                 int                      `json:"rooms,omitempty"`
	Currency              res.SoftRef              `json:"currency,omitempty"`
	ConvertAmounts        bool                     `json:"convertAmounts"`
	ConvertAmountsTask    res.SoftRef              `json:"convertAmountsTask,omitempty"`
	IndexGuestProfileTask res.SoftRef              `json:"indexGuestProfileTask,omitempty"`
	Chain                 res.SoftRef              `json:"chain,omitempty"`
	Group                 res.SoftRef              `json:"group,omitempty"`
	Reseller              res.SoftRef              `json:"reseller,omitempty"`
	Teams                 res.SoftRef              `json:"teams,omitempty"`
	CreatedAt             string                   `json:"createdAt,omitempty"`
	UpdatedAt             string                   `json:"updatedAt,omitempty"`
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

func (t EntityType) String() string { return string(t) }

type EntitySelector struct {
	EntityID uuid.UUID
}

func (s EntitySelector) RID() string { return "authority.entities." + s.EntityID.String() }

type EntityID uuid.UUID

func (s EntityID) RID() string {
	return "authority.entities." + uuid.UUID(s).String()
}

func (s EntityID) String() string {
	return uuid.UUID(s).String()
}

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

type BuildEntityESQueryParams struct {
	Conditions         *models.SearchConditions `json:"conditions,omitempty"`
	RawConditions      json.RawMessage          `json:"-"`
	GuestIDs           uuid.UUIDs               `json:"guestIds,omitempty"`
	DisableEntityScope bool                     `json:"disableEntityScope,omitempty"`
}

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
