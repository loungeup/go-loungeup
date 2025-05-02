package client

import (
	"encoding/json"

	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/resmodels"
	"github.com/loungeup/go-loungeup/transport"
)

//go:generate mockgen -source integration.go -destination=./mocks/mock_integration.go -package=mocks

type IntegrationsManager interface {
	ReadEntityIntegration(selector *resmodels.EntityIntegrationSelector) (*resmodels.EntityIntegration, error)
	UpdateEntityIntegration(selector *resmodels.EntityIntegrationSelector, params any) (resprot.Response, error)
	ReadEntityIntegrations(selector *resmodels.EntityIntegrationsSelector) ([]*resmodels.EntityIntegration, error)
	ReadIntegration(selector *models.IntegrationSelector) (*models.Integration, error)
	ReadIntegrations(selector *models.IntegrationsSelector) ([]*models.Integration, error)
	FetchFromProvider(selector *resmodels.EntityIntegrationSelector, params any) (json.RawMessage, error)
	FetchLatestEntityIntegrationRoomTypes(selector *resmodels.LatestIntegrationSelector) ([]*models.RoomType, error)
	CreateTicket(selector *resmodels.EntityIntegrationSelector, params any) (json.RawMessage, error)
	SendToProvider(selector *resmodels.EntityIntegrationSelector, params any) (json.RawMessage, error)
}

type IntegrationsClient struct {
	base *BaseClient
}

func NewIntegrationsClient(base *BaseClient) *IntegrationsClient {
	return &IntegrationsClient{
		base: base,
	}
}

func (c *IntegrationsClient) ReadEntityIntegration(
	selector *resmodels.EntityIntegrationSelector,
) (*resmodels.EntityIntegration, error) {
	return c.readEntityIntegrationByRID(selector.RID())
}

func (c *IntegrationsClient) UpdateEntityIntegration(
	selector *resmodels.EntityIntegrationSelector,
	params any,
) (resprot.Response, error) {
	response := c.base.resClient.Request("call."+selector.RID()+".patch",
		resprot.Request{Params: params})

	if response.HasError() {
		return resprot.Response{}, response.Error
	}

	return response, nil
}

func (c *IntegrationsClient) ReadEntityIntegrations(
	selector *resmodels.EntityIntegrationsSelector,
) ([]*resmodels.EntityIntegration, error) {
	cacheKey := selector.RID() + "?" + selector.EncodedQuery()

	if cachedResult, ok := c.base.ReadCache(cacheKey).([]*resmodels.EntityIntegration); ok {
		return cachedResult, nil
	}

	references, err := transport.GetRESCollection[res.Ref](
		c.base.resClient,
		selector.RID(),
		resprot.Request{Query: selector.EncodedQuery()},
	)
	if err != nil {
		return nil, err
	}

	result := []*resmodels.EntityIntegration{}

	for _, reference := range references {
		model, err := c.readEntityIntegrationByRID(string(reference))
		if err != nil {
			return nil, err
		}

		result = append(result, model)
	}

	defer c.base.WriteCache(cacheKey, result)

	return result, nil
}

func (c *IntegrationsClient) ReadIntegration(selector *models.IntegrationSelector) (*models.Integration, error) {
	return c.readIntegrationByRID(selector.RID())
}

func (c *IntegrationsClient) ReadIntegrations(selector *models.IntegrationsSelector) ([]*models.Integration, error) {
	cacheKey := selector.RID() + "?" + selector.EncodedQuery()

	if cachedResult, ok := c.base.ReadCache(cacheKey).([]*models.Integration); ok {
		return cachedResult, nil
	}

	references, err := transport.GetRESCollection[res.Ref](
		c.base.resClient,
		selector.RID(),
		resprot.Request{Query: selector.EncodedQuery()},
	)
	if err != nil {
		return nil, err
	}

	result := []*models.Integration{}

	for _, reference := range references {
		model, err := c.readIntegrationByRID(string(reference))
		if err != nil {
			return nil, err
		}

		result = append(result, model)
	}

	defer c.base.WriteCache(cacheKey, result)

	return result, nil
}

func (c *IntegrationsClient) FetchFromProvider(
	selector *resmodels.EntityIntegrationSelector,
	params any,
) (json.RawMessage, error) {
	return transport.CallRESResult[json.RawMessage](
		c.base.resClient,
		selector.RID()+".fetch-from-provider",
		resprot.Request{Params: params},
	)
}

func (c *IntegrationsClient) FetchLatestEntityIntegrationRoomTypes(
	selector *resmodels.LatestIntegrationSelector,
) ([]*models.RoomType, error) {
	return transport.CallRESResult[[]*models.RoomType](
		c.base.resClient,
		selector.RID()+".fetch-room-types",
		resprot.Request{
			Query: selector.EncodedQuery(),
		},
	)
}

func (c *IntegrationsClient) CreateTicket(
	selector *resmodels.EntityIntegrationSelector,
	params any,
) (json.RawMessage, error) {
	return transport.CallRESResult[json.RawMessage](
		c.base.resClient,
		selector.RID()+".create-ticket",
		resprot.Request{Params: params},
	)
}

func (c *IntegrationsClient) SendToProvider(
	selector *resmodels.EntityIntegrationSelector,
	params any,
) (json.RawMessage, error) {
	return transport.CallRESResult[json.RawMessage](
		c.base.resClient,
		selector.RID()+".send-to-provider",
		resprot.Request{Params: params},
	)
}

func (c *IntegrationsClient) readEntityIntegrationByRID(resourceID string) (*resmodels.EntityIntegration, error) {
	if cachedResult, ok := c.base.ReadCache(resourceID).(*resmodels.EntityIntegration); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*resmodels.EntityIntegration](c.base.resClient, resourceID, resprot.Request{})
	if err != nil {
		return nil, err
	}

	if !result.IntegrationReference.IsValid() {
		return result, nil
	}

	relatedIntegration, err := c.readIntegrationByRID(string(result.IntegrationReference))
	if err != nil {
		return nil, err
	}

	result.Integration = relatedIntegration

	defer c.base.WriteCache(resourceID, result)

	return result, nil
}

func (c *IntegrationsClient) readIntegrationByRID(resourceID string) (*models.Integration, error) {
	if cachedResult, ok := c.base.ReadCache(resourceID).(*models.Integration); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*models.Integration](c.base.resClient, resourceID, resprot.Request{})
	if err != nil {
		return nil, err
	}

	defer c.base.WriteCache(resourceID, result)

	return result, nil
}
