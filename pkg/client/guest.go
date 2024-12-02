package client

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/transport"
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
	rid := selector.rid()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(rid).(*models.Guest); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*models.Guest](c.baseClient.resClient, rid, resprot.Request{})
	if err != nil {
		return nil, err
	}

	defer c.baseClient.eventuallyWriteCache(rid, result)

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

func (s *GuestSelector) rid() string {
	query := url.Values{}

	if s.Expand {
		query.Add("expand", "true")
	}

	if s.Redirect {
		query.Add("redirect", "true")
	}

	result := "guestprofile.entities." + s.EntityID.String() + ".guests." + s.GuestID.String()
	if len(query) > 0 {
		result += "?" + query.Encode()
	}

	return result
}
