package testdata

import (
	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/resmodels"
)

var (
	EntityCustomFields = &resmodels.EntityCustomFields{
		User: res.NewDataValue(map[string]resmodels.EntityCustomField{
			"vipLevel": {
				Label: "VIP Level",
				Type:  resmodels.EntityCustomFieldTypeText,
			},
		}),
		Visit: res.NewDataValue(map[string]resmodels.EntityCustomField{
			"adultsCounts": {
				Label: "Adults Counts",
				Type:  resmodels.EntityCustomFieldTypeNumber,
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

	EntityCustomFieldsSelector = &resmodels.EntityCustomFieldsSelector{
		EntityID: uuid.MustParse(Entity.ID),
	}
)
