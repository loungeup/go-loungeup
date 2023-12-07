package client

import "github.com/loungeup/go-loungeup/pkg/transport"

// Client is used to interact with our services.
type Client struct {
	Internal internalClient
}

// Option used to configure a Client.
type Option func(*Client)

// NewWithTransport returns a Client with the given transport and options.
func NewWithTransport(transport *transport.Transport) *Client {
	return &Client{
		Internal: internalClient{
			Entities: entitiesClient{
				resClient: transport.RESClient,
			},
			Guests: guestsClient{
				resClient: transport.RESClient,
			},
		},
	}
}
