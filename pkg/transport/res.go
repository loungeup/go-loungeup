package transport

import (
	"encoding/json"

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

// CallRESResult from the resource ID.
func CallRESResult(client RESRequester, resourceID string, request resprot.Request) (json.RawMessage, error) {
	result := json.RawMessage(`{}`)

	response := client.Request("call."+resourceID, request)
	if response.HasError() {
		return result, response.Error
	}

	return response.Result, nil
}

// GetRESCollection from the resource ID.
func GetRESCollection[T any](client RESRequester, resourceID string, request resprot.Request) ([]T, error) {
	var result []T

	response := client.Request("get."+resourceID, request)
	if response.HasError() {
		return result, response.Error
	}

	if _, err := response.ParseCollection(&result); err != nil {
		return result, err
	}

	return result, nil
}

// GetRESModel from the resource ID.
func GetRESModel[T any](client RESRequester, resourceID string) (T, error) {
	var result T

	response := client.Request("get."+resourceID, resprot.Request{})

	if response.HasError() {
		return result, response.Error
	}

	if _, err := response.ParseModel(&result); err != nil {
		return result, err
	}

	return result, nil
}
