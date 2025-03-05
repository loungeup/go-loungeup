package testdata

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/client/models"
)

var (
	EntityIntegration = &models.EntityIntegration{
		ID:                   uuid.New(),
		EntityID:             uuid.New(),
		IntegrationReference: res.Ref(IntegrationSelector.RID()),
		Integration:          Integration,
		Values: models.NewDataValue(models.EntityIntegrationValues{
			"username": "john.doe",
		}),
		Status: "valid",
	}

	EntityIntegrationCollection = `[
		{"rid": "` + EntityIntegrationSelector.RID() + `"}
	]`

	EntityIntegrationModel = `{
		"id": "` + EntityIntegration.ID.String() + `",
		"entityId": "` + EntityIntegration.EntityID.String() + `",
		"integration": {"rid": "` + IntegrationSelector.RID() + `"},
		"values": {
			"data": {
				"username": "john.doe"
			}
		},
		"status": "valid"
	}`

	EntityIntegrationSelector = &models.EntityIntegrationSelector{
		EntityID:      EntityIntegration.EntityID,
		IntegrationID: EntityIntegration.ID,
	}

	EntityIntegrationsSelector = &models.EntityIntegrationsSelector{
		EntityID:             EntityIntegration.EntityID,
		IntegrationsSelector: IntegrationsSelector,
	}

	Integration = &models.Integration{
		Name:     "mews",
		Category: "pms",
		Unique:   true,
		Params: models.NewDataValue(models.IntegrationParams{
			{
				Name:     "clientSecret",
				Type:     "string",
				Format:   "password",
				Required: true,
			},
		}),
		Provider: models.NewDataValue(models.IntegrationProvider{
			Name: "mews",
			Properties: map[string]any{
				"matchedBookingFields": []any{"arrival", "departure"},
			},
		}),
	}

	IntegrationCollection = `[
		{"rid": "` + IntegrationSelector.RID() + `"}
	]`

	IntegrationModel = `{
		"name": "mews",
		"category": "pms",
		"unique": true,
		"parameters": {
			"data": [
				{
					"name": "clientSecret",
					"type": "string",
					"format": "password",
					"required": true
				}
			]
		},
		"provider": {
			"data": {
				"name": "mews",
				"properties": {
					"matchedBookingFields": ["arrival", "departure"]
				}
			}
		}
	}`

	IntegrationSelector = &models.IntegrationSelector{
		Name: "mews",
	}

	IntegrationsSelector = &models.IntegrationsSelector{
		Category: "pms",
	}

	ProviderResult      = json.RawMessage(`{"foo": "bar"}`)
	ProviderResultModel = `{"foo": "bar"}`
)
