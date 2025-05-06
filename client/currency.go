package client

import (
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/transport"
)

//go:generate mockgen -source currency.go -destination=./mocks/mock_currency.go -package=mocks

type CurrenciesManager interface {
	ReadCurrencyRates(selector *models.CurrencyRatesSelector) (*models.CurrencyRates, error)
}

type CurrenciesClient struct {
	base *BaseClient
}

func NewCurrenciesClient(base *BaseClient) *CurrenciesClient {
	return &CurrenciesClient{base}
}

func (c *CurrenciesClient) ReadCurrencyRates(selector *models.CurrencyRatesSelector) (*models.CurrencyRates, error) {
	cacheKey := selector.RID() + "?" + selector.EncodedQuery()

	if cachedResult, ok := c.base.ReadCache(cacheKey).(*models.CurrencyRates); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*models.CurrencyRates](
		c.base.resClient,
		selector.RID(),
		resprot.Request{Query: selector.EncodedQuery()},
	)
	if err != nil {
		return nil, err
	}

	defer c.base.WriteCache(cacheKey, result)

	return result, nil
}
