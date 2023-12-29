package client

import (
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

type EntitySelector struct {
	ID uuid.UUID `json:"id"`
}

// entitiesClient provides methods to interact with entities.
type entitiesClient struct {
	resClient transport.RESRequester
}

func (c *entitiesClient) ReadEntity(selector EntitySelector) (Entity, error) {
	return transport.GetRESModel[Entity](c.resClient, authorityServiceName+".entities."+selector.ID.String())
}

func (c *entitiesClient) ReadEntityAccounts(selector EntitySelector) ([]Entity, error) {
	resourceID := authorityServiceName + ".entities." + selector.ID.String() + ".accounts"

	accountReferences, err := transport.GetRESCollection[res.Ref](c.resClient, resourceID)
	if err != nil {
		return nil, err
	}

	result := []Entity{}

	for _, accountReference := range accountReferences {
		relatedAccount, err := transport.GetRESModel[Entity](c.resClient, string(accountReference))
		if err != nil {
			return nil, err
		}

		result = append(result, relatedAccount)
	}

	return result, nil
}
