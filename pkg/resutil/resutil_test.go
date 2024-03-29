package resutil

import (
	"encoding/json"
	"testing"

	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/restest"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

var testAccessHandler = res.Access(res.AccessGranted)

func TestAddNATSMessageHandler(t *testing.T) {
	service := res.NewService("test")
	service.Handle("test.foo", testAccessHandler)
	service.SetOnServe(func(service *res.Service) {
		assert.NoError(t, AddNATSMessageHandler(service, "test.foo", func(message *nats.Msg) error { return nil }))
	})

	session := restest.NewSession(t, service)
	defer session.Close()

	session.SendMessage("test.foo", "", nil)
}

func TestHandleCollectionQueryRequest(t *testing.T) {
	const testQuery = "foo=bar"

	service := res.NewService("users-manager")
	service.Handle("users",
		testAccessHandler,
		res.GetCollection(func(request res.CollectionRequest) {
			request.QueryCollection([]string{"john.doe", "jane.doe"}, testQuery)
		}),
	)
	service.Handle("users.latest",
		testAccessHandler,
		res.GetModel(func(request res.ModelRequest) {
			request.QueryModel("jane.doe", testQuery)
		}),
		res.Call("delete", func(request res.CallRequest) {
			request.OK(nil)

			assert.NoError(t, HandleCollectionQueryRequest(
				service,
				"users-manager.users",
				func(request res.QueryRequest) ([]string, error) {
					return []string{"john.doe"}, nil
				},
			))

			assert.NoError(t, HandleModelQueryRequest(
				service,
				"users-manager.users.latest",
				func(request res.QueryRequest) (string, error) {
					return "john.doe", nil
				},
			))
		}),
	)

	session := restest.NewSession(t, service)
	defer session.Close()

	session.Call("users-manager.users.latest?"+testQuery, "delete", nil).Response()

	var latestUserQueryEventSubject string
	session.GetMsg().AssertQueryEvent("users-manager.users.latest", &latestUserQueryEventSubject)

	var usersQueryEventSubject string
	session.GetMsg().AssertQueryEvent("users-manager.users", &usersQueryEventSubject)

	session.QueryRequest(latestUserQueryEventSubject, testQuery).Response().AssertResult(json.RawMessage(`{
		"model": "john.doe"
	}`))
	session.QueryRequest(usersQueryEventSubject, testQuery).Response().AssertResult(json.RawMessage(`{
		"collection": [
			"john.doe"
		]
	}`))
}
