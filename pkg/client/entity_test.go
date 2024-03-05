package client

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/transport"
	"github.com/stretchr/testify/assert"
)

func TestReadEntityAccounts(t *testing.T) {
	want := []Entity{
		{
			ID:   uuid.MustParse("acec14d0-1897-478b-ac80-009ad0b9508a"),
			Name: "Foo",
		},
	}

	got, err := NewWithTransport(
		&transport.Transport{
			RESClient: &resClientMock{
				requestFunc: func(resourceID string, request resprot.Request) resprot.Response {
					if strings.Contains(resourceID, "accounts") {
						return resprot.Response{
							Result: json.RawMessage(`{
								"collection": [
									{"rid": "authority.entities.acec14d0-1897-478b-ac80-009ad0b9508a"}
								]
							}`),
						}
					}

					return resprot.Response{
						Result: json.RawMessage(`{
							"model": {
								"id": "acec14d0-1897-478b-ac80-009ad0b9508a",
								"name": "Foo"
							}
						}`),
					}
				},
			},
		},
	).Internal.Entities.ReadEntityAccounts(EntityAccountsSelector{
		EntitySelector: EntitySelector{
			ID: uuid.MustParse("acec14d0-1897-478b-ac80-009ad0b9508a"),
		},
		Limit:  10,
		Offset: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

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

func TestEntityAccountsSelector(t *testing.T) {
	tests := map[string]struct {
		in               EntityAccountsSelector
		wantResourceID   string
		wantEncodedQuery string
	}{
		"simple": {
			in: EntityAccountsSelector{
				EntitySelector: EntitySelector{
					ID: uuid.MustParse("acec14d0-1897-478b-ac80-009ad0b9508a"),
				},
			},
			wantResourceID:   "authority.entities.acec14d0-1897-478b-ac80-009ad0b9508a.accounts",
			wantEncodedQuery: "",
		},
		"with limit": {
			in: EntityAccountsSelector{
				EntitySelector: EntitySelector{
					ID: uuid.MustParse("acec14d0-1897-478b-ac80-009ad0b9508a"),
				},
				Limit: 10,
			},
			wantResourceID:   "authority.entities.acec14d0-1897-478b-ac80-009ad0b9508a.accounts",
			wantEncodedQuery: "limit=10",
		},
		"with offset": {
			in: EntityAccountsSelector{
				EntitySelector: EntitySelector{
					ID: uuid.MustParse("acec14d0-1897-478b-ac80-009ad0b9508a"),
				},
				Limit:  10,
				Offset: 10,
			},
			wantResourceID:   "authority.entities.acec14d0-1897-478b-ac80-009ad0b9508a.accounts",
			wantEncodedQuery: "limit=10&offset=10",
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.wantResourceID, tt.in.resourceID())
			assert.Equal(t, tt.wantEncodedQuery, tt.in.encodeQuery())
		})
	}
}

func TestReadRoomTypes(t *testing.T) {
	roomTypeID := "5b6ad30a-f765-4247-9925-a13380ba284c"
	entityID := "acec14d0-1897-478b-ac80-009ad0b9508a"
	want := []RoomType{
		{
			ID:       uuid.MustParse(roomTypeID),
			EntityID: uuid.MustParse(entityID),
			Name:     "Foo",
		},
	}

	got, err := NewWithTransport(
		&transport.Transport{
			RESClient: &resClientMock{
				requestFunc: func(resourceID string, request resprot.Request) resprot.Response {
					if strings.Contains(resourceID, roomTypeID) {
						return resprot.Response{
							Result: json.RawMessage(`{
								"model": {
									"id": "` + roomTypeID + `",
									"entityId": "` + entityID + `",
									"name": "Foo"
								}
							}`),
						}
					}

					return resprot.Response{
						Result: json.RawMessage(`{
							"collection": [
								{"rid": "authority.entities.` + entityID + `.room-types.` + roomTypeID + `"}
							]
						}`),
					}
				},
			},
		},
	).Internal.Entities.ReadRoomTypes(RoomTypesSelector{
		EntitySelector: EntitySelector{
			ID: uuid.MustParse(entityID),
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
