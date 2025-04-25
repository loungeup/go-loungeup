package client

import (
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/transport"
)

type proxyDBClient struct{ baseClient *Client }

func (client *proxyDBClient) ReadBooking(selector *models.BookingSelector) (*models.Booking, error) {
	return transport.GetRESModel[*models.Booking](client.baseClient.resClient, selector.RID(), resprot.Request{})
}

func (client *proxyDBClient) ReadBookingById(selector *models.BookingSelectorById) (*models.Booking, error) {
	return transport.GetRESModel[*models.Booking](client.baseClient.resClient, selector.RID(), resprot.Request{})
}

func (client *proxyDBClient) ReadEntityMetadatas(selector *models.EntityMetadatasSelector) (*models.EntityMetadatas, error) {
	return transport.GetRESModel[*models.EntityMetadatas](client.baseClient.resClient, selector.RID(), resprot.Request{})
}
