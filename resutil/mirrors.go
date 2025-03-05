package resutil

import (
	"log/slog"
	"reflect"
	"time"

	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/errors"
	"github.com/loungeup/go-loungeup/log"
	"github.com/loungeup/go-loungeup/pagination"
)

type MirrorModelProvider[Model, Selector any] interface {
	MakeModelRID(model Model) string
	ParseModelSelector(resource res.Resource) (Selector, error)
	ReadSourceModel(selector Selector) (Model, error)
}

func UseModelHandlerMirror[Model, Selector any](
	provider MirrorModelProvider[Model, Selector],
	next res.ModelHandler,
) res.ModelHandler {
	return func(request res.ModelRequest) {
		selector, err := provider.ParseModelSelector(request)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		source, err := provider.ReadSourceModel(selector)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		if isEmptyValue(source) {
			next(request)

			return
		}

		sourceResource, err := request.Service().Resource(provider.MakeModelRID(source))
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		next(&modelRequestMirror{
			Resource:    sourceResource,
			baseRequest: request,
		})
	}
}

type MirrorCallProvider[Model, Selector any] interface {
	MakeMirrorModelRIDsPager(selector Selector) *pagination.Pager[[]string, string]
	ParseModelSelector(resource res.Resource) (Selector, error)
	ReadSourceModel(selector Selector) (Model, error)
}

func UseCallHandlerMirror[Model, Selector any](
	provider MirrorCallProvider[Model, Selector],
	next res.CallHandler,
) res.CallHandler {
	return func(request res.CallRequest) {
		selector, err := provider.ParseModelSelector(request)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		source, err := provider.ReadSourceModel(selector)
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)

			return
		}

		if isEmptyValue(source) {
			next(&callRequestMirror{
				CallRequest:          request,
				mirrorModelRIDsPager: provider.MakeMirrorModelRIDsPager(selector),
			})

			return
		}

		errors.LogAndWriteRESError(log.Default(), request, &errors.Error{
			Code:    errors.CodeConflict,
			Message: "This resource is managed by another entity",
		})
	}
}

// callRequestMirror is an extension of res.CallRequest. It sends all RES events to mirrors. See:
// https://github.com/jirenius/go-res/blob/master/resource.go#L47-L102.
type callRequestMirror struct {
	res.CallRequest

	mirrorModelRIDsPager *pagination.Pager[[]string, string]
}

func (m *callRequestMirror) Event(event string, payload any) {
	m.CallRequest.Event(event, payload)
	m.applyToMirrorResources(func(resource res.Resource) { resource.Event(event, payload) })
}

func (m *callRequestMirror) ChangeEvent(props map[string]any) {
	m.CallRequest.ChangeEvent(props)
	m.applyToMirrorResources(func(resource res.Resource) { resource.ChangeEvent(props) })
}

func (m *callRequestMirror) AddEvent(value any, idx int) {
	m.CallRequest.AddEvent(value, idx)
	m.applyToMirrorResources(func(resource res.Resource) { resource.AddEvent(value, idx) })
}

func (m *callRequestMirror) RemoveEvent(idx int) {
	m.CallRequest.RemoveEvent(idx)
	m.applyToMirrorResources(func(resource res.Resource) { resource.RemoveEvent(idx) })
}

func (m *callRequestMirror) ReaccessEvent() {
	m.CallRequest.ReaccessEvent()
	m.applyToMirrorResources(func(resource res.Resource) { resource.ReaccessEvent() })
}

func (m *callRequestMirror) ResetEvent() {
	m.CallRequest.ResetEvent()
	m.applyToMirrorResources(func(resource res.Resource) { resource.ResetEvent() })
}

func (m *callRequestMirror) QueryEvent(f func(request res.QueryRequest)) {
	m.CallRequest.QueryEvent(f)
	m.applyToMirrorResources(func(resource res.Resource) { resource.QueryEvent(f) })
}

func (m *callRequestMirror) CreateEvent(value any) {
	m.CallRequest.CreateEvent(value)
	m.applyToMirrorResources(func(resource res.Resource) { resource.CreateEvent(value) })
}

func (m *callRequestMirror) DeleteEvent() {
	m.CallRequest.DeleteEvent()
	m.applyToMirrorResources(func(resource res.Resource) { resource.DeleteEvent() })
}

func (r *callRequestMirror) applyToMirrorResources(f func(resource res.Resource)) {
	for r.mirrorModelRIDsPager.Next() {
		for _, rid := range r.mirrorModelRIDsPager.Page() {
			if err := r.Service().With(rid, func(resource res.Resource) { f(resource) }); err != nil {
				log.Default().Error("Could not process mirror resource",
					slog.Any("error", err),
					slog.String("rid", rid),
				)
			}
		}
	}

	if err := r.mirrorModelRIDsPager.Err(); err != nil {
		log.Default().Error("Could not paginate mirror model RIDs", slog.Any("error", err))
	}
}

// modelRequestMirror is a res.ModelRequest. The goal is to be able to set a custom res.Resource and to use the base
// request for other methods. See: https://github.com/jirenius/go-res/blob/master/request.go#L65-L71.
type modelRequestMirror struct {
	res.Resource

	baseRequest res.ModelRequest
}

var _ (res.ModelRequest) = (*modelRequestMirror)(nil)

func (m *modelRequestMirror) Model(model any) { m.baseRequest.Model(model) }

func (m *modelRequestMirror) QueryModel(model any, query string) {
	m.baseRequest.QueryModel(model, query)
}

func (m *modelRequestMirror) NotFound() { m.baseRequest.NotFound() }

func (m *modelRequestMirror) InvalidQuery(message string) { m.baseRequest.InvalidQuery(message) }

func (m *modelRequestMirror) Error(err error) { m.baseRequest.Error(err) }

func (m *modelRequestMirror) Timeout(d time.Duration) { m.baseRequest.Timeout(d) }

func (m *modelRequestMirror) ForValue() bool { return m.baseRequest.ForValue() }

func isEmptyValue[T any](value T) bool {
	var emptyValue T

	return reflect.DeepEqual(value, emptyValue)
}
