package matcher

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMatcherGet(t *testing.T) {
	matcher := Matcher[any]{
		true:   NewMatchable("test"),
		1.0:    NewMatchable("test"),
		1:      NewMatchable("test"),
		"test": NewMatchable(1),
		uuid.MustParse("cb612eba-0206-4b6e-b5a3-68bff833a7b5"): NewMatchable("test"),
	}

	tests := map[string]struct {
		in   any
		want *Matchable
	}{
		"bool":      {in: true, want: NewMatchable("test")},
		"float":     {in: 1.0, want: NewMatchable("test")},
		"int":       {in: 1, want: NewMatchable("test")},
		"nil":       {in: nil, want: nil},
		"string":    {in: "test", want: NewMatchable(1)},
		"uuid":      {in: uuid.MustParse("cb612eba-0206-4b6e-b5a3-68bff833a7b5"), want: NewMatchable("test")},
		"not found": {in: "not found", want: nil},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, matcher.Get(tt.in))
		})
	}
}

func TestMatchableMarshalJSON(t *testing.T) {
	tests := map[string]struct {
		in   *Matchable
		want string
	}{
		"bool":   {in: NewMatchable(true), want: "true"},
		"float":  {in: NewMatchable(1.0), want: "1"},
		"int":    {in: NewMatchable(1), want: "1"},
		"nil":    {in: NewMatchable(nil), want: "null"},
		"string": {in: NewMatchable("test"), want: `"test"`},
		"uuid":   {in: NewMatchable(uuid.MustParse("cb612eba-0206-4b6e-b5a3-68bff833a7b5")), want: `"cb612eba-0206-4b6e-b5a3-68bff833a7b5"`},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			got, err := json.Marshal(tt.in)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestMatchableUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		in   json.RawMessage
		want *Matchable
	}{
		"simple": {in: json.RawMessage(`true`), want: NewMatchable(true)},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			got := NewMatchable(nil)
			assert.NoError(t, json.Unmarshal(tt.in, got))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMatchableIsNil(t *testing.T) {
	tests := map[string]struct {
		in   *Matchable
		want bool
	}{
		"nil":     {in: NewMatchable(nil), want: true},
		"not nil": {in: NewMatchable(true), want: false},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.in.IsNil())
		})
	}
}

func TestCastMatchable(t *testing.T) {
	tests := map[string]struct {
		in     *Matchable
		assert func(t *testing.T, in *Matchable)
	}{
		"bool": {
			in:     NewMatchable(true),
			assert: func(t *testing.T, in *Matchable) { assert.Equal(t, true, in.Bool()) },
		},
		"float64": {
			in:     NewMatchable(1.0),
			assert: func(t *testing.T, in *Matchable) { assert.Equal(t, 1.0, in.Float64()) },
		},
		"int": {
			in:     NewMatchable(1),
			assert: func(t *testing.T, in *Matchable) { assert.Equal(t, 1, in.Int()) },
		},
		"string": {
			in:     NewMatchable("test"),
			assert: func(t *testing.T, in *Matchable) { assert.Equal(t, "test", in.String()) },
		},
		"uuid": {
			in: NewMatchable(uuid.MustParse("cb612eba-0206-4b6e-b5a3-68bff833a7b5")),
			assert: func(t *testing.T, in *Matchable) {
				assert.Equal(t, uuid.MustParse("cb612eba-0206-4b6e-b5a3-68bff833a7b5"), in.UUID())
			},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			tt.assert(t, tt.in)
		})
	}
}
