package resresultsets

import "github.com/google/uuid"

type ResultSet struct {
	ID         uuid.UUID
	Collection any
}
