package transport

import (
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
)

// Here, we are trying to use as much as possible from the resprot package. The resprot package provides functions and
// types to work with the (low-level) RES protocol.

// RESRequester is the interface used to execute a request using the RES protocol. It wraps the Request method.
type RESRequester interface {
	Request(resourceID string, request resprot.Request) resprot.Response
}

// RESClient used to interact with NATS services using the RES protocol.
type RESClient struct {
	natsConnection res.Conn
}

// NewRESClient returns a client to interact with NATS services using the RES protocol.
func NewRESClient(natsConnection res.Conn) *RESClient {
	return &RESClient{
		natsConnection: natsConnection,
	}
}

var _ (RESRequester) = (*RESClient)(nil)

func (c *RESClient) Request(resourceID string, request resprot.Request) resprot.Response {
	return resprot.SendRequest(c.natsConnection, resourceID, request, defaultNATSTimeout)
}
