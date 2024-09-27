package client

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
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

func (c *guestsClient) ReadIndexableGuestSelectors(
	selector *models.IndexableGuestsSelector,
) ([]*res.DataValue[*models.GuestSelector], error) {
	return transport.GetRESCollection[*res.DataValue[*models.GuestSelector]](
		c.baseClient.resClient,
		selector.RID(),
		resprot.Request{
			Query: selector.EncodedQuery(),
		},
	)
}

func (c *guestsClient) ReadGuest(selector *models.GuestSelector) (*models.Guest, error) {
	return c.readGuestByRID(selector.RID())
}

func (c *guestsClient) readGuestByRID(rid string) (*models.Guest, error) {
	if cachedResult, ok := c.baseClient.eventuallyReadCache(rid).(*models.Guest); ok {
		return cachedResult, nil
	}

	guest, err := transport.GetRESModel[*models.Guest](c.baseClient.resClient, rid, resprot.Request{})
	if err != nil {
		return nil, err
	}

	defer c.baseClient.eventuallyWriteCache(rid, guest)

	return guest, nil
}
