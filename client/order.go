package client

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/resmodels"
	"github.com/loungeup/go-loungeup/transport"
)

type OrdersManager interface {
	ReadOne(selector *resmodels.OrderSelector) (*resmodels.Order, error)
}

type OrdersClient struct {
	base *BaseClient
}

func NewOrdersClient(base *BaseClient) *OrdersClient { return &OrdersClient{base} }

var _ (OrdersManager) = (*OrdersClient)(nil)

func (client *OrdersClient) ReadOne(selector *resmodels.OrderSelector) (*resmodels.Order, error) {
	return transport.GetRESModel[*resmodels.Order](
		client.base.resClient,
		makeOrderRID(selector.EntityID, selector.LegacyBookingID, selector.OrderID),
		resprot.Request{},
	)
}

func makeOrderRID(entityID uuid.UUID, legacyBookingID uint64, orderID uuid.UUID) string {
	return fmt.Sprintf(
		"bookings-manager.entities.%s.bookings.%d.orders.%s",
		entityID.String(),
		legacyBookingID,
		orderID.String(),
	)
}
