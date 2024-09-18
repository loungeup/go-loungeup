package client

import (
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

// currencyRate client provides methods to interact with entities.
type currencyClient struct{ baseClient *Client }

func (c *currencyClient) ReadCurrencyRates(selector *models.CurrencyRatesSelector) (*models.CurrencyRates, error) {
	cacheKey := selector.RID() + "?" + selector.EncodedQuery()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(cacheKey).(*models.CurrencyRates); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*models.CurrencyRates](
		c.baseClient.resClient,
		selector.RID(),
		resprot.Request{Query: selector.EncodedQuery()},
	)
	if err != nil {
		return nil, err
	}

	defer c.baseClient.eventuallyWriteCache(cacheKey, result)

	return result, nil
}
