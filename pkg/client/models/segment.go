package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type SegmentSelector struct {
	EntityID  uuid.UUID
	SegmentID uuid.UUID
}

type BuildSegmentESQueryResponse struct {
	Query   json.RawMessage `json:"query"`
	Version string          `json:"version"`
}
