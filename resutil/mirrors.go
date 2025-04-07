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

const errorCodeManagedByAnotherEntity = "loungeup.resourceManagedByAnotherEntity"

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

type MirrorCallProvider[Model, Selector any, PageReader pagination.PageReader[[]string, string]] interface {
	MakeMirrorModelRIDsPager(selector Selector) *pagination.Pager[[]string, string, PageReader]
	ParseModelSelector(resource res.Resource) (Selector, error)
	ReadSourceModel(selector Selector) (Model, error)
}

func UseCallHandlerMirror[Model, Selector any, PageReader pagination.PageReader[[]string, string]](
	provider MirrorCallProvider[Model, Selector, PageReader],
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
			next(&callRequestMirror[PageReader]{
				CallRequest:          request,
				mirrorModelRIDsPager: provider.MakeMirrorModelRIDsPager(selector),
			})

			return
		}

		errors.LogAndWriteRESError(log.Default(), request, &errors.Error{
			Code:    errorCodeManagedByAnotherEntity,
			Message: "This resource is managed by another entity",
		})
	}
}

// callRequestMirror is an extension of res.CallRequest. It sends all RES events to mirrors. See:
// https://github.com/jirenius/go-res/blob/master/resource.go#L47-L102.
type callRequestMirror[PageReader pagination.PageReader[[]string, string]] struct {
	res.CallRequest

	mirrorModelRIDsPager *pagination.Pager[[]string, string, PageReader]
}

func (m *callRequestMirror[PageReader]) Event(event string, payload any) {
	m.CallRequest.Event(event, payload)
	m.applyToMirrorResources(func(resource res.Resource) { resource.Event(event, payload) })
}

func (m *callRequestMirror[PageReader]) ChangeEvent(props map[string]any) {
	m.CallRequest.ChangeEvent(props)
	m.applyToMirrorResources(func(resource res.Resource) { resource.ChangeEvent(props) })
}

func (m *callRequestMirror[PageReader]) AddEvent(value any, idx int) {
	m.CallRequest.AddEvent(value, idx)
	m.applyToMirrorResources(func(resource res.Resource) { resource.AddEvent(value, idx) })
}

func (m *callRequestMirror[PageReader]) RemoveEvent(idx int) {
	m.CallRequest.RemoveEvent(idx)
	m.applyToMirrorResources(func(resource res.Resource) { resource.RemoveEvent(idx) })
}

func (m *callRequestMirror[PageReader]) ReaccessEvent() {
	m.CallRequest.ReaccessEvent()
	m.applyToMirrorResources(func(resource res.Resource) { resource.ReaccessEvent() })
}

func (m *callRequestMirror[PageReader]) ResetEvent() {
	m.CallRequest.ResetEvent()
	m.applyToMirrorResources(func(resource res.Resource) { resource.ResetEvent() })
}

func (m *callRequestMirror[PageReader]) QueryEvent(f func(request res.QueryRequest)) {
	m.CallRequest.QueryEvent(f)
	m.applyToMirrorResources(func(resource res.Resource) { resource.QueryEvent(f) })
}

func (m *callRequestMirror[PageReader]) CreateEvent(value any) {
	m.CallRequest.CreateEvent(value)
	m.applyToMirrorResources(func(resource res.Resource) { resource.CreateEvent(value) })
}

func (m *callRequestMirror[PageReader]) DeleteEvent() {
	m.CallRequest.DeleteEvent()
	m.applyToMirrorResources(func(resource res.Resource) { resource.DeleteEvent() })
}

func (m *callRequestMirror[PageReader]) applyToMirrorResources(f func(resource res.Resource)) {
	for m.mirrorModelRIDsPager.Next() {
		for _, rid := range m.mirrorModelRIDsPager.Page() {
			if err := m.Service().With(rid, func(resource res.Resource) { f(resource) }); err != nil {
				log.Default().Error("Could not process mirror resource",
					slog.Any("error", err),
					slog.String("rid", rid),
				)
			}
		}
	}

	if err := m.mirrorModelRIDsPager.Err(); err != nil {
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
