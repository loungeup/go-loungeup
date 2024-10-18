// Package resutil provides RES utility functions, complementing the more common ones in the github.com/jirenius/go-res
// package.
package resutil

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"reflect"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/pkg/errors"
	"github.com/loungeup/go-loungeup/pkg/log"
	"github.com/nats-io/nats.go"
)

const natsMessagesChannelSize = 64

type Deletable[T any] struct {
	Deleted bool
	Value   T
}

func (d *Deletable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, res.DeleteAction.RawMessage) {
		d.Deleted = true

		return nil
	}

	return json.Unmarshal(data, &d.Value)
}

type RefSlice []res.Ref

func (s RefSlice) Strings() []string {
	result := []string{}
	for _, ref := range s {
		result = append(result, string(ref))
	}

	return result
}

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

// CompareModels and returns a map of the differences between them.
// The result can be used with the https://pkg.go.dev/github.com/jirenius/go-res#Request.ChangeEvent method.
func CompareModels[T any](previous, current T) map[string]any {
	mapFromPrevious, err := convertToMap(previous)
	if err != nil {
		log.Default().Error("Could not convert previous model to a map",
			slog.String("error", err.Error()),
			slog.Any("previous", previous),
		)

		return nil
	}

	mapFromCurrent, err := convertToMap(current)
	if err != nil {
		log.Default().Error("Could not convert current model to a map",
			slog.Any("current", current),
			slog.String("error", err.Error()),
		)

		return nil
	}

	return compareMaps(mapFromPrevious, mapFromCurrent)
}

// convertToMap converts a value to a map[string]any.
func convertToMap[T any](value T) (map[string]any, error) {
	encodedValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	result := map[string]any{}
	if err := json.Unmarshal(encodedValue, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// compareMaps and returns a map of the differences between them.
// The result can be used with the https://pkg.go.dev/github.com/jirenius/go-res#Request.ChangeEvent method.
func compareMaps(previous, current map[string]any) map[string]any {
	result := map[string]any{}

	for key := range previous {
		if _, found := current[key]; !found {
			result[key] = res.DeleteAction
		}
	}

	for key, value := range current {
		previousValue, found := previous[key]
		if !found {
			result[key] = value
			continue
		}

		if reflect.DeepEqual(value, previousValue) {
			continue
		}

		result[key] = value
	}

	return result
}

func HandleCollectionQueryRequest[Collection ~[]Model, Model any](
	service *res.Service,
	rid string,
	handler func(request res.QueryRequest) (Collection, error),
) error {
	return handleQueryRequest(service, rid, func(request res.QueryRequest) {
		response, err := handler(request)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)
			return
		}

		request.Collection(response)
	})
}

func HandleModelQueryRequest[Model any](
	service *res.Service,
	rid string,
	handler func(request res.QueryRequest) (Model, error),
) error {
	return handleQueryRequest(service, rid, func(request res.QueryRequest) {
		response, err := handler(request)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)
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
func MapRefs[S ~[]E, E any](s S, f func(e E) res.Ref) RefSlice {
	result := RefSlice{}

	for _, e := range s {
		ref := f(e)
		if !ref.IsValid() {
			continue
		}

		result = append(result, ref)
	}

	return result
}

// ParseUUIDPathParam from the resource with the given key.
func ParseUUIDPathParam(resource res.Resource, key string) (uuid.UUID, error) {
	result, err := uuid.Parse(resource.PathParam(key))
	if err != nil {
		return uuid.Nil, &errors.Error{
			Code:            errors.CodeInvalid,
			Message:         "Invalid '" + key + "' path parameter",
			UnderlyingError: err,
		}
	}

	return result, nil
}
