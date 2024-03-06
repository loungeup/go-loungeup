//nolint:gomnd
package testdata

import (
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/client/models"
)

var (
	RoomType = models.RoomType{
		ID:                   uuid.New(),
		EntityID:             uuid.New(),
		Name:                 "Standard",
		Code:                 "STD",
		Capacity:             100,
		CapacitySafetyMargin: 10,
		CreatedAt:            time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:            time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	RoomTypeCollection = `[
		{"rid": "` + RoomTypeSelector.RID() + `"}
	]`

	RoomTypeModel = `{
		"id": "` + RoomType.ID.String() + `",
		"entityId": "` + RoomType.EntityID.String() + `",
		"name": "Standard",
		"code": "STD",
		"capacity": 100,
		"capacitySafetyMargin": 10,
		"createdAt": "2020-01-01T00:00:00Z",
		"updatedAt": "2020-01-01T00:00:00Z"
	}`

	RoomTypeSelector = models.RoomTypeSelector{
		EntityID:      RoomType.EntityID,
		IntegrationID: RoomType.ID,
	}

	RoomTypesSelector = models.RoomTypesSelector{
		EntityID: RoomType.EntityID,
	}
)
