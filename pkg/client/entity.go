package client

import (
	"encoding/json"
	"fmt"

	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

// entitiesClient provides methods to interact with entities.
type entitiesClient struct{ baseClient *Client }

func (c *entitiesClient) BuildESQuery(
	selector *models.EntitySelector,
	params *models.BuildEntityESQueryParams,
) (json.RawMessage, error) {
	return transport.CallRESResult[json.RawMessage](
		c.baseClient.resClient,
		"guestprofile.entities."+selector.EntityID.String()+".build-elasticsearch-query",
		resprot.Request{
			Params: params,
			Token:  json.RawMessage(`{"agentRoles": ["service"]}`),
		},
	)
}

func (c *entitiesClient) ReadEntity(selector *models.EntitySelector) (*models.Entity, error) {
	return c.readEntityByRID(selector.RID())
}

func (c *entitiesClient) ReadEntityAccounts(selector *models.EntityAccountsSelector) ([]*models.Entity, error) {
	cacheKey := selector.RID() + "?" + selector.EncodedQuery()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(cacheKey).([]*models.Entity); ok {
		return cachedResult, nil
	}

	references, err := transport.GetRESCollection[res.Ref](
		c.baseClient.resClient,
		selector.RID(),
		resprot.Request{
			Query: selector.EncodedQuery(),
		},
	)
	if err != nil {
		return nil, err
	}

	result := []*models.Entity{}

	for _, reference := range references {
		account, err := c.readEntityByRID(string(reference))
		if err != nil {
			return nil, err
		}

		result = append(result, account)
	}

	defer c.baseClient.eventuallyWriteCache(cacheKey, result)

	return result, nil
}

func (c *entitiesClient) ReadAccountParents(selector *models.EntitySelector) ([]*models.Entity, error) {
	cacheKey := selector.RID() + ".parent-entities"

	if cachedResult, ok := c.baseClient.eventuallyReadCache(cacheKey).([]*models.Entity); ok {
		return cachedResult, nil
	}

	entity, err := c.readEntityByRID(selector.RID())
	if err != nil {
		return nil, err
	}

	if entity.Type != models.EntityTypeAccount {
		return nil, fmt.Errorf("entity is not an account")
	}

	result := []*models.Entity{}

	if entity.Chain != "" {
		chain, err := c.readEntityByRID(string(entity.Chain))
		if err != nil {
			return nil, err
		}

		result = append(result, chain)
	}

	if entity.Group != "" {
		group, err := c.readEntityByRID(string(entity.Group))
		if err != nil {
			return nil, err
		}

		result = append(result, group)
	}

	defer c.baseClient.eventuallyWriteCache(cacheKey, result)

	return result, nil
}

func (c *entitiesClient) ReadEntityCustomFields(
	selector *models.EntityCustomFieldsSelector,
) (*models.EntityCustomFields, error) {
	if cachedResult, ok := c.baseClient.eventuallyReadCache(selector.RID()).(*models.EntityCustomFields); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*models.EntityCustomFields](
		c.baseClient.resClient,
		selector.RID(),
		resprot.Request{},
	)
	if err != nil {
		return nil, err
	}

	defer c.baseClient.eventuallyWriteCache(selector.RID(), result)

	return result, nil
}

func (c *entitiesClient) PatchEntity(selector *models.EntitySelector, updates *models.EntityUpdates) error {
	encodedUpdates, err := json.Marshal(updates)
	if err != nil {
		return fmt.Errorf("could not encode updates: %w", err)
	}

	if response := c.baseClient.resClient.Request(
		"call."+selector.RID()+".patch",
		resprot.Request{Params: json.RawMessage(encodedUpdates)},
	); response.HasError() {
		return response.Error
	}

	return nil
}

func (c *entitiesClient) readEntityByRID(resourceID string) (*models.Entity, error) {
	if cachedResult, ok := c.baseClient.eventuallyReadCache(resourceID).(*models.Entity); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*models.Entity](c.baseClient.resClient, resourceID, resprot.Request{})
	if err != nil {
		return nil, err
	}

	defer c.baseClient.eventuallyWriteCache(resourceID, result)

	return result, nil
}
