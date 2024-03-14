package client

import (
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/cache"
	"github.com/loungeup/go-loungeup/pkg/transport"
	"github.com/loungeup/go-loungeup/pkg/transporttest"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	clientToTest := NewWithTransport(
		transport.New(
			transport.WithRESClient(&transporttest.RESClientMock{
				RequestFunc: func(resourceID string, request resprot.Request) resprot.Response {
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
