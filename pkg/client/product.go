package client

import (
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

type productsClient struct{ baseClient *Client }

func (c *productsClient) ReadProducts(selector *models.ProductsSelector) ([]*models.Product, error) {
	cacheKey := selector.RID() + "?" + selector.EncodedQuery()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(cacheKey).([]*models.Product); ok {
		return cachedResult, nil
	}

	references, err := transport.GetRESCollection[res.Ref](c.baseClient.resClient, selector.RID(), resprot.Request{})
	if err != nil {
		return nil, err
	}

	result := []*models.Product{}

	for _, reference := range references {
		product, err := c.readProductByRID(string(reference))
		if err != nil {
			return nil, err
		}

		result = append(result, product)
	}

	defer c.baseClient.eventuallyWriteCache(cacheKey, result)

	return result, nil
}

func (c *productsClient) ReadProduct(selector *models.ProductSelector) (*models.Product, error) {
	return c.readProductByRID(selector.RID())
}

func (c *productsClient) readProductByRID(rid string) (*models.Product, error) {
	if cachedResult, ok := c.baseClient.eventuallyReadCache(rid).(*models.Product); ok {
		return cachedResult, nil
	}

	product, err := transport.GetRESModel[*models.Product](c.baseClient.resClient, rid)
	if err != nil {
		return nil, err
	}

	defer c.baseClient.eventuallyWriteCache(rid, product)

	return product, nil
}
