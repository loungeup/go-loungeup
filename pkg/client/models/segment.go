package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type SegmentSelector struct {
	EntityID  uuid.UUID
	SegmentID uuid.UUID
}

func (s *SegmentSelector) RID() string {
	return fmt.Sprintf("guestprofile.entities.%s.segments.%s.build-elasticsearch-query",
		s.EntityID.String(),
		s.SegmentID.String(),
	)
}

type SegmentQuery struct {
	Criteria []*SegmentCriteria `json:"criteria"`
}

type SegmentCriteria struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

type BuildESQueryResponse struct {
	Query json.RawMessage `json:"query"`
}
