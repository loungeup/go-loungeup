package client

import (
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/transport"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	clientToTest := NewWithTransport(transport.New(
		transport.WithRESClient(&resClientMock{
			requestFunc: func(resourceID string, request resprot.Request) resprot.Response {
				return resprot.Response{}
			},
		}),
	))

	require.NotNil(t, clientToTest.Internal.Entities.resClient)
	require.NotNil(t, clientToTest.Internal.Guests.resClient)

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
