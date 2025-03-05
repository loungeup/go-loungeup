package client

import (
	"strings"
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/client/testdata"
	"github.com/loungeup/go-loungeup/transport"
	"github.com/loungeup/go-loungeup/transporttest"
	"github.com/stretchr/testify/assert"
)

func TestReadEntity(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.EntitySelector.RID()):
					return transporttest.NewRESModelResponse(testdata.EntityModel)
				default:
					return transporttest.NewRESModelResponse(`{}`)
				}
			},
		},
	}).Internal.Entities.ReadEntity(testdata.EntitySelector)
	assert.NoError(t, err)
	assert.Equal(t, testdata.Entity, got)
}

func TestReadEntityAccounts(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.EntityAccountsSelector.RID()):
					return transporttest.NewRESCollectionResponse(testdata.EntityCollection)
				case strings.HasSuffix(resourceID, testdata.EntitySelector.RID()):
					return transporttest.NewRESModelResponse(testdata.EntityModel)
				default:
					return transporttest.NewRESModelResponse(`{}`)
				}
			},
		},
	}).Internal.Entities.ReadEntityAccounts(testdata.EntityAccountsSelector)
	assert.NoError(t, err)
	assert.Equal(t, []*models.Entity{testdata.Entity}, got)
}

func TestReadEntityCustomFields(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.EntityCustomFieldsSelector.RID()):
					return transporttest.NewRESModelResponse(testdata.EntityCustomFieldsModel)
				default:
					return transporttest.NewRESModelResponse(`{}`)
				}
			},
		},
	}).Internal.Entities.ReadEntityCustomFields(testdata.EntityCustomFieldsSelector)
	assert.NoError(t, err)
	assert.Equal(t, testdata.EntityCustomFields, got)
}

func TestReadAccountParents(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.AccountChainSelector.RID()):
					return transporttest.NewRESModelResponse(testdata.ChainModel)
				case strings.HasSuffix(resourceID, testdata.AccountGroupSelector.RID()):
					return transporttest.NewRESModelResponse(testdata.GroupModel)
				case strings.HasSuffix(resourceID, testdata.EntitySelector.RID()):
					return transporttest.NewRESModelResponse(testdata.EntityModel)
				default:
					return transporttest.NewRESModelResponse(`{}`)
				}
			},
		},
	}).Internal.Entities.ReadAccountParents(testdata.EntitySelector)
	assert.NoError(t, err)
	assert.Equal(t, []*models.Entity{testdata.EntityChain, testdata.EntityGroup}, got)
}
