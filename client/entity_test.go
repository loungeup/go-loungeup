package client

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/cache"
	cacheMocks "github.com/loungeup/go-loungeup/cache/mocks"
	"github.com/loungeup/go-loungeup/resmodels"
	"github.com/loungeup/go-loungeup/transport"
	transportMocks "github.com/loungeup/go-loungeup/transport/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func initTest(t *testing.T) (*transportMocks.MockRESRequester, *cacheMocks.MockReadWriter) {
	ctrl := gomock.NewController(t)

	resClient := transportMocks.NewMockRESRequester(ctrl)
	cache := cacheMocks.NewMockReadWriter(ctrl)

	return resClient, cache
}

func newTransport(resClient transport.RESRequester, cache cache.ReadWriter) *Client {
	return NewWithTransport(&transport.Transport{
		RESClient: resClient,
	}, WithCache(cache))
}

type modelEntity struct {
	Model resmodels.Entity `json:"model"`
}

type modelEntityCustomFields struct {
	Model resmodels.EntityCustomFields `json:"model"`
}

type result struct {
	Collections []collection `json:"collection"`
}

type collection struct {
	Rid string `json:"rid"`
}

func TestReadEntity(t *testing.T) {
	uuid := uuid.New()
	entityAccount := createEntity(uuid.String())

	resClient, cache := initTest(t)

	t.Run("ReadEntity succeess", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		resourceID := "authority.entities." + uuid.String()
		expected := &entityAccount

		r := entityToRESresp(entityAccount)

		resClient.EXPECT().Request("get."+resourceID, resprot.Request{}).Return(r)

		resp, err := transportClient.Entities.ReadEntity(&resmodels.EntitySelector{EntityID: uuid})

		assert.NoError(t, err)
		assert.Equal(t, expected.ID, resp.ID)
	})

	t.Run("ReadEntity with cache", func(t *testing.T) {
		transportClient := newTransport(resClient, cache)

		resourceID := "authority.entities." + uuid.String()
		expected := &entityAccount

		cache.EXPECT().Read(resourceID).Return(expected)

		resp, err := transportClient.Entities.ReadEntity(&resmodels.EntitySelector{EntityID: uuid})

		assert.NoError(t, err)
		assert.Equal(t, expected.ID, resp.ID)
	})

	t.Run("ReadEntity with error: fail GetRESModel", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		resClient.EXPECT().Request("get.authority.entities."+uuid.String(), resprot.Request{}).Return(resprot.Response{})

		resp, err := transportClient.Entities.ReadEntity(&resmodels.EntitySelector{EntityID: uuid})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestReadEntityAccounts(t *testing.T) {
	parentAccountUUID := uuid.New()
	account1UUID := uuid.New()
	account2UUID := uuid.New()

	_ = createEntity(parentAccountUUID.String())

	account1 := createEntity(account1UUID.String())
	account1.Chain = res.SoftRef(parentAccountUUID.String())

	account2 := createEntity(account2UUID.String())
	account2.Chain = res.SoftRef(parentAccountUUID.String())

	resClient, _ := initTest(t)

	t.Run("ReadEntityAccounts success", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		resourceID := "authority.entities." + parentAccountUUID.String() + ".accounts"
		resourceIDEntity1 := "authority.entities." + account1UUID.String()
		resourceIDEntity2 := "authority.entities." + account2UUID.String()

		expected := []*resmodels.Entity{
			&account1,
			&account2,
		}

		resgateResp := result{
			Collections: []collection{
				{Rid: resourceIDEntity1},
				{Rid: resourceIDEntity2},
			},
		}
		respJSON, err := json.Marshal(resgateResp)
		assert.NoError(t, err)

		r := resprot.Response{
			Result: respJSON,
		}

		assert.NoError(t, err)

		// Entity 1
		resgateRespEntity1 := modelEntity{
			Model: account1,
		}
		respJSONEntity1, err := json.Marshal(resgateRespEntity1)
		assert.NoError(t, err)

		r1 := resprot.Response{
			Result: respJSONEntity1,
		}

		// Entity 2
		resgateRespEntity2 := modelEntity{
			Model: account2,
		}
		respJSONEntity2, err := json.Marshal(resgateRespEntity2)
		assert.NoError(t, err)

		r2 := resprot.Response{
			Result: respJSONEntity2,
		}

		resClient.EXPECT().Request("get."+resourceID, gomock.Any()).Return(r)
		resClient.EXPECT().Request("get."+resourceIDEntity1, resprot.Request{}).Return(r1)
		resClient.EXPECT().Request("get."+resourceIDEntity2, resprot.Request{}).Return(r2)

		selector := &resmodels.EntityAccountsSelector{
			EntityID: parentAccountUUID,
			Limit:    25,
			Offset:   0,
		}
		accountsResp, err := transportClient.Entities.ReadEntityAccounts(selector)
		assert.NoError(t, err)

		for i, account := range accountsResp {
			assert.Equal(t, expected[i].ID, account.ID)
		}
	})

	t.Run("ReadEntityAccounts with error: fail GetRESCollection", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		resourceID := "authority.entities." + parentAccountUUID.String() + ".accounts"

		resClient.EXPECT().Request("get."+resourceID, gomock.Any()).Return(resprot.Response{})

		selector := &resmodels.EntityAccountsSelector{
			EntityID: parentAccountUUID,
			Limit:    25,
			Offset:   0,
		}
		_, err := transportClient.Entities.ReadEntityAccounts(selector)
		assert.Error(t, err)
	})

	t.Run("ReadEntityAccounts with error: fail to get child entity", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		resourceID := "authority.entities." + parentAccountUUID.String() + ".accounts"
		resourceIDEntity1 := "authority.entities." + account1UUID.String()

		resgateResp := result{
			Collections: []collection{
				{Rid: resourceIDEntity1},
			},
		}
		respJSON, err := json.Marshal(resgateResp)
		assert.NoError(t, err)

		r := resprot.Response{
			Result: respJSON,
		}

		resClient.EXPECT().Request("get."+resourceID, gomock.Any()).Return(r)
		resClient.EXPECT().Request("get."+resourceIDEntity1, resprot.Request{}).Return(resprot.Response{})

		selector := &resmodels.EntityAccountsSelector{
			EntityID: parentAccountUUID,
			Limit:    25,
			Offset:   0,
		}
		_, err = transportClient.Entities.ReadEntityAccounts(selector)
		assert.Error(t, err)
	})
}

