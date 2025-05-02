package client

import (
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/transport"
)

//go:generate mockgen -source product.go -destination=./mocks/mock_product.go -package=mocks

type ProductsManager interface {
	ReadProducts(selector *models.ProductsSelector) ([]*models.Product, error)
	ReadProduct(selector *models.ProductSelector) (*models.Product, error)
}

type ProductsClient struct {
	base *BaseClient
}

func NewProductsClient(base *BaseClient) *ProductsClient {
	return &ProductsClient{
		base: base,
	}
}

func (c *ProductsClient) ReadProducts(selector *models.ProductsSelector) ([]*models.Product, error) {
	cacheKey := selector.RID() + "?" + selector.EncodedQuery()

	if cachedResult, ok := c.base.ReadCache(cacheKey).([]*models.Product); ok {
		return cachedResult, nil
	}

	references, err := transport.GetRESCollection[res.Ref](c.base.resClient, selector.RID(), resprot.Request{
		Query: selector.EncodedQuery(),
	})
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

	defer c.base.WriteCache(cacheKey, result)

	return result, nil
}

func (c *ProductsClient) ReadProduct(selector *models.ProductSelector) (*models.Product, error) {
	return c.readProductByRID(selector.RID())
}

func (c *ProductsClient) readProductByRID(rid string) (*models.Product, error) {
	if cachedResult, ok := c.base.ReadCache(rid).(*models.Product); ok {
		return cachedResult, nil
	}

	product, err := transport.GetRESModel[*models.Product](c.base.resClient, rid, resprot.Request{})
	if err != nil {
		return nil, err
	}

	defer c.base.WriteCache(rid, product)

	return product, nil
}
