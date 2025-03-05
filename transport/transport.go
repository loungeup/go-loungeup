package transport

import "net/http"

// Transport used to interact with the LoungeUp servers.
type Transport struct {
	// HTTPClient used to interact with HTTP services.
	HTTPClient HTTPDoer

	// RESClient used to interact with NATS services using the RES protocol.
	RESClient RESRequester
}

// Option used to configure a Transport.
type Option func(*Transport)

// New returns a Transport with the given options. By default, the HTTP transport is set to http.DefaultClient.
func New(options ...Option) *Transport {
	result := &Transport{
		HTTPClient: http.DefaultClient,
	}

	for _, option := range options {
		option(result)
	}

	return result
}

// WithHTTPClient is an option to set the HTTP client of a Transport.
func WithHTTPClient(httpClient HTTPDoer) Option {
	return func(t *Transport) { t.HTTPClient = httpClient }
}

// WithRESClient is an option to set the RES client of a Transport.
func WithRESClient(resClient RESRequester) Option {
	return func(t *Transport) { t.RESClient = resClient }
}
