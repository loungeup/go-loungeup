package models

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
