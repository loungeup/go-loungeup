// Package matcher provides a way to match internal and external IDs.
package matcher

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/errors"
)

var InvalidMatchError = &errors.Error{
	Code:    errors.CodeInvalid,
	Message: "A match must have valid internal and external IDs",
}

type Matcher[TKey comparable] map[TKey]*Matchable

func (m Matcher[TKey]) Get(key TKey) *Matchable {
	if result, ok := m[key]; ok {
		return result
	}

	return nil
}

type Matchable struct{ value any }

func NewMatchable(value any) *Matchable { return &Matchable{value} }

var (
	_ json.Marshaler   = (*Matchable)(nil)
	_ json.Unmarshaler = (*Matchable)(nil)
)

func (m *Matchable) MarshalJSON() ([]byte, error)    { return json.Marshal(m.value) }
func (m *Matchable) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, &m.value) }

func (m *Matchable) IsNil() bool { return m == nil || m.value == nil }

func (m *Matchable) Bool() bool       { return castAs[bool](m.value) }
func (m *Matchable) Float64() float64 { return castAs[float64](m.value) }
func (m *Matchable) Int() int         { return castAs[int](m.value) }
func (m *Matchable) String() string   { return castAs[string](m.value) }
func (m *Matchable) UUID() uuid.UUID  { return castAs[uuid.UUID](m.value) }

func castAs[T any](v any) T {
	var result T

	result, _ = v.(T)

	return result
}