func TestReadAccountParents(t *testing.T) {
	resClient, _ := initTest(t)

	t.Run("ReadAccountParents success account has only a chain", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		uuidAccount := uuid.New()

		chainEntity := createEntity(uuidAccount.String())
		accountEntity := createEntity(uuidAccount.String())

		resourceID := "authority.entities." + accountEntity.ID
		chainResourceID := "authority.entities." + chainEntity.ID

		accountEntity.Chain = res.SoftRef(chainResourceID)

		expected := []*resmodels.Entity{
			&chainEntity,
		}

		resClient.EXPECT().Request("get."+resourceID, resprot.Request{}).Return(entityToRESresp(accountEntity))
		resClient.EXPECT().Request("get."+chainResourceID, resprot.Request{}).Return(entityToRESresp(chainEntity))

		resp, err := transportClient.Entities.ReadAccountParents(&resmodels.EntitySelector{EntityID: uuidAccount})
		assert.NoError(t, err)
		assert.Equal(t, expected[0].ID, resp[0].ID)
	})

	t.Run("ReadAccountParents success account has only a group", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		uuidAccount := uuid.New()

		groupEntity := createEntity(uuid.New().String())
		accountEntity := createEntity(uuidAccount.String())

		resourceID := "authority.entities." + accountEntity.ID
		groupResourceID := "authority.entities." + groupEntity.ID

		accountEntity.Group = res.SoftRef(groupResourceID)

		expected := []*resmodels.Entity{
			&groupEntity,
		}

		resClient.EXPECT().Request("get."+resourceID, resprot.Request{}).Return(entityToRESresp(accountEntity))
		resClient.EXPECT().Request("get."+groupResourceID, resprot.Request{}).Return(entityToRESresp(groupEntity))

		resp, err := transportClient.Entities.ReadAccountParents(&resmodels.EntitySelector{EntityID: uuidAccount})
		assert.NoError(t, err)
		assert.Equal(t, expected[0].ID, resp[0].ID)
	})

	t.Run("ReadAccountParents success account has both chain and group", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		uuidAccount := uuid.New()

		groupEntity := createEntity(uuid.New().String())
		chainEntity := createEntity(uuidAccount.String())
		accountEntity := createEntity(uuidAccount.String())

		resourceID := "authority.entities." + accountEntity.ID
		chainResourceID := "authority.entities." + chainEntity.ID
		groupResourceID := "authority.entities." + groupEntity.ID

		accountEntity.Chain = res.SoftRef(chainResourceID)
		accountEntity.Group = res.SoftRef(groupResourceID)

		expected := []*resmodels.Entity{
			&chainEntity,
			&groupEntity,
		}

		resClient.EXPECT().Request("get."+resourceID, resprot.Request{}).Return(entityToRESresp(accountEntity))
		resClient.EXPECT().Request("get."+chainResourceID, resprot.Request{}).Return(entityToRESresp(chainEntity))
		resClient.EXPECT().Request("get."+groupResourceID, resprot.Request{}).Return(entityToRESresp(groupEntity))

		resp, err := transportClient.Entities.ReadAccountParents(&resmodels.EntitySelector{EntityID: uuidAccount})
		assert.NoError(t, err)
		assert.Equal(t, expected[0].ID, resp[0].ID)
		assert.Equal(t, expected[1].ID, resp[1].ID)
	})

	t.Run("ReadAccountParents with error: fail to get chain entity", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		uuidAccount := uuid.New()
		accountEntity := createEntity(uuidAccount.String())
		resourceID := "authority.entities." + accountEntity.ID

		chainEntity := createEntity(uuid.New().String())
		accountEntity.Chain = res.SoftRef(chainEntity.ID)

		resClient.EXPECT().Request("get."+resourceID, resprot.Request{}).Return(entityToRESresp(accountEntity))
		resClient.EXPECT().Request("get."+chainEntity.ID, resprot.Request{}).Return(resprot.Response{
			Error: &res.Error{
				Code:    "internal",
				Message: "bruh this is a error",
			},
		})

		_, err := transportClient.Entities.ReadAccountParents(&resmodels.EntitySelector{EntityID: uuidAccount})
		assert.Error(t, err)
	})

	t.Run("ReadAccountParents with error: fail to get group entity", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		uuidAccount := uuid.New()
		accountEntity := createEntity(uuidAccount.String())
		resourceID := "authority.entities." + accountEntity.ID

		groupEntity := createEntity(uuid.New().String())
		accountEntity.Group = res.SoftRef(groupEntity.ID)

		resClient.EXPECT().Request("get."+resourceID, resprot.Request{}).Return(entityToRESresp(accountEntity))
		resClient.EXPECT().Request("get."+groupEntity.ID, resprot.Request{}).Return(resprot.Response{
			Error: &res.Error{
				Code:    "internal",
				Message: "bruh this is a error",
			},
		})

		_, err := transportClient.Entities.ReadAccountParents(&resmodels.EntitySelector{EntityID: uuidAccount})
		assert.Error(t, err)
	})

	t.Run("ReadAccountParents with error: error read entity account", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		uuidAccount := uuid.New()
		accountEntity := createEntity(uuidAccount.String())
		resourceID := "authority.entities." + accountEntity.ID

		resClient.EXPECT().Request("get."+resourceID, resprot.Request{}).Return(resprot.Response{
			Error: &res.Error{
				Code:    "internal",
				Message: "bruh this is a error",
			},
		})

		_, err := transportClient.Entities.ReadAccountParents(&resmodels.EntitySelector{EntityID: uuidAccount})
		assert.Error(t, err)
	})

	t.Run("ReadAccountParents with error: entity is not an account", func(t *testing.T) {
		transportClient := newTransport(resClient, nil)

		uuidAccount := uuid.New()
		entity := createEntity(uuidAccount.String())
		entity.Type = resmodels.EntityTypeGroup
		resourceID := "authority.entities." + entity.ID

		resClient.EXPECT().Request("get."+resourceID, resprot.Request{}).Return(entityToRESresp(entity))

		_, err := transportClient.Entities.ReadAccountParents(&resmodels.EntitySelector{EntityID: uuidAccount})
		assert.Error(t, err)
	})
}

