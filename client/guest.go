package client

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/resresultsets"
	"github.com/loungeup/go-loungeup/transport"
)

type guestsClient struct{ baseClient *Client }

func (c *guestsClient) AnonymizeGuests(entityID uuid.UUID, guestIDs []uuid.UUID) error {
	resourceID := "call.guestprofile.entities." + entityID.String() + ".guests.anonymize"

	response := c.baseClient.resClient.Request(resourceID, resprot.Request{
		Params: map[string]any{
			"guests": guestIDs,
		},
		Token: map[string]any{
			"agentId":    "",
			"agentRoles": []string{"service"},
		},
	})
	if response.HasError() {
		return response.Error
	}

	return nil
}

func (c *guestsClient) CountGuests(
	entityID uuid.UUID,
	request *models.SearchGuestsRequest,
) (*models.CountGuestsResponse, error) {
	result, err := transport.CallRESResult[*models.CountGuestsResponse](
		c.baseClient.resClient,
		"guestprofile.entities."+entityID.String()+".guests.count",
		resprot.Request{Params: request},
	)
	if err != nil {
		return nil, fmt.Errorf("could not execute request: %w", err)
	}

	return result, nil
}

func (c *guestsClient) ReadOne(selector *GuestSelector) (*models.Guest, error) {
	cacheKey := selector.makeCacheKey()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(cacheKey).(*models.Guest); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*models.Guest](c.baseClient.resClient, selector.makeRID(), resprot.Request{
		Query: selector.makeEncodedQuery(),
	})
	if err != nil {
		return nil, err
	}

	defer c.baseClient.eventuallyWriteCache(cacheKey, result)

	return result, nil
}

type GuestSelector struct {
	GuestID  uuid.UUID
	EntityID uuid.UUID

	// Expand the composition of the guest.
	Expand bool

	// Redirect to the root guest.
	Redirect bool
}

func (selector *GuestSelector) makeCacheKey() string {
	return selector.makeRID() + selector.makeEncodedQuery()
}

func (selector *GuestSelector) makeRID() string {
	return "guestprofile.entities." + selector.EntityID.String() + ".guests." + selector.GuestID.String()
}

func (selector *GuestSelector) makeEncodedQuery() string {
	query := url.Values{}

	if selector.Expand {
		query.Add("expand", "true")
	}

	if selector.Redirect {
		query.Add("redirect", "true")
	}

	return query.Encode()
}

func (c *guestsClient) SearchByContact(
	selector *models.SearchByContactSelector,
) (*resresultsets.KeysetPaginationModel, error) {
	return transport.CallRESResult[*resresultsets.KeysetPaginationModel](
		c.baseClient.resClient,
		"guestprofile.entities."+selector.EntityID.String()+".guests.search-by-contact",
		resprot.Request{Params: selector},
	)
}
