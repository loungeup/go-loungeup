package resresultsets

import (
	"github.com/google/uuid"
)

type resultSet struct {
	serviceName string
	id          uuid.UUID
	collection  any
}

func (s *resultSet) rid() string { return s.serviceName + ".result-sets." + s.id.String() }
