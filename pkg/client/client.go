package client

import (
	"github.com/loungeup/go-loungeup/pkg/cache"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

// Client is used to interact with our services.
type Client struct {
	Internal internalClient

	// The following fields are used internally by sub-clients.
	cache     cache.ReadWriter
	resClient transport.RESRequester
}

// Option used to configure a Client.
type Option func(*Client)

// NewWithTransport returns a Client with the given transport and options.
func NewWithTransport(transport *transport.Transport, options ...Option) *Client {
	result := &Client{
		resClient: transport.RESClient,
	}

	result.Internal = internalClient{
		Entities:     &entitiesClient{baseClient: result},
		Guests:       &guestsClient{baseClient: result},
		Integrations: &integrationsClient{baseClient: result},
		ProxyDB:      &proxyDBClient{baseClient: result},
		RoomTypes:    &roomTypesClient{baseClient: result},
	}

	for _, option := range options {
		option(result)
	}

	return result
}

func WithCache(cache cache.ReadWriter) Option { return func(c *Client) { c.cache = cache } }

func (c *Client) eventuallyReadCache(key string) any {
	if c.cache == nil {
		return nil
	}

	return c.cache.Read(key)
}

func (c *Client) eventuallyWriteCache(key string, value any) {
	if c.cache == nil {
		return
	}

	c.cache.Write(key, value)
}
