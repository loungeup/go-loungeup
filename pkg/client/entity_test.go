package client

import (
	"strings"
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/client/testdata"
	"github.com/loungeup/go-loungeup/pkg/transport"
	"github.com/loungeup/go-loungeup/pkg/transporttest"
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
