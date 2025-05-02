package resutil

import (
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/errors"
	"github.com/loungeup/go-loungeup/log"
	"github.com/loungeup/go-loungeup/resmodels"
)

// EntityCallRequest wraps a res.CallRequest with an associated Entity. This structure allows handlers to access the
// entity directly, reducing the need for repeated entity lookups and improving code clarity.
type EntityCallRequest struct {
	res.CallRequest

	Entity *resmodels.Entity
}

type EntityCallHandler func(request *EntityCallRequest)

type EntityCallRequestConfig struct {
	// PathParamKey used to extract the entity ID from the request path parameters.
	PathParamKey string
}

type EntityCallRequestOption func(config *EntityCallRequestConfig)

type EntityReader interface {
	ReadEntity(selector *resmodels.EntitySelector) (*resmodels.Entity, error)
}

func WithEntityCallRequestPathParamKey(key string) EntityCallRequestOption {
	return func(config *EntityCallRequestConfig) { config.PathParamKey = key }
}

// WithEntityCallHandler wraps a handler to include entity resolution logic. This function ensures that the entity is
// resolved and passed to the handler, simplifying the handler's responsibilities and centralizing error handling.
func WithEntityCallHandler(next EntityCallHandler, reader EntityReader, options ...EntityCallRequestOption) res.CallHandler {
	config := &EntityCallRequestConfig{
		PathParamKey: "entityID",
	}

	for _, option := range options {
		option(config)
	}

	return func(request res.CallRequest) {
		entityID, err := ParseUUIDPathParam(request, config.PathParamKey)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		entity, err := reader.ReadEntity(&resmodels.EntitySelector{EntityID: entityID})
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		next(&EntityCallRequest{
			CallRequest: request,
			Entity:      entity,
		})
	}
}
