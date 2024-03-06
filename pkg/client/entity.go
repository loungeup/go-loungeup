package client

import (
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

// entitiesClient provides methods to interact with entities.
type entitiesClient struct{ baseClient *Client }

func (c *entitiesClient) ReadEntity(selector models.EntitySelector) (models.Entity, error) {
	return c.readEntityByRID(selector.RID())
}

func (c *entitiesClient) ReadEntityAccounts(selector models.EntityAccountsSelector) ([]models.Entity, error) {
	cacheKey := selector.RID() + "?" + selector.EncodedQuery()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(cacheKey).([]models.Entity); ok {
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

	result := []models.Entity{}

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

func (c *entitiesClient) ReadEntityCustomFields(
	selector models.EntityCustomFieldsSelector,
) (models.EntityCustomFields, error) {
	if cachedResult, ok := c.baseClient.eventuallyReadCache(selector.RID()).(models.EntityCustomFields); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[models.EntityCustomFields](
		c.baseClient.resClient,
		selector.RID(),
	)
	if err != nil {
		return models.EntityCustomFields{}, err
	}

	defer c.baseClient.eventuallyWriteCache(selector.RID(), result)

	return result, nil
}

func (c *entitiesClient) readEntityByRID(resourceID string) (models.Entity, error) {
	if cachedResult, ok := c.baseClient.eventuallyReadCache(resourceID).(models.Entity); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[models.Entity](c.baseClient.resClient, resourceID)
	if err != nil {
		return models.Entity{}, err
	}

	defer c.baseClient.eventuallyWriteCache(resourceID, result)

	return result, nil
}
