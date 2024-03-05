package client

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

// Entity as represented with the RES protocol.
type Entity struct {
	ID         uuid.UUID      `json:"id"`
	LegacyID   int            `json:"legacyId,omitempty"`
	Type       EntityType     `json:"type"`
	Name       string         `json:"name,omitempty"`
	Slug       string         `json:"slug,omitempty"`
	Image      string         `json:"image,omitempty"`
	Languages  *res.DataValue `json:"languages,omitempty"`
	Timezone   string         `json:"timezone,omitempty"`
	Country    string         `json:"country,omitempty"`
	PostalCode string         `json:"postalCode,omitempty"`
	City       string         `json:"city,omitempty"`
	Address    string         `json:"address,omitempty"`
	Rooms      int            `json:"rooms,omitempty"`
	Currency   res.SoftRef    `json:"currency,omitempty"`
	Chain      res.SoftRef    `json:"chain,omitempty"`
	Group      res.SoftRef    `json:"group,omitempty"`
	Reseller   res.SoftRef    `json:"reseller,omitempty"`
	Teams      res.Ref        `json:"teams,omitempty"`
	CreatedAt  string         `json:"createdAt"`
	UpdatedAt  string         `json:"updatedAt"`
}

type EntityType string

const (
	EntityTypeAccount  EntityType = "account"
	EntityTypeChain    EntityType = "chain"
	EntityTypeGroup    EntityType = "group"
	EntityTypeReseller EntityType = "reseller"
)

type EntityCustomFields struct {
	User  RESDataValue[map[string]EntityCustomField] `json:"user,omitempty"`
	Visit RESDataValue[map[string]EntityCustomField] `json:"visit,omitempty"`
}

type EntityCustomField struct {
	Label string                `json:"label,omitempty"`
	Type  EntityCustomFieldType `json:"type,omitempty"`
}

type EntityCustomFieldType string

const (
	EntityCustomFieldTypeBoolean EntityCustomFieldType = "boolean"
	EntityCustomFieldTypeDate    EntityCustomFieldType = "date"
	EntityCustomFieldTypeList    EntityCustomFieldType = "list"
	EntityCustomFieldTypeNumber  EntityCustomFieldType = "number"
	EntityCustomFieldTypeText    EntityCustomFieldType = "text"
)

type RoomType struct {
	ID                   uuid.UUID `json:"id"`
	EntityID             uuid.UUID `json:"entityId"`
	Name                 string    `json:"name,omitempty"`
	Code                 string    `json:"code,omitempty"`
	Capacity             int       `json:"capacity,omitempty"`
	CapacitySafetyMargin int       `json:"capacitySafetyMargin,omitempty"`
	CreatedAt            time.Time `json:"createdAt,omitempty"`
	UpdatedAt            time.Time `json:"updatedAt,omitempty"`
}

type EntitySelector struct {
	ID uuid.UUID `json:"id"`
}

func (s EntitySelector) resourceID() string {
	return authorityServiceName + ".entities." + s.ID.String()
}

// entitiesClient provides methods to interact with entities.
type entitiesClient struct{ baseClient *Client }

func (c *entitiesClient) ReadEntity(s EntitySelector) (Entity, error) {
	resourceID := s.resourceID()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(resourceID).(Entity); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[Entity](c.baseClient.resClient, resourceID)
	if err != nil {
		return Entity{}, err
	}

	defer c.baseClient.eventuallyWriteCache(resourceID, result)

	return result, nil
}

type EntityAccountsSelector struct {
	EntitySelector

	Limit, Offset int
}

func (s EntityAccountsSelector) encodeQuery() string {
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

func (s EntityAccountsSelector) resourceID() string {
	return s.EntitySelector.resourceID() + ".accounts"
}

func (c *entitiesClient) ReadEntityAccounts(s EntityAccountsSelector) ([]Entity, error) {
	encodedQuery := s.encodeQuery()
	resourceID := s.resourceID()

	cacheKey := resourceID + "?" + encodedQuery

	if cachedResult, ok := c.baseClient.eventuallyReadCache(cacheKey).([]Entity); ok {
		return cachedResult, nil
	}

	accountReferences, err := transport.GetRESCollection[res.Ref](c.baseClient.resClient, resourceID, resprot.Request{
		Query: encodedQuery,
	})
	if err != nil {
		return nil, err
	}

	result := []Entity{}

	for _, accountReference := range accountReferences {
		relatedAccount, err := transport.GetRESModel[Entity](c.baseClient.resClient, string(accountReference))
		if err != nil {
			return nil, err
		}

		result = append(result, relatedAccount)
	}

	defer c.baseClient.eventuallyWriteCache(cacheKey, result)

	return result, nil
}

func (c *entitiesClient) ReadEntityCustomFields(s EntitySelector) (EntityCustomFields, error) {
	resourceID := s.resourceID() + ".custom-fields"

	if cachedResult, ok := c.baseClient.eventuallyReadCache(resourceID).(EntityCustomFields); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[EntityCustomFields](c.baseClient.resClient, s.resourceID()+".custom-fields")
	if err != nil {
		return EntityCustomFields{}, err
	}

	defer c.baseClient.eventuallyWriteCache(resourceID, result)

	return result, nil
}

type RoomTypesSelector struct {
	EntitySelector
}

func (s RoomTypesSelector) resourceID() string {
	return s.EntitySelector.resourceID() + ".room-types"
}

func (c *entitiesClient) ReadRoomTypes(s RoomTypesSelector) ([]RoomType, error) {
	resourceID := s.resourceID()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(resourceID).([]RoomType); ok {
		return cachedResult, nil
	}

	references, err := transport.GetRESCollection[res.Ref](c.baseClient.resClient, resourceID, resprot.Request{})
	if err != nil {
		return nil, err
	}

	result := []RoomType{}

	for _, reference := range references {
		relatedRoomType, err := transport.GetRESModel[RoomType](c.baseClient.resClient, string(reference))
		if err != nil {
			return nil, err
		}

		result = append(result, relatedRoomType)
	}

	defer c.baseClient.eventuallyWriteCache(resourceID, result)

	return result, nil
}