func TestReadEntityCustomFields(t *testing.T) {
	transport, cache := initTest(t)

	t.Run("ReadEntityCustomFields success without cache", func(t *testing.T) {
		transportClient := newTransport(transport, nil)

		entityID := uuid.New()
		createEntity(entityID.String())
		resourceID := "proxy-db.entities." + entityID.String() + ".custom-fields"

		cFields := resmodels.EntityCustomFields{
			User: res.NewDataValue(map[string]resmodels.EntityCustomField{
				"field1": {
					Label: "stade-toulousain-overated",
					Type:  resmodels.EntityCustomFieldTypeText,
				},
			}),
		}
		expected := &cFields

		resResp := customFieldsToRESresp(cFields)
		transport.EXPECT().Request("get."+resourceID, resprot.Request{}).Return(resResp)

		req := &resmodels.EntityCustomFieldsSelector{
			EntityID: uuid.UUID(entityID),
		}
		resp, err := transportClient.Entities.ReadEntityCustomFields(req)
		assert.NoError(t, err)
		assert.Equal(t, expected.User.Data["field1"].Label, resp.User.Data["field1"].Label)
	})

	t.Run("ReadEntityCustomFields success with cache", func(t *testing.T) {
		transportClient := newTransport(transport, cache)

		entityID := uuid.New()
		createEntity(entityID.String())
		resourceID := "proxy-db.entities." + entityID.String() + ".custom-fields"

		cFields := resmodels.EntityCustomFields{
			User: res.NewDataValue(map[string]resmodels.EntityCustomField{
				"field1": {
					Label: "stade-toulousain-overated",
					Type:  resmodels.EntityCustomFieldTypeText,
				},
			}),
		}
		expected := &cFields

		cache.EXPECT().Read(resourceID).Return(expected)

		req := &resmodels.EntityCustomFieldsSelector{
			EntityID: uuid.UUID(entityID),
		}
		resp, err := transportClient.Entities.ReadEntityCustomFields(req)
		assert.NoError(t, err)
		assert.Equal(t, expected.User.Data["field1"].Label, resp.User.Data["field1"].Label)
	})

	t.Run("ReadEntityCustomFields success: empty cache and write in cache", func(t *testing.T) {
		transportClient := newTransport(transport, cache)

		entityID := uuid.New()
		createEntity(entityID.String())
		resourceID := "proxy-db.entities." + entityID.String() + ".custom-fields"

		cFields := resmodels.EntityCustomFields{
			User: res.NewDataValue(map[string]resmodels.EntityCustomField{
				"field1": {
					Label: "stade-toulousain-overated",
					Type:  resmodels.EntityCustomFieldTypeText,
				},
			}),
		}
		expected := &cFields

		resResp := customFieldsToRESresp(cFields)

		cache.EXPECT().Read(resourceID).Return(nil)
		transport.EXPECT().Request("get."+resourceID, resprot.Request{}).Return(resResp)
		cache.EXPECT().Write(resourceID, expected)

		req := &resmodels.EntityCustomFieldsSelector{
			EntityID: uuid.UUID(entityID),
		}
		resp, err := transportClient.Entities.ReadEntityCustomFields(req)
		assert.NoError(t, err)
		assert.Equal(t, expected.User.Data["field1"].Label, resp.User.Data["field1"].Label)
	})

	t.Run("ReadEntityCustomFields with error: fail GetRESModel", func(t *testing.T) {
		transportClient := newTransport(transport, nil)

		entityID := uuid.New()
		createEntity(entityID.String())
		resourceID := "proxy-db.entities." + entityID.String() + ".custom-fields"

		transport.EXPECT().Request("get."+resourceID, resprot.Request{}).Return(resprot.Response{})

		req := &resmodels.EntityCustomFieldsSelector{
			EntityID: uuid.UUID(entityID),
		}

		_, err := transportClient.Entities.ReadEntityCustomFields(req)
		assert.Error(t, err)
	})
}

