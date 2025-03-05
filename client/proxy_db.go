package client

import (
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/transport"
)

type proxyDBClient struct{ baseClient *Client }

func (c *proxyDBClient) ReadBooking(selector *models.BookingSelector) (*models.Booking, error) {
	return transport.GetRESModel[*models.Booking](c.baseClient.resClient, selector.RID(), resprot.Request{})
}
