package testdata

import (
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/translations"
)

var (
	Product = &models.Product{
		ID:       uuid.MustParse("56872372-7e3d-412a-a3b7-42c8c181147c"),
		EntityID: uuid.MustParse("c9503e44-9551-422e-88c9-bd3e0037d27a"),
		Category: "upgrade",
		Name: models.NewDataValue(translations.Translations{
			"en": "A name",
		}),
		Description: models.NewDataValue(translations.Translations{
			"en": "A description",
		}),
		Image:         "https://example.com/image.jpg",
		Configuration: models.NewDataValue[any]("foo"),
		CreatedAt:     time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		DisabledAt:    time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		EnabledAt:     time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	ProductsCollection = `[
		{
			"rid": "bookings-manager.entities.c9503e44-9551-422e-88c9-bd3e0037d27a.products.56872372-7e3d-412a-a3b7-42c8c181147c"
		}
	]`
	ProductModel = `{
		"id": "56872372-7e3d-412a-a3b7-42c8c181147c",
		"entityId": "c9503e44-9551-422e-88c9-bd3e0037d27a",
		"category": "upgrade",
		"name": {
			"data": {
				"en": "A name"
			}
		},
		"description": {
			"data": {
				"en": "A description"
			}
		},
		"image": "https://example.com/image.jpg",
		"configuration": {
			"data": "foo"
		},
		"createdAt": "2023-01-01T00:00:00Z",
		"disabledAt": "2023-01-01T00:00:00Z",
		"enabledAt": "2023-01-01T00:00:00Z",
		"updatedAt": "2023-01-01T00:00:00Z"
	}`

	ProductSelector = &models.ProductSelector{
		EntityID:  uuid.MustParse("c9503e44-9551-422e-88c9-bd3e0037d27a"),
		ProductID: uuid.MustParse("56872372-7e3d-412a-a3b7-42c8c181147c"),
	}
	ProductsSelector = &models.ProductsSelector{
		EntityID: uuid.MustParse("c9503e44-9551-422e-88c9-bd3e0037d27a"),
	}
)
