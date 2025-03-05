package uuidutil

import "github.com/google/uuid"

func Last(s uuid.UUIDs) uuid.UUID {
	if len(s) == 0 {
		return uuid.Nil
	}

	return s[len(s)-1]
}
