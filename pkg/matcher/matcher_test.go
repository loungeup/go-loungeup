package matcher

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMatcherGet(t *testing.T) {
	matcher := Matcher{
		"test": "test",
	}

	tests := map[string]struct {
		in, want Matchable
	}{
		"found":     {in: "test", want: "test"},
		"not found": {in: "not found", want: ""},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, matcher.Get(tt.in))
		})
	}
}

func TestCastMatchable(t *testing.T) {
	tests := map[string]struct {
		in     Matchable
		assert func(t *testing.T, in Matchable)
	}{
		"bool": {
			in:     "true",
			assert: func(t *testing.T, in Matchable) { assert.Equal(t, true, in.Bool()) },
		},
		"float64": {
			in:     "1.0",
			assert: func(t *testing.T, in Matchable) { assert.Equal(t, 1.0, in.Float64()) },
		},
		"int": {
			in:     "1",
			assert: func(t *testing.T, in Matchable) { assert.Equal(t, 1, in.Int()) },
		},
		"string": {
			in:     "test",
			assert: func(t *testing.T, in Matchable) { assert.Equal(t, "test", in.String()) },
		},
		"uuid": {
			in: "cb612eba-0206-4b6e-b5a3-68bff833a7b5",
			assert: func(t *testing.T, in Matchable) {
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
