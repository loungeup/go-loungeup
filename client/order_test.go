package client

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pointer"
	"github.com/loungeup/go-loungeup/resmodels"
	"github.com/loungeup/go-loungeup/transport"
	"github.com/loungeup/go-loungeup/transporttest"
	"github.com/stretchr/testify/require"
)

func TestReadOneOrder(t *testing.T) {
	testOrder := &resmodels.Order{
		ID:              uuid.MustParse("214d56f8-b539-465d-aea0-117043c56399"),
		EntityID:        uuid.MustParse("f718f848-bf81-4f08-a692-53b60884cfc0"),
		LegacyBookingID: 1234,
		ProductID:       uuid.MustParse("238c5dd1-ba24-4965-93c0-4e0e8cc33bd8"),
		Price:           49.99,
		Quantity:        12,
		ConvertedPrice: pointer.From(res.NewDataValue(map[string]float64{
			"fr": 39.98,
		})),
		Metadata:    res.NewDataValue(json.RawMessage(`{"foo": "bar"}`)),
		CreatedAt:   time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		CompletedAt: pointer.From(time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC)),
		RunAt:       time.Date(2025, time.January, 3, 0, 0, 0, 0, time.UTC),
	}

	client := NewWithTransport(
		transport.New(
			transport.WithRESClient(&transporttest.RESClientMock{
				RequestFunc: func(subject string, request resprot.Request) resprot.Response {
					expectedSubject := "get." + makeOrderRID(testOrder.EntityID, testOrder.LegacyBookingID, testOrder.ID)
					require.Equal(t, expectedSubject, subject)

					require.Empty(t, request)

					return transporttest.NewRESModelResponse(`{
						"id": "214d56f8-b539-465d-aea0-117043c56399",
						"entityId": "f718f848-bf81-4f08-a692-53b60884cfc0",
						"legacyBookingId": 1234,
						"productId": "238c5dd1-ba24-4965-93c0-4e0e8cc33bd8",
						"price": 49.99,
						"quantity": 12,
						"convertedPrice": {
							"data": {
								"fr": 39.98
							}
						},
						"metadata": {
							"data": {"foo": "bar"}
						},
						"createdAt": "2025-01-01T00:00:00Z",
						"completedAt": "2025-01-02T00:00:00Z",
						"runAt": "2025-01-03T00:00:00Z"
					}`)
				},
			}),
		),
	)

	got, err := client.Orders.ReadOne(&resmodels.OrderSelector{
		EntityID:        testOrder.EntityID,
		LegacyBookingID: testOrder.LegacyBookingID,
		OrderID:         testOrder.ID,
	})
	require.NoError(t, err)
	require.Equal(t, testOrder, got)
}
