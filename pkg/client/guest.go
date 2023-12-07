package client

import (
	"github.com/google/uuid"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

type guestsClient struct {
	resClient transport.RESRequester
}

func (c *guestsClient) AnonymizeGuests(entityID uuid.UUID, guestIDs []uuid.UUID) error {
	response := c.resClient.Request("call.guestprofile.entities."+entityID.String()+".guests.anonymize", resprot.Request{
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
