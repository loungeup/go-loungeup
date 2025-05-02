package resmodels

import (
	"github.com/google/uuid"
	"github.com/jirenius/go-res"
)

type EntityCustomFields struct {
	User  res.DataValue[map[string]EntityCustomField] `json:"user,omitempty"`
	Visit res.DataValue[map[string]EntityCustomField] `json:"visit,omitempty"`
}

type EntityCustomField struct {
	Label string                `json:"label,omitempty"`
	Type  EntityCustomFieldType `json:"type,omitempty"`
}

type EntityCustomFieldType string

const (
	EntityCustomFieldTypeBoolean EntityCustomFieldType = "boolean"
	EntityCustomFieldTypeDate    EntityCustomFieldType = "date"
	EntityCustomFieldTypeList    EntityCustomFieldType = "list"
	EntityCustomFieldTypeNumber  EntityCustomFieldType = "number"
	EntityCustomFieldTypeText    EntityCustomFieldType = "text"
)

type EntityCustomFieldsSelector struct {
	EntityID uuid.UUID
}

func (s EntityCustomFieldsSelector) RID() string {
	return "proxy-db.entities." + s.EntityID.String() + ".custom-fields"
}
