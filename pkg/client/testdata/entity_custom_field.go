package testdata

import "github.com/loungeup/go-loungeup/pkg/client/models"

var (
	EntityCustomFields = &models.EntityCustomFields{
		User: models.NewDataValue(map[string]models.EntityCustomField{
			"vipLevel": {
				Label: "VIP Level",
				Type:  models.EntityCustomFieldTypeText,
			},
		}),
		Visit: models.NewDataValue(map[string]models.EntityCustomField{
			"adultsCounts": {
				Label: "Adults Counts",
				Type:  models.EntityCustomFieldTypeNumber,
			},
		}),
	}

	EntityCustomFieldsModel = `{
		"user": {
			"data": {
				"vipLevel": {
					"label": "VIP Level",
					"type": "text"
				}
			}
		},
		"visit": {
			"data": {
				"adultsCounts": {
					"label": "Adults Counts",
					"type": "number"
				}
			}
		}
	}`

	EntityCustomFieldsSelector = &models.EntityCustomFieldsSelector{
		EntityID: Entity.ID,
	}
)
