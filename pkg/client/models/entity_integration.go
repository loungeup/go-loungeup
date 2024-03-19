package models

import (
	"strings"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
)

type EntityIntegration struct {
	ID                   uuid.UUID                 `json:"id,omitempty"`
	EntityID             uuid.UUID                 `json:"entityId,omitempty"`
	IntegrationReference res.Ref                   `json:"integration,omitempty"`
	Integration          *Integration              `json:"-"`
	Values               DataValue[map[string]any] `json:"values,omitempty"`
	Status               string                    `json:"status,omitempty"`
}

type EntityIntegrationSelector struct {
	EntityID, IntegrationID uuid.UUID
}

func (s EntityIntegrationSelector) RID() string {
	return "integrations.entities." + s.EntityID.String() + ".integrations." + s.IntegrationID.String()
}

type EntityIntegrationsSelector struct {
	*IntegrationsSelector

	EnabledFeatures []string
	EntityID        uuid.UUID
}

func (s EntityIntegrationsSelector) EncodedQuery() string {
	return "category=" + s.Category + "&enabledFeatures=" + strings.Join(s.EnabledFeatures, ",")
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
