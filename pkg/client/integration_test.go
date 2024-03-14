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

func TestReadEntityIntegration(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &resClientMock{
			requestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.EntityIntegrationSelector.RID()):
					return newModelResponse(testdata.EntityIntegrationModel)
				case strings.HasSuffix(resourceID, testdata.IntegrationSelector.RID()):
					return newModelResponse(testdata.IntegrationModel)
				default:
					return newModelResponse(`{}`)
				}
			},
		},
	}).Internal.Integrations.ReadEntityIntegration(testdata.EntityIntegrationSelector)
	assert.NoError(t, err)
	assert.Equal(t, testdata.EntityIntegration, got)
}

func TestReadEntityIntegrations(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &resClientMock{
			requestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.EntityIntegrationsSelector.RID()):
					return newCollectionResponse(testdata.EntityIntegrationCollection)
				case strings.HasSuffix(resourceID, testdata.EntityIntegrationSelector.RID()):
					return newModelResponse(testdata.EntityIntegrationModel)
				case strings.HasSuffix(resourceID, testdata.IntegrationSelector.RID()):
					return newModelResponse(testdata.IntegrationModel)
				default:
					return newModelResponse(`{}`)
				}
			},
		},
	}).Internal.Integrations.ReadEntityIntegrations(testdata.EntityIntegrationsSelector)
	assert.NoError(t, err)
	assert.Equal(t, []*models.EntityIntegration{testdata.EntityIntegration}, got)
}

func TestReadIntegration(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &resClientMock{
			requestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.IntegrationSelector.RID()):
					return newModelResponse(testdata.IntegrationModel)
				default:
					return newModelResponse(`{}`)
				}
			},
		},
	}).Internal.Integrations.ReadIntegration(testdata.IntegrationSelector)
	assert.NoError(t, err)
	assert.Equal(t, testdata.Integration, got)
}

func TestReadIntegrations(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &resClientMock{
			requestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.IntegrationsSelector.RID()):
					return newCollectionResponse(testdata.IntegrationCollection)
				case strings.HasSuffix(resourceID, testdata.IntegrationSelector.RID()):
					return newModelResponse(testdata.IntegrationModel)
				default:
					return newModelResponse(`{}`)
				}
			},
		},
	}).Internal.Integrations.ReadIntegrations(testdata.IntegrationsSelector)
	assert.NoError(t, err)
	assert.Equal(t, []*models.Integration{testdata.Integration}, got)
}

func TestFetchFromProvider(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &resClientMock{
			requestFunc: func(resourceID string, request resprot.Request) resprot.Response {
				return newResultResponse(testdata.ProviderResultModel)
			},
		},
	}).Internal.Integrations.FetchFromProvider(testdata.EntityIntegrationSelector, nil)
	assert.NoError(t, err)
	assert.Equal(t, testdata.ProviderResult, got)
}

func TestSendToProvider(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &resClientMock{
			requestFunc: func(resourceID string, request resprot.Request) resprot.Response {
				return newResultResponse(testdata.ProviderResultModel)
			},
		},
	}).Internal.Integrations.SendToProvider(testdata.EntityIntegrationSelector, nil)
	assert.NoError(t, err)
	assert.Equal(t, testdata.ProviderResult, got)
}
