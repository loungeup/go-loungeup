package client

import (
	"encoding/json"

	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

// integrationsClient used to interact with the integrations service using the RES protocol.
type integrationsClient struct{ baseClient *Client }

func (c *integrationsClient) ReadEntityIntegration(
	selector *models.EntityIntegrationSelector,
) (*models.EntityIntegration, error) {
	return c.readEntityIntegrationByRID(selector.RID())
}

func (c *integrationsClient) ReadEntityIntegrations(
	selector *models.EntityIntegrationsSelector,
) ([]*models.EntityIntegration, error) {
	cacheKey := selector.RID() + "?" + selector.EncodedQuery()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(cacheKey).([]*models.EntityIntegration); ok {
		return cachedResult, nil
	}

	references, err := transport.GetRESCollection[res.Ref](
		c.baseClient.resClient,
		selector.RID(),
		resprot.Request{Query: selector.EncodedQuery()},
	)
	if err != nil {
		return nil, err
	}

	result := []*models.EntityIntegration{}

	for _, reference := range references {
		model, err := c.readEntityIntegrationByRID(string(reference))
		if err != nil {
			return nil, err
		}

		result = append(result, model)
	}

	defer c.baseClient.eventuallyWriteCache(cacheKey, result)

	return result, nil
}

func (c *integrationsClient) ReadIntegration(selector *models.IntegrationSelector) (*models.Integration, error) {
	return c.readIntegrationByRID(selector.RID())
}

func (c *integrationsClient) ReadIntegrations(selector *models.IntegrationsSelector) ([]*models.Integration, error) {
	cacheKey := selector.RID() + "?" + selector.EncodedQuery()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(cacheKey).([]*models.Integration); ok {
		return cachedResult, nil
	}

	references, err := transport.GetRESCollection[res.Ref](
		c.baseClient.resClient,
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

	defer c.baseClient.eventuallyWriteCache(cacheKey, result)

	return result, nil
}

func (c *integrationsClient) FetchFromProvider(
	selector *models.EntityIntegrationSelector,
	params any,
) (json.RawMessage, error) {
	return transport.CallRESResult[json.RawMessage](
		c.baseClient.resClient,
		selector.RID()+".fetch-from-provider",
		resprot.Request{Params: params},
	)
}

func (c *integrationsClient) FetchLatestEntityIntegrationRoomTypes(
	selector *models.LatestIntegrationSelector,
) ([]*models.RoomType, error) {
	return transport.CallRESResult[[]*models.RoomType](
		c.baseClient.resClient,
		selector.RID()+".fetch-room-types",
		resprot.Request{
			Query: selector.EncodedQuery(),
		},
	)
}

func (c *integrationsClient) CreateTicket(
	selector *models.EntityIntegrationSelector,
	params any,
) (json.RawMessage, error) {
	return transport.CallRESResult[json.RawMessage](
		c.baseClient.resClient,
		selector.RID()+".tickets.create",
		resprot.Request{Params: params},
	)
}

func (c *integrationsClient) SendToProvider(
	selector *models.EntityIntegrationSelector,
	params any,
) (json.RawMessage, error) {
	return transport.CallRESResult[json.RawMessage](
		c.baseClient.resClient,
		selector.RID()+".send-to-provider",
		resprot.Request{Params: params},
	)
}

func (c *integrationsClient) readEntityIntegrationByRID(resourceID string) (*models.EntityIntegration, error) {
	if cachedResult, ok := c.baseClient.eventuallyReadCache(resourceID).(*models.EntityIntegration); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*models.EntityIntegration](c.baseClient.resClient, resourceID)
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

	defer c.baseClient.eventuallyWriteCache(resourceID, result)

	return result, nil
}

func (c *integrationsClient) readIntegrationByRID(resourceID string) (*models.Integration, error) {
	if cachedResult, ok := c.baseClient.eventuallyReadCache(resourceID).(*models.Integration); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*models.Integration](c.baseClient.resClient, resourceID)
	if err != nil {
		return nil, err
	}

	defer c.baseClient.eventuallyWriteCache(resourceID, result)

	return result, nil
}
