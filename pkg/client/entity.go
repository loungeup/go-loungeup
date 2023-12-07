package client

import (
	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

// Entity as represented with the RES protocol.
type Entity struct {
	ID         string         `json:"id"`
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
	response := c.resClient.Request("get."+authorityServiceName+".entities."+selector.ID.String(), resprot.Request{})
	if response.HasError() {
		return Entity{}, response.Error
	}

	result := Entity{}
	if _, err := response.ParseModel(&result); err != nil {
		return Entity{}, err
	}

	return result, nil
}
