package client

import (
	"strings"
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/client/testdata"
	"github.com/loungeup/go-loungeup/resmodels"
	"github.com/loungeup/go-loungeup/transport"
	"github.com/loungeup/go-loungeup/transporttest"
	"github.com/stretchr/testify/assert"
)

func TestReadEntityIntegration(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.EntityIntegrationSelector.RID()):
					return transporttest.NewRESModelResponse(testdata.EntityIntegrationModel)
				case strings.HasSuffix(resourceID, testdata.IntegrationSelector.RID()):
					return transporttest.NewRESModelResponse(testdata.IntegrationModel)
				default:
					return transporttest.NewRESModelResponse(`{}`)
				}
			},
		},
	}, nil).Integrations.ReadEntityIntegration(testdata.EntityIntegrationSelector)
	assert.NoError(t, err)
	assert.Equal(t, testdata.EntityIntegration, got)
}

func TestReadEntityIntegrations(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.EntityIntegrationsSelector.RID()):
					return transporttest.NewRESCollectionResponse(testdata.EntityIntegrationCollection)
				case strings.HasSuffix(resourceID, testdata.EntityIntegrationSelector.RID()):
					return transporttest.NewRESModelResponse(testdata.EntityIntegrationModel)
				case strings.HasSuffix(resourceID, testdata.IntegrationSelector.RID()):
					return transporttest.NewRESModelResponse(testdata.IntegrationModel)
				default:
					return transporttest.NewRESModelResponse(`{}`)
				}
			},
		},
	}, nil).Integrations.ReadEntityIntegrations(testdata.EntityIntegrationsSelector)
	assert.NoError(t, err)
	assert.Equal(t, []*resmodels.EntityIntegration{testdata.EntityIntegration}, got)
}

func TestReadIntegration(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.IntegrationSelector.RID()):
					return transporttest.NewRESModelResponse(testdata.IntegrationModel)
				default:
					return transporttest.NewRESModelResponse(`{}`)
				}
			},
		},
	}, nil).Integrations.ReadIntegration(testdata.IntegrationSelector)
	assert.NoError(t, err)
	assert.Equal(t, testdata.Integration, got)
}

func TestReadIntegrations(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.IntegrationsSelector.RID()):
					return transporttest.NewRESCollectionResponse(testdata.IntegrationCollection)
				case strings.HasSuffix(resourceID, testdata.IntegrationSelector.RID()):
					return transporttest.NewRESModelResponse(testdata.IntegrationModel)
				default:
					return transporttest.NewRESModelResponse(`{}`)
				}
			},
		},
	}, nil).Integrations.ReadIntegrations(testdata.IntegrationsSelector)
	assert.NoError(t, err)
	assert.Equal(t, []*models.Integration{testdata.Integration}, got)
}

func TestFetchFromProvider(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, request resprot.Request) resprot.Response {
				return transporttest.NewRESResultResponse(testdata.ProviderResultModel)
			},
		},
	}, nil).Integrations.FetchFromProvider(testdata.EntityIntegrationSelector, nil)
	assert.NoError(t, err)
	assert.Equal(t, testdata.ProviderResult, got)
}

func TestSendToProvider(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, request resprot.Request) resprot.Response {
				return transporttest.NewRESResultResponse(testdata.ProviderResultModel)
			},
		},
	}, nil).Integrations.SendToProvider(testdata.EntityIntegrationSelector, nil)
	assert.NoError(t, err)
	assert.Equal(t, testdata.ProviderResult, got)
}
