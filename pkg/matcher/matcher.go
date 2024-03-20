// Package matcher provides a way to match internal and external IDs.
package matcher

import (
	"strconv"

	"github.com/google/uuid"
)

type Matcher map[Matchable]Matchable

func (m Matcher) Get(key Matchable) Matchable {
	if result, ok := m[key]; ok {
		return result
	}

	return Matchable("")
}

type Matchable string

func (m Matchable) String() string { return string(m) }

func (m Matchable) Bool() bool {
	result, _ := strconv.ParseBool(m.String())
	return result
}

func (m Matchable) Float64() float64 {
	result, _ := strconv.ParseFloat(m.String(), 64)
	return result
}

func (m Matchable) Int() int {
	result, _ := strconv.ParseInt(m.String(), 10, 64)
	return int(result)
}

func (m Matchable) UUID() uuid.UUID {
	result, _ := uuid.Parse(m.String())
	return result
}
