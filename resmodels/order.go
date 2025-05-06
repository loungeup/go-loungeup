package resmodels

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
)

type Order struct {
	ID              uuid.UUID                          `json:"id"`
	EntityID        uuid.UUID                          `json:"entityId"`
	LegacyBookingID uint64                             `json:"legacyBookingId"`
	ProductID       uuid.UUID                          `json:"productId"`
	Price           float64                            `json:"price"`
	Quantity        uint64                             `json:"quantity"`
	ConvertedPrice  *res.DataValue[map[string]float64] `json:"convertedPrice"`
	Metadata        res.DataValue[json.RawMessage]     `json:"metadata"`
	CreatedAt       time.Time                          `json:"createdAt"`
	CompletedAt     *time.Time                         `json:"completedAt"`
	RunAt           time.Time                          `json:"runAt"`
}

type OrderSelector struct {
	EntityID        uuid.UUID
	LegacyBookingID uint64
	OrderID         uuid.UUID
}
