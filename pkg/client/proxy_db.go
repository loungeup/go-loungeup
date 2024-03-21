package client

import (
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

type proxyDBClient struct{ baseClient *Client }

func (c *proxyDBClient) ReadBooking(selector *models.BookingSelector) (*models.Booking, error) {
	return transport.GetRESModel[*models.Booking](c.baseClient.resClient, selector.RID())
}
