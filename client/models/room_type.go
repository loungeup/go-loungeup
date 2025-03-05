package models

import (
	"time"

	"github.com/google/uuid"
)

type RoomType struct {
	ID                   uuid.UUID `json:"id"`
	EntityID             uuid.UUID `json:"entityId"`
	Name                 string    `json:"name,omitempty"`
	Code                 string    `json:"code,omitempty"`
	Capacity             int       `json:"capacity,omitempty"`
	CapacitySafetyMargin int       `json:"capacitySafetyMargin,omitempty"`
	CreatedAt            time.Time `json:"createdAt,omitempty"`
	UpdatedAt            time.Time `json:"updatedAt,omitempty"`
}

func (t *RoomType) EffectiveCapacity() int { return t.Capacity - t.CapacitySafetyMargin }

func (t *RoomType) IsAvailable() bool { return t.EffectiveCapacity() > 0 }

type RoomTypeSelector struct {
	EntityID      uuid.UUID
	IntegrationID uuid.UUID
}

func (s RoomTypeSelector) RID() string {
	return "authority.entities." + s.EntityID.String() + ".room-types." + s.IntegrationID.String()
}

type RoomTypesSelector struct {
	EntityID uuid.UUID
}

func (s RoomTypesSelector) RID() string {
	return "authority.entities." + s.EntityID.String() + ".room-types"
}
