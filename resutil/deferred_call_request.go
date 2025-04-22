package resutil

import "github.com/jirenius/go-res"

// DeferredCallRequest embeds a res.CallRequest and adds state to it. It stores the state of the request and send the
// reply when you explicitly call the Reply method.
type DeferredCallRequest struct {
	res.CallRequest

	err    error
	rid    string
	result any
}

// Error sets the error of the request.
func (request *DeferredCallRequest) Error(err error) { request.err = err }

// Resource sets the rid of the request.
func (request *DeferredCallRequest) Resource(rid string) { request.rid = rid }

// OK sets the result of the request.
func (request *DeferredCallRequest) OK(result any) { request.result = result }

func (request *DeferredCallRequest) GetError() error { return request.err }
func (request *DeferredCallRequest) GetRID() string  { return request.rid }
func (request *DeferredCallRequest) GetResult() any  { return request.result }

func (request *DeferredCallRequest) Reply() {
	if err := request.GetError(); err != nil {
		request.CallRequest.Error(request.err)
	} else if rid := request.GetRID(); rid != "" {
		request.CallRequest.Resource(request.rid)
	} else {
		request.CallRequest.OK(request.GetResult())
	}
}
