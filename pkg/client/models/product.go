package models

import (
	"net/url"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/translations"
)

type Product struct {
	ID            string                               `json:"id"`
	EntityID      string                               `json:"entityId"`
	Category      string                               `json:"category"`
	Name          DataValue[translations.Translations] `json:"name"`
	Description   DataValue[translations.Translations] `json:"description"`
	Image         string                               `json:"image"`
	Configuration DataValue[any]                       `json:"configuration"`
	CreatedAt     string                               `json:"createdAt"`
	DisabledAt    string                               `json:"disabledAt"`
	EnabledAt     string                               `json:"enabledAt"`
	UpdatedAt     string                               `json:"updatedAt"`
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
