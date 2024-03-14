package client

import (
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

type proxyDBClient struct{ baseClient *Client }

func (c *proxyDBClient) ReadBooking(selector *models.BookingSelector) (*models.Booking, error) {
	if cachedResult, ok := c.baseClient.eventuallyReadCache(selector.RID()).(*models.Booking); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*models.Booking](c.baseClient.resClient, selector.RID())
	if err != nil {
		return nil, err
	}

	defer c.baseClient.eventuallyWriteCache(selector.RID(), result)

	return result, nil
}
