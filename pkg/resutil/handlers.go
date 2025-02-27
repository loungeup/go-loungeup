package resutil

import (
	"log/slog"
	"net/url"

	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/pkg/errors"
	"github.com/loungeup/go-loungeup/pkg/log"
)

type (
	GetCollectionProvider[Collection ~[]Model, Model, Selector any] interface {
		MakeCollectionQuery(selector Selector) url.Values
		ParseCollectionSelector(resource res.Resource) (Selector, error)
		ReadCollection(selector Selector) (Collection, error)
	}

	GetModelProvider[Model, Selector any] interface {
		MakeModelQuery(selector Selector) url.Values
		ParseModelSelector(resource res.Resource) (Selector, error)
		ReadModel(selector Selector) (Model, error)
	}

	CreateModelProvider[Model, Params any] interface {
		CreateModel(params Params) (Model, error)
		MakeModelRID(model Model) string
		ParseModelParams(request res.CallRequest) (Params, error)
	}

	UpdateModelProvider[Model, Selector, Params any] interface {
		MakeModelRID(model Model) string
		ParseModelParams(request res.CallRequest) (Params, error)
		ParseModelSelector(resource res.Resource) (Selector, error)
		ReadModel(selector Selector) (Model, error)
		UpdateModel(existingModel Model, params Params) (Model, error)
	}

	DeleteModelProvider[Selector any] interface {
		ParseModelSelector(resource res.Resource) (Selector, error)
		DeleteModel(selector Selector) error
	}
)

func UseGetCollectionHandler[Collection ~[]Model, Model, Selector any](
	provider GetCollectionProvider[Collection, Model, Selector],
) res.CollectionHandler {
	return func(request res.CollectionRequest) {
		selector, err := provider.ParseCollectionSelector(request)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		collection, err := provider.ReadCollection(selector)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		request.QueryCollection(
			eventuallyMapCollection(provider, collection),
			provider.MakeCollectionQuery(selector).Encode(),
		)
	}
}

func UseGetModelHandler[Model, Selector any](provider GetModelProvider[Model, Selector]) res.ModelHandler {
	return func(request res.ModelRequest) {
		selector, err := provider.ParseModelSelector(request)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		model, err := provider.ReadModel(selector)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		request.QueryModel(
			eventuallyMapModel(provider, model),
			provider.MakeModelQuery(selector).Encode(),
		)
	}
}

func UseCreateModelHandler[Model, Params any](provider CreateModelProvider[Model, Params]) res.CallHandler {
	return func(request res.CallRequest) {
		params, err := provider.ParseModelParams(request)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		model, err := provider.CreateModel(params)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		request.Resource(provider.MakeModelRID(model))
	}
}

func UseUpdateModelHandler[Model, Selector, Updates any](
	provider UpdateModelProvider[Model, Selector, Updates],
) res.CallHandler {
	return func(request res.CallRequest) {
		params, err := provider.ParseModelParams(request)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		selector, err := provider.ParseModelSelector(request)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		existingModel, err := provider.ReadModel(selector)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		updatedModel, err := provider.UpdateModel(existingModel, params)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		request.Resource(provider.MakeModelRID(updatedModel))
		request.ChangeEvent(CompareModels(
			eventuallyMapModel(provider, existingModel),
			eventuallyMapModel(provider, updatedModel),
		))
	}
}

func UseDeleteModelHandler[Selector any](provider DeleteModelProvider[Selector]) res.CallHandler {
	return func(request res.CallRequest) {
		selector, err := provider.ParseModelSelector(request)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		if err := provider.DeleteModel(selector); err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		request.OK(nil)
		request.DeleteEvent()
	}
}

func WithCallHandlerHooks(handler res.CallHandler, hooks []res.CallHandler) res.CallHandler {
	return func(request res.CallRequest) {
		handler(request)

		for _, hook := range hooks {
			hook(request)
		}
	}
}

func WithCollectionQueryEventHandler[Collection ~[]Model, Model, Selector any](
	makeRIDFunc func(resource res.Resource) string,
	provider GetCollectionProvider[Collection, Model, Selector],
) res.CallHandler {
	return func(request res.CallRequest) {
		rid := makeRIDFunc(request)

		if err := handleQueryRequest(request.Service(), rid, func(request res.QueryRequest) {
			selector, err := provider.ParseCollectionSelector(request)
			if err != nil {
				errors.LogAndWriteRESError(log.Default(), request, err)

				return
			}

			result, err := provider.ReadCollection(selector)
			if err != nil {
				errors.LogAndWriteRESError(log.Default(), request, err)

				return
			}

			request.Model(eventuallyMapCollection(provider, result))
		}); err != nil {
			log.Default().Error("Could not handle collection query request",
				slog.Any("error", err),
				slog.String("rid", rid),
			)
		}
	}
}

func WithModelQueryEventHandler[Model, Selector any](
	makeRIDFunc func(resource res.Resource) string,
	provider GetModelProvider[Model, Selector],
) res.CallHandler {
	return func(request res.CallRequest) {
		rid := makeRIDFunc(request)

		if err := handleQueryRequest(request.Service(), rid, func(request res.QueryRequest) {
			selector, err := provider.ParseModelSelector(request)
			if err != nil {
				errors.LogAndWriteRESError(log.Default(), request, err)

				return
			}

			result, err := provider.ReadModel(selector)
			if err != nil {
				errors.LogAndWriteRESError(log.Default(), request, err)

				return
			}

			request.Model(eventuallyMapModel(provider, result))
		}); err != nil {
			log.Default().Error("Could not handle model query request",
				slog.Any("error", err),
				slog.String("rid", rid),
			)
		}
	}
}

type (
	collectionMapper[Collection ~[]Model, Model any] interface {
		MapCollection(collection Collection) any
	}

	modelMapper[Model any] interface {
		MapModel(model Model) any
	}
)

func eventuallyMapCollection[Collection ~[]Model, Model any](provider any, collection Collection) any {
	if mapper, ok := provider.(collectionMapper[Collection, Model]); ok {
		return mapper.MapCollection(collection)
	}

	return collection
}

func eventuallyMapModel[Model any](provider any, model Model) any {
	if mapper, ok := provider.(modelMapper[Model]); ok {
		return mapper.MapModel(model)
	}

	return model
}
