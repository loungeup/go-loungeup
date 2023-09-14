// client package provides a client for the LoungeUp API.
package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res/resprot"
	"github.com/nats-io/nats.go"
)

const defaultNATSRequestTimeout = 60 * time.Second

type Client struct {
	GuestProfile *guestProfileClient
}

type Option func(*Client)

func New(options ...Option) *Client {
	result := &Client{
		GuestProfile: &guestProfileClient{},
	}

	for _, option := range options {
		option(result)
	}

	return result
}

func WithNatsConnection(natsConnection *nats.Conn) Option {
	return func(c *Client) {
		c.GuestProfile.natsConnection = natsConnection
	}
}

type guestProfileClient struct {
	natsConnection *nats.Conn
}

func (c *guestProfileClient) AnonymizeGuests(relatedEntity uuid.UUID, guestsToAnonymize []uuid.UUID) error {
	encodedRequest, err := json.Marshal(newRESRequestWithParams(map[string]any{
		"guests": guestsToAnonymize,
	}))
	if err != nil {
		return err
	}

	receivedMessage, err := c.natsConnection.Request(
		"call.guestprofile.entities."+relatedEntity.String()+".guests.anonymize",
		encodedRequest,
		defaultNATSRequestTimeout,
	)
	if err != nil {
		return err
	}

	return parseRESResponse(receivedMessage.Data)
}

func newRESRequestWithParams(paramsToSend any) *resprot.Request {
	return &resprot.Request{
		Params: paramsToSend,
		Token: map[string]any{
			"agentId":    "",
			"agentRoles": []string{"service"},
		},
	}
}

func parseRESResponse(rawResponse json.RawMessage) error {
	parsedResponse := resprot.ParseResponse(rawResponse)
	if parsedResponse.HasError() {
		return parsedResponse.Error
	}

	return nil
}
