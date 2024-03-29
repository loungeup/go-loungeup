// Package resutil provides RES utility functions, complementing the more common ones in the github.com/jirenius/go-res
// package.
package resutil

import (
	"github.com/jirenius/go-res"
	"github.com/nats-io/nats.go"
)

const natsMessagesChannelSize = 64

// AddNATSMessageHandler to a RES service. The handler will be called for each message received on the given subject.
func AddNATSMessageHandler(
	service *res.Service,
	subject string,
	handler func(message *nats.Msg) error,
) error {
	messages := make(chan *nats.Msg, natsMessagesChannelSize)
	if _, err := service.Conn().ChanSubscribe(subject, messages); err != nil {
		return err
	}

	go func() {
		for message := range messages {
			if logger := service.Logger(); logger != nil {
				logger.Tracef("S--> %s: %s", message.Subject, string(message.Data))
			}

			if err := handler(message); err != nil {
				if logger := service.Logger(); logger != nil {
					logger.Errorf("Could not handle NATS message: %s", err)
				}
			}
		}
	}()

	return nil
}

func HandleCollectionQueryRequest[T any](
	service *res.Service,
	rid string,
	handler func(request res.QueryRequest) ([]T, error),
) error {
	return handleQueryRequest(service, rid, func(request res.QueryRequest) {
		response, err := handler(request)
		if err != nil {
			request.Error(err)
			return
		}

		request.Collection(response)
	})
}

func HandleModelQueryRequest[T any](
	service *res.Service,
	rid string,
	handler func(request res.QueryRequest) (T, error),
) error {
	return handleQueryRequest(service, rid, func(request res.QueryRequest) {
		response, err := handler(request)
		if err != nil {
			request.Error(err)
			return
		}

		request.Model(response)
	})
}

// https://resgate.io/docs/specification/res-service-protocol/#query-request
func handleQueryRequest(
	service *res.Service,
	rid string,
	handler func(request res.QueryRequest),
) error {
	return service.With(rid, func(resource res.Resource) {
		resource.QueryEvent(func(request res.QueryRequest) {
			if request == nil {
				return // https://github.com/jirenius/go-res/blob/372a82d603a13d7601f8b14e74eccaebd325ee61/resource.go#L336-L339
			}

			handler(request)
		})
	})
}

// MapRefs from a slice of elements.
// The function f is called for each element in the slice, and the resulting reference is added to the result.
func MapRefs[E any, S []E](s S, f func(e E) res.Ref) []res.Ref {
	result := []res.Ref{}

	for _, e := range s {
		ref := f(e)
		if !ref.IsValid() {
			continue
		}

		result = append(result, ref)
	}

	return result
}
