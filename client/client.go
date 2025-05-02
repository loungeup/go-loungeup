// Package client provides a client to interact with the LoungeUp API.
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
	Bookings      BookingsManager
	ComputedAttrs ComputedAttrsManager
	Entities      EntitiesManager
	Guests        GuestsManager
	Integrations  IntegrationsManager
	Products      ProductsManager
	ProxyDB       ProxyDBManager
	RoomTypes     RoomTypesManager
	Segments      SegmentsManager
	Currency      CurrencyManager
}

type BaseClient struct {
	httpAPIKey string
	httpClient transport.HTTPDoer
	httpAPIURL string
	resClient  transport.RESRequester
	cache      cache.ReadWriter
}

// Option used to configure a Client.
type Option func(*BaseClient)

// NewWithTransport returns a Client with the given transport and options.
func NewWithTransport(t *transport.Transport, c cache.ReadWriter, options ...Option) *Client {
	base := &BaseClient{
		httpClient: t.HTTPClient,
		resClient:  t.RESClient,
		cache:      c,
	}

	result := &Client{
		Bookings:      NewBookingsClient(base),
		ComputedAttrs: NewComputedAttrsClient(base),
		Entities:      NewEntitiesClient(base),
		Guests:        NewGuestsClient(base),
		Integrations:  NewIntegrationsClient(base),
		Products:      NewProductsClient(base),
		ProxyDB:       NewProxyDBClient(base),
		RoomTypes:     NewRoomTypesClient(base),
		Segments:      NewSegmentsClient(base),
		Currency:      NewCurrencyClient(base),
	}

	for _, option := range options {
		option(base)
	}

	return result
}

func WithHTTPAPIKey(key string) Option {
	return func(b *BaseClient) { b.httpAPIKey = key }
}

func WithHTTPAPIURL(url string) Option {
	return func(b *BaseClient) { b.httpAPIURL = url }
}

func (b *BaseClient) ExecuteHTTPRequest(request *http.Request) (*http.Response, error) {
	request.Header.Set("Authorization", "Bearer "+b.httpAPIKey)

	response, err := b.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %s", response.Status)
	}

	return response, nil
}

func (b *BaseClient) ReadCache(key string) any {
	if b.cache == nil {
		return nil
	}

	return b.cache.Read(key)
}

func (b *BaseClient) WriteCache(key string, value any) {
	if b.cache == nil {
		return
	}

	b.cache.Write(key, value)
}

func (b *BaseClient) WriteCacheWithDuration(key string, value any, duration time.Duration) {
	if b.cache == nil {
		return
	}

	b.cache.WriteWithDuration(key, value, duration)
}
