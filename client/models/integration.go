package models

type Integration struct {
	Name     string                         `json:"name,omitempty"`
	Category string                         `json:"category,omitempty"`
	Unique   bool                           `json:"unique,omitempty"`
	Params   DataValue[IntegrationParams]   `json:"parameters,omitempty"`
	Provider DataValue[IntegrationProvider] `json:"provider,omitempty"`
}

type IntegrationParams []IntegrationParam

type IntegrationParam struct {
	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	Type        string               `json:"type,omitempty"`
	Format      string               `json:"format,omitempty"`
	ReadOnly    bool                 `json:"readOnly,omitempty"`
	Required    bool                 `json:"required,omitempty"`
	Default     any                  `json:"default,omitempty"`
	Enum        IntegrationParamEnum `json:"enum,omitempty"`
	Items       *IntegrationParam    `json:"items,omitempty"`
}

type IntegrationParamEnum []IntegrationParamEnumValue

type IntegrationParamEnumValue struct {
	// Key is machine-readable.
	Key any `json:"key"`

	// Value is human-readable.
	Value string `json:"value,omitempty"`
}

type IntegrationProvider struct {
	Name       string         `json:"name,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

type IntegrationSelector struct {
	Name string
}

func (s IntegrationSelector) RID() string { return "integrations.integrations." + s.Name }

type IntegrationsSelector struct {
	Category string
}

func (s IntegrationsSelector) EncodedQuery() string { return "category=" + s.Category }
func (s IntegrationsSelector) RID() string          { return "integrations.integrations" }