func TestPatchEntity(t *testing.T) {
	transport, _ := initTest(t)

	t.Run("PatchEntity success", func(t *testing.T) {
		transportClient := newTransport(transport, nil)

		entityID := uuid.New()
		selector := &resmodels.EntitySelector{
			EntityID: uuid.UUID(entityID),
		}
		// rid := resmodels.EntityID(entityID)

		updates := &resmodels.EntityUpdates{
			ConvertAmountsTaskRID:    "convertAmountsTaskRID",
			IndexGuestProfile:        true,
			IndexGuestProfileTaskRID: "indexGuestProfileTaskRID",
		}

		encodedUpdates, err := json.Marshal(updates)
		assert.NoError(t, err)

		resourceID := "call.authority.entities." + entityID.String() + ".patch"

		transport.EXPECT().Request(resourceID, resprot.Request{Params: json.RawMessage(encodedUpdates)}).Return(resprot.Response{})

		err = transportClient.Entities.PatchEntity(selector, updates)
	})

	t.Run("PatchEntity with error: response error", func(t *testing.T) {
		transportClient := newTransport(transport, nil)

		entityID := uuid.New()
		selector := &resmodels.EntitySelector{
			EntityID: uuid.UUID(entityID),
		}

		updates := &resmodels.EntityUpdates{
			ConvertAmountsTaskRID:    "convertAmountsTaskRID",
			IndexGuestProfile:        true,
			IndexGuestProfileTaskRID: "indexGuestProfileTaskRID",
		}

		encodedUpdates, err := json.Marshal(updates)
		assert.NoError(t, err)

		resourceID := "call.authority.entities." + entityID.String() + ".patch"

		transport.EXPECT().Request(resourceID, resprot.Request{Params: json.RawMessage(encodedUpdates)}).Return(resprot.Response{
			Error: &res.Error{
				Code:    "internal",
				Message: "bruh this is a error",
			},
		})

		err = transportClient.Entities.PatchEntity(selector, updates)
		assert.Error(t, err)
	})
}

