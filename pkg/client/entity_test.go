package client

import (
	"encoding/json"
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/transport"
	"github.com/stretchr/testify/assert"
)

func TestReadEntityCustomFields(t *testing.T) {
	want := EntityCustomFields{
		User: RESDataValue[map[string]EntityCustomField]{
			Data: map[string]EntityCustomField{
				"foo": {
					Label: "Foo",
					Type:  EntityCustomFieldTypeText,
				},
			},
		},
		Visit: RESDataValue[map[string]EntityCustomField]{
			Data: map[string]EntityCustomField{
				"bar": {
					Label: "Bar",
					Type:  EntityCustomFieldTypeBoolean,
				},
			},
		},
	}

	got, err := NewWithTransport(
		&transport.Transport{
			RESClient: &resClientMock{
				requestFunc: func(resourceID string, request resprot.Request) resprot.Response {
					return resprot.Response{
						Result: json.RawMessage(`{
							"model": {
								"user": {
									"data": {
										"foo": {
											"label": "Foo",
											"type": "text"
										}
									}
								},
								"visit": {
									"data": {
										"bar": {
											"label": "Bar",
											"type": "boolean"
										}
									}
								}
							}
						}`),
					}
				},
			},
		},
	).Internal.Entities.ReadEntityCustomFields(EntitySelector{})
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
