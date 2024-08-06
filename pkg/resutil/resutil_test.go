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

func TestCompareModels(t *testing.T) {
	type UserPhone struct {
		Number string `json:"number,omitempty"`
	}

	type User struct {
		ID        string         `json:"id"`
		FirstName string         `json:"firstName,omitempty"`
		LastName  string         `json:"lastName,omitempty"`
		Emails    *res.DataValue `json:"emails,omitempty"`
		Phones    *res.DataValue `json:"phones,omitempty"`
	}

	tests := map[string]struct {
		previous, current *User
		want              string
	}{
		"simple": {
			previous: &User{
				ID:       "u-1",
				LastName: "Doe",
				Emails: &res.DataValue{
					Data: []string{
						"john.doe@loungeup.com",
						"john@doe.com",
					},
				},
				Phones: &res.DataValue{
					Data: []*UserPhone{
						{
							Number: "+33 612345678",
						},
					},
				},
			},
			current: &User{
				ID:        "u-1",
				FirstName: "Jane",
				LastName:  "",
				Emails: &res.DataValue{
					Data: []string{
						"john.doe@gmail.com",
						"john@doe.com",
					},
				},
				Phones: &res.DataValue{
					Data: []*UserPhone{
						{
							Number: "+33 512345678",
						},
					},
				},
			},
			want: `{
				"firstName": "Jane",
				"lastName": {"action": "delete"},
				"emails": {
					"data": [
						"john.doe@gmail.com",
						"john@doe.com"
					]
				},
				"phones": {
					"data": [
						{
							"number": "+33 512345678"
						}
					]
				}
			}`,
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			got, err := json.Marshal(CompareModels(tt.previous, tt.current))
			assert.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
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

func TestMapRefs(t *testing.T) {
	tests := map[string]struct {
		in   []string
		want RefSlice
	}{
		"empty": {
			in:   []string{},
			want: RefSlice{},
		},
		"simple reference": {
			in: []string{"foo"},
			want: RefSlice{
				res.Ref("foo"),
			},
		},
		"ignore invalid reference": {
			in: []string{"foo", "."},
			want: RefSlice{
				res.Ref("foo"),
			},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, MapRefs(tt.in, func(e string) res.Ref {
				return res.Ref(e)
			}))
		})
	}
}
