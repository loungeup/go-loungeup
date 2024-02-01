package client

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

// Entity as represented with the RES protocol.
type Entity struct {
	ID         uuid.UUID      `json:"id"`
	LegacyID   int            `json:"legacyId,omitempty"`
	Type       string         `json:"type"`
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

	Limit, Offset uint
}

func (s EntityAccountsSelector) resourceID() string {
	result := s.EntitySelector.resourceID() + ".accounts"
	if s.Limit > 0 {
		result += "?limit=" + fmt.Sprint(s.Limit)
	}

	if s.Offset > 0 {
		result += "&offset=" + fmt.Sprint(s.Offset)
	}

	return result
}

func (c *entitiesClient) ReadEntityAccounts(s EntityAccountsSelector) ([]Entity, error) {
	resourceID := s.resourceID()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(resourceID).([]Entity); ok {
		return cachedResult, nil
	}

	accountReferences, err := transport.GetRESCollection[res.Ref](c.baseClient.resClient, resourceID)
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

	defer c.baseClient.eventuallyWriteCache(resourceID, result)

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
