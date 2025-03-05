package models

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/translations"
)

type Product struct {
	ID            uuid.UUID                            `json:"id"`
	EntityID      uuid.UUID                            `json:"entityId"`
	Category      string                               `json:"category"`
	Name          DataValue[translations.Translations] `json:"name"`
	Description   DataValue[translations.Translations] `json:"description"`
	Image         string                               `json:"image"`
	Configuration DataValue[any]                       `json:"configuration"`
	CreatedAt     time.Time                            `json:"createdAt"`
	DisabledAt    time.Time                            `json:"disabledAt"`
	EnabledAt     time.Time                            `json:"enabledAt"`
	UpdatedAt     time.Time                            `json:"updatedAt"`
}

type ProductSelector struct {
	EntityID  uuid.UUID
	ProductID uuid.UUID
}

func (s *ProductSelector) RID() string {
	return "bookings-manager.entities." + s.EntityID.String() + ".products." + s.ProductID.String()
}

type ProductsSelector struct {
	EntityID uuid.UUID
	Category string
	Enabled  string

	RoomTypeIDOrSourceRoomTypeID uuid.UUID
}

func (s *ProductsSelector) RID() string {
	return "bookings-manager.entities." + s.EntityID.String() + ".products"
}

func (s *ProductsSelector) EncodedQuery() string {
	query := url.Values{}

	if s.Category != "" {
		query.Set("category", s.Category)
	}

	if s.Enabled != "" {
		query.Set("enabled", s.Enabled)
	}

	if s.RoomTypeIDOrSourceRoomTypeID != uuid.Nil {
		query.Set("roomTypeIdOrSourceRoomTypeId", s.RoomTypeIDOrSourceRoomTypeID.String())
	}

	return query.Encode()
}
