package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/loungeup/go-loungeup/cache"
	"github.com/loungeup/go-loungeup/transport"
)

// Client is used to interact with our services.
type Client struct {
	Internal *internalClient

	// The following fields are used internally by sub-clients.
	cache      cache.ReadWriter
	httpAPIKey string
	httpClient transport.HTTPDoer
	httpAPIURL string
	resClient  transport.RESRequester
}

// Option used to configure a Client.
type Option func(*Client)

// NewWithTransport returns a Client with the given transport and options.
func NewWithTransport(transport *transport.Transport, options ...Option) *Client {
	result := &Client{
		httpClient: transport.HTTPClient,
		resClient:  transport.RESClient,
	}

	result.Internal = &internalClient{
		Bookings:      &bookingsClient{baseClient: result},
		ComputedAttrs: &computedAttrsClient{baseClient: result},
		Currency:      &currencyClient{baseClient: result},
		Entities:      &entitiesClient{baseClient: result},
		Guests:        &guestsClient{baseClient: result},
		Integrations:  &integrationsClient{baseClient: result},
		Products:      &productsClient{baseClient: result},
		ProxyDB:       &proxyDBClient{baseClient: result},
		RoomTypes:     &roomTypesClient{baseClient: result},
		Segments:      &segmentsClient{baseClient: result},
	}

	for _, option := range options {
		option(result)
	}

	return result
}

func WithCache(cache cache.ReadWriter) Option { return func(c *Client) { c.cache = cache } }
func WithHTTPAPIKey(key string) Option        { return func(c *Client) { c.httpAPIKey = key } }
func WithHTTPAPIURL(url string) Option        { return func(c *Client) { c.httpAPIURL = url } }

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

func (c *Client) eventuallyWriteCacheWithDuration(key string, value any, duration time.Duration) {
	if c.cache == nil {
		return
	}

	c.cache.WriteWithDuration(key, value, duration)
}

func (c *Client) executeHTTPRequest(request *http.Request) (*http.Response, error) {
	request.Header.Set("Authorization", "Bearer "+c.httpAPIKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %s", response.Status)
	}

	return response, nil
}
