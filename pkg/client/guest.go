package client

import (
	"github.com/google/uuid"
	"github.com/jirenius/go-res/resprot"
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
