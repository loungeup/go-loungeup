package client

import (
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/transport"
)

//go:generate mockgen -source proxy_db.go -destination=./mocks/mock_proxy_db.go -package=mocks
type ProxyDBManager interface {
	ReadBooking(selector *models.BookingSelector) (*models.Booking, error)
	ReadBookingById(selector *models.BookingSelectorById) (*models.Booking, error)
}

type ProxyDBClient struct {
	base *BaseClient
}

func NewProxyDBClient(base *BaseClient) *ProxyDBClient {
	return &ProxyDBClient{
		base: base,
	}
}

func (c *ProxyDBClient) ReadBooking(selector *models.BookingSelector) (*models.Booking, error) {
	return transport.GetRESModel[*models.Booking](c.base.resClient, selector.RID(), resprot.Request{})
}

func (client *ProxyDBClient) ReadBookingById(selector *models.BookingSelectorById) (*models.Booking, error) {
	return transport.GetRESModel[*models.Booking](client.base.resClient, selector.RID(), resprot.Request{})
}

func (client *ProxyDBClient) ReadEntityMetadatas(selector *models.EntityMetadatasSelector) (*models.EntityMetadatas, error) {
	return transport.GetRESModel[*models.EntityMetadatas](client.base.resClient, selector.RID(), resprot.Request{})
}
