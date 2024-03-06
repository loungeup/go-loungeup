package client

import (
	"encoding/json"
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/cache"
	"github.com/loungeup/go-loungeup/pkg/transport"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	clientToTest := NewWithTransport(
		transport.New(
			transport.WithRESClient(&resClientMock{
				requestFunc: func(resourceID string, request resprot.Request) resprot.Response {
					return resprot.Response{}
				},
			}),
		),
		WithCache(&cache.Mock{}),
	)

	require.NotNil(t, clientToTest.cache)
	require.NotNil(t, clientToTest.resClient)

	// We are not testing each sub-client here because they are just encoding the request, using the transport layer and
	// parsing the response.
}

type resClientMock struct {
	requestFunc func(resourceID string, request resprot.Request) resprot.Response
}

var _ (transport.RESRequester) = (*resClientMock)(nil)

func (c *resClientMock) Request(resourceID string, request resprot.Request) resprot.Response {
	return c.requestFunc(resourceID, request)
}

// newCollectionResponse creates a new response with the specified collection.
func newCollectionResponse(collection string) resprot.Response {
	return resprot.Response{Result: json.RawMessage(`{"collection":` + collection + `}`)}
}

// newModelResponse creates a new resprot.Response with the given model.
func newModelResponse(model string) resprot.Response {
	return resprot.Response{Result: json.RawMessage(`{"model":` + model + `}`)}
}

// newResultResponse creates a new resprot.Response with the given result.
func newResultResponse(result string) resprot.Response {
	return resprot.Response{Result: json.RawMessage(result)}
}