func TestBuildESQueryEntity(t *testing.T) {
	t.Run("BuildESQueryEntity success", func(t *testing.T) {
		transport, _ := initTest(t)
		transportClient := newTransport(transport, nil)

		entityID := uuid.New()
		params := &resmodels.BuildEntityESQueryParams{}

		transport.EXPECT().Request("call.guestprofile.entities."+entityID.String()+".build-elasticsearch-query", resprot.Request{
			Params: params,
			Token:  json.RawMessage(`{"agentRoles": ["service"]}`),
		}).Return(resprot.Response{
			Result: json.RawMessage(`{
				"result": {
					"bool": {
						"filter": {
							"bool": {
								"must": [
									{
										"term": {
											"guest.account.entityId": "` + entityID.String() + `"
										}
									}
								]
							}
						}
					}
				}
			}`),
		})

		resp, err := transportClient.Entities.BuildESQueryEntity(&resmodels.EntitySelector{
			EntityID: entityID,
		}, params)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
	})
}

func createEntity(uuid string) resmodels.Entity {
	lang := res.NewDataValue([]string{"en"})

	return resmodels.Entity{
		ID:             uuid,
		LegacyID:       1,
		Type:           resmodels.EntityTypeAccount,
		Name:           "Test Account",
		Slug:           "testaccount",
		Image:          "https://example.com/image.jpg",
		Languages:      &lang,
		Timezone:       "Europe/Paris",
		Country:        "FR",
		PostalCode:     "31520",
		City:           "Ramonville-Saint-Agne",
		Address:        "12 avenue de l'Europe",
		Rooms:          100,
		Currency:       res.SoftRef("authority.currencies.eur"),
		ConvertAmounts: true,
		Chain:          res.SoftRef(""),
		Group:          res.SoftRef(""),
		CreatedAt:      "2020-01-01T00:00:00Z",
		UpdatedAt:      "2020-01-01T00:00:00Z",
	}
}

func entityToRESresp(e resmodels.Entity) resprot.Response {
	model := modelEntity{
		Model: e,
	}

	respJSON, err := json.Marshal(model)
	if err != nil {
		panic(err)
	}

	return resprot.Response{
		Result: respJSON,
	}
}

func customFieldsToRESresp(e resmodels.EntityCustomFields) resprot.Response {
	model := modelEntityCustomFields{
		Model: e,
	}

	respJSON, err := json.Marshal(model)
	if err != nil {
		panic(err)
	}

	return resprot.Response{
		Result: respJSON,
	}
}
