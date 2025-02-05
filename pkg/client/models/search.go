package models

import (
	"time"

	"github.com/google/uuid"
)

type SearchConditions struct {
	Logic    string             `json:"logic"`
	Criteria []*SearchCriterion `json:"criteria"`
}

type SearchCriterion struct {
	Logic    string            `json:"logic"`
	Criteria []*SearchCriteria `json:"criteria"`
}

type SearchCriteria struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

type StructuredValueCriterion[T any] struct {
	Preferred bool      `json:"preferred,omitempty"`
	Value     T         `json:"value,omitempty"`
	UpdatedAt time.Time `json:"updated,omitempty"`
	From      string    `json:"from,omitempty"`
}

type PMSIDCriterion struct {
	ID  string `json:"id,omitempty"`
	PMS string `json:"pms,omitempty"`
}
type MessengerIDCriterion struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type PhoneCriterion struct {
	Phone       string `json:"phone,omitempty"`
	Type        string `json:"type,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
}

type SearchByContactMatch string

const (
	SearchByContactMatchOne SearchByContactMatch = "one"
	SearchByContactMatchAll SearchByContactMatch = "all"
)

type SearchByContactSelector struct {
	EntityID uuid.UUID `json:"-"`

	Match SearchByContactMatch `json:"match,omitempty"`

	Emails      []StructuredValueCriterion[Email]              `json:"emails,omitempty"`
	Phones      []StructuredValueCriterion[PhoneCriterion]     `json:"phones,omitempty"`
	MessengerID StructuredValueCriterion[MessengerIDCriterion] `json:"messengerId,omitempty"`
	PMSID       StructuredValueCriterion[PMSIDCriterion]       `json:"pmsId,omitempty"`
	Credentials StructuredValueCriterion[Credentials]          `json:"credentials,omitempty"`

	LastKey string `json:"lastKey,omitempty"`
	Size    int    `json:"size,omitempty"`
}
