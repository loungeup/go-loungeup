package models

import "github.com/google/uuid"

type EntityMetadatasSelector struct {
	EntityID uuid.UUID
}

func (s *EntityMetadatasSelector) RID() string {
	return "proxy-db.entities." + s.EntityID.String() + ".metadata"
}

type EntityMetadatas struct {
	Metadata map[string]any `json:"metadata"`
}
