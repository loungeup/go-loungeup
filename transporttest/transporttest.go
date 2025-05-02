// Package transporttest provides utilities to test the transport package.
package transporttest

import (
	"encoding/json"
	"net/http"

	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/transport"
)

type HTTPClientMock struct {
	DoFunc func(request *http.Request) (*http.Response, error)
}

var _ transport.HTTPDoer = (*HTTPClientMock)(nil)

func (c *HTTPClientMock) Do(request *http.Request) (*http.Response, error) {
	return c.DoFunc(request)
}

type RESClientMock struct {
	RequestFunc transport.RESRequestHandler
}

var _ transport.RESRequester = (*RESClientMock)(nil)

func (c *RESClientMock) Request(resourceID string, request resprot.Request) resprot.Response {
	return c.RequestFunc(resourceID, request)
}

// NewRESCollectionResponse creates a new response with the specified collection.
func NewRESCollectionResponse(collection string) resprot.Response {
	return resprot.Response{Result: json.RawMessage(`{"collection":` + collection + `}`)}
}

// NewRESInternalErrorResponse creates a new response with an internal error.
func NewRESInternalErrorResponse() resprot.Response {
	return resprot.Response{Error: res.ErrInternalError}
}

// NewRESModelResponse creates a new resprot.Response with the given model.
func NewRESModelResponse(model string) resprot.Response {
	return resprot.Response{Result: json.RawMessage(`{"model":` + model + `}`)}
}

// NewRESResultResponse creates a new resprot.Response with the given result.
func NewRESResultResponse(result string) resprot.Response {
	return resprot.Response{Result: json.RawMessage(result)}
}

func UseHandlers(handlersPerSubject map[string]transport.RESRequestHandler) transport.RESRequestHandler {
	return func(subject string, request resprot.Request) resprot.Response {
		if handler, ok := handlersPerSubject[subject]; ok {
			return handler(subject, request)
		}

		panic("unexpected subject: '" + subject + "'")
	}
}

func UseJSONModelHandler[T any](model T) transport.RESRequestHandler {
	return func(_ string, _ resprot.Request) resprot.Response {
		encodedValue, _ := json.Marshal(model)

		return NewRESModelResponse(string(encodedValue))
	}
}
