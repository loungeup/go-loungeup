package resmodels

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/errors"
)

var (
	ErrorEntityIntegrationValueNotFound = &errors.Error{
		Code:    errors.CodeInvalid,
		Message: "Entity integration value not found",
	}

	ErrorInvalidEntityIntegrationValueType = &errors.Error{
		Code:    errors.CodeInvalid,
		Message: "Invalid entity integration value type",
	}
)

type EntityIntegration struct {
	ID                   uuid.UUID                              `json:"id,omitempty"`
	EntityID             uuid.UUID                              `json:"entityId,omitempty"`
	IntegrationReference res.Ref                                `json:"integration,omitempty"`
	Integration          *models.Integration                    `json:"-"`
	Values               res.DataValue[EntityIntegrationValues] `json:"values,omitempty"`
	Status               string                                 `json:"status,omitempty"`
}

type EntityIntegrationValues map[string]any

func GetEntityIntegrationValue[T any](values EntityIntegrationValues, key string) (T, error) {
	var result T

	value, ok := values[key]
	if !ok {
		return result, ErrorEntityIntegrationValueNotFound
	}

	encodedValue, err := json.Marshal(value)
	if err != nil {
		return result, ErrorInvalidEntityIntegrationValueType
	}

	if err := json.Unmarshal(encodedValue, &result); err != nil {
		return result, ErrorInvalidEntityIntegrationValueType
	}

	return result, nil
}

type EntityIntegrationSelector struct {
	EntityID, IntegrationID uuid.UUID
}

func (s EntityIntegrationSelector) RID() string {
	return "integrations.entities." + s.EntityID.String() + ".integrations." + s.IntegrationID.String()
}

type EntityIntegrationsSelector struct {
	*models.IntegrationsSelector

	EnabledFeatures []string
	EntityID        uuid.UUID
	Matchers        map[string][]string
}

func (s EntityIntegrationsSelector) EncodedQuery() string {
	query := "category=" + s.Category + "&enabledFeatures=" + strings.Join(s.EnabledFeatures, ",")

	if len(s.Matchers) > 0 {
		encoded, _ := json.Marshal(s.Matchers)
		if len(encoded) > 0 {
			query += "&matchers=" + string(encoded)
		}
	}

	return query
}

func (s EntityIntegrationsSelector) RID() string {
	return "integrations.entities." + s.EntityID.String() + ".integrations"
}

type LatestIntegrationSelector struct {
	*EntityIntegrationsSelector
}

func (s LatestIntegrationSelector) RID() string {
	return "integrations.entities." + s.EntityID.String() + ".latest-integration"
}
