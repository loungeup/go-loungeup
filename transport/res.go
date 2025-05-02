package transport

import (
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/log"
)

//go:generate mockgen -source res.go -destination=./mocks/mock_res.go -package=mocks

// Here, we are trying to use as much as possible from the resprot package. The resprot package provides functions and
// types to work with the (low-level) RES protocol.

type RESRequestHandler func(subject string, request resprot.Request) resprot.Response

// RESRequester is the interface used to execute a request using the RES protocol. It wraps the Request method.
type RESRequester interface {
	Request(subject string, request resprot.Request) resprot.Response
}

// RESClient used to interact with NATS services using the RES protocol.
type RESClient struct {
	disableRetries bool
	natsConnection res.Conn
	natsTimeout    time.Duration
}

type RESClientOption func(c *RESClient)

// NewRESClient returns a client to interact with NATS services using the RES protocol.
func NewRESClient(natsConnection res.Conn, options ...RESClientOption) *RESClient {
	const defaultNATSTimeout = 4 * time.Second

	result := &RESClient{
		natsConnection: natsConnection,
		natsTimeout:    defaultNATSTimeout,
	}

	for _, option := range options {
		option(result)
	}

	return result
}

func WithRESClientNATSTimeout(natsTimeout time.Duration) RESClientOption {
	return func(c *RESClient) { c.natsTimeout = natsTimeout }
}

func WithRESClientWithoutRetries() RESClientOption {
	return func(c *RESClient) { c.disableRetries = true }
}

var _ (RESRequester) = (*RESClient)(nil)

func (c *RESClient) Request(subject string, request resprot.Request) resprot.Response {
	if c.disableRetries {
		return c.requestOnce(subject, request)
	}

	return c.requestWithRetries(subject, request)
}

func (c *RESClient) requestWithRetries(subject string, request resprot.Request) resprot.Response {
	var lastResponse resprot.Response

	_ = backoff.RetryNotify(
		func() error {
			lastResponse = c.requestOnce(subject, request)
			if !lastResponse.HasError() {
				return nil
			}

			switch lastResponse.Error.Code {
			case res.CodeInternalError,
				res.CodeNotFound,
				res.CodeTimeout:
				return lastResponse.Error
			default:
				return backoff.Permanent(lastResponse.Error) // Do not retry for other codes.
			}
		},
		backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(c.natsTimeout)),
		func(err error, retryingIn time.Duration) {
			log.Default().Error("Could not request RES service. Retrying...",
				slog.Any("error", err),
				slog.String("subject", subject),
				slog.String("retryingIn", retryingIn.String()),
			)
		},
	)

	return lastResponse
}

func (c *RESClient) requestOnce(subject string, request resprot.Request) resprot.Response {
	return resprot.SendRequest(c.natsConnection, subject, request, c.natsTimeout)
}

// CallRESResult from the resource ID.
func CallRESResult[T any](client RESRequester, resourceID string, request resprot.Request) (T, error) {
	var result T

	response := client.Request("call."+resourceID, request)
	if response.HasError() {
		return result, response.Error
	}

	return result, response.ParseResult(&result)
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
func GetRESModel[T any](client RESRequester, resourceID string, request resprot.Request) (T, error) {
	var result T

	response := client.Request("get."+resourceID, request)

	if response.HasError() {
		return result, response.Error
	}

	if _, err := response.ParseModel(&result); err != nil {
		return result, err
	}

	return result, nil
}
