package client

import (
	"strings"
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/client/testdata"
	"github.com/loungeup/go-loungeup/pkg/transport"
	"github.com/stretchr/testify/assert"
)

func TestReadEntity(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &resClientMock{
			requestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.EntitySelector.RID()):
					return newModelResponse(testdata.EntityModel)
				default:
					return newModelResponse(`{}`)
				}
			},
		},
	}).Internal.Entities.ReadEntity(testdata.EntitySelector)
	assert.NoError(t, err)
	assert.Equal(t, testdata.Entity, got)
}

func TestReadEntityAccounts(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &resClientMock{
			requestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.EntityAccountsSelector.RID()):
					return newCollectionResponse(testdata.EntityCollection)
				case strings.HasSuffix(resourceID, testdata.EntitySelector.RID()):
					return newModelResponse(testdata.EntityModel)
				default:
					return newModelResponse(`{}`)
				}
			},
		},
	}).Internal.Entities.ReadEntityAccounts(testdata.EntityAccountsSelector)
	assert.NoError(t, err)
	assert.Equal(t, []models.Entity{testdata.Entity}, got)
}

func TestReadEntityCustomFields(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &resClientMock{
			requestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.EntityCustomFieldsSelector.RID()):
					return newModelResponse(testdata.EntityCustomFieldsModel)
				default:
					return newModelResponse(`{}`)
				}
			},
		},
	}).Internal.Entities.ReadEntityCustomFields(testdata.EntityCustomFieldsSelector)
	assert.NoError(t, err)
	assert.Equal(t, testdata.EntityCustomFields, got)
}
