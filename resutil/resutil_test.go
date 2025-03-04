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

func TestMarshalJSONWithDataValues(t *testing.T) {
	type MetaValue struct {
		Value string `json:"value"`
	}

	tests := map[string]struct {
		in   any
		want string
	}{
		"simple": {
			in:   map[string]any{"id": "b0f71e4c-3dd5-4b5b-bbf6-869fcf05c1df"},
			want: `{"id":"b0f71e4c-3dd5-4b5b-bbf6-869fcf05c1df"}`,
		},
		"with data values": {
			in: map[string]any{
				"id":        "80b2e36e-bbaa-4616-9d0d-b3061333186e",
				"firstName": &MetaValue{Value: "John"},
				"languages": []*MetaValue{{Value: "en"}, {Value: "fr"}},
			},
			want: `{
				"id": "80b2e36e-bbaa-4616-9d0d-b3061333186e",
				"firstName": {
					"data": {"value": "John"}
				},
				"languages": {
					"data": [
						{"value": "en"},
						{"value": "fr"}
					]
				}
			}`,
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			got, err := MarshalJSONWithDataValues(tt.in)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestDeletable(t *testing.T) {
	tests := map[string]struct {
		in   []byte
		want *Deletable[string]
	}{
		"delete": {
			in:   res.DeleteAction.RawMessage,
			want: &Deletable[string]{Deleted: true},
		},
		"string": {
			in:   []byte(`"foo"`),
			want: &Deletable[string]{Value: "foo"},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			got := &Deletable[string]{}
			assert.NoError(t, json.Unmarshal(tt.in, got))
			assert.Equal(t, tt.want, got)
		})
	}
}

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
		ID        string                       `json:"id"`
		FirstName string                       `json:"firstName,omitempty"`
		LastName  string                       `json:"lastName,omitempty"`
		Emails    *res.DataValue[[]string]     `json:"emails,omitempty"`
		Phones    *res.DataValue[[]*UserPhone] `json:"phones,omitempty"`
	}

	tests := map[string]struct {
		previous, current *User
		want              string
	}{
		"simple": {
			previous: &User{
				ID:       "u-1",
				LastName: "Doe",
				Emails: &res.DataValue[[]string]{
					Data: []string{
						"john.doe@loungeup.com",
						"john@doe.com",
					},
				},
				Phones: &res.DataValue[[]*UserPhone]{
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
				Emails: &res.DataValue[[]string]{
					Data: []string{
						"john.doe@gmail.com",
						"john@doe.com",
					},
				},
				Phones: &res.DataValue[[]*UserPhone]{
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
