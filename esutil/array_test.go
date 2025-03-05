package esutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayMarshalJSON(t *testing.T) {
	tests := map[string]struct {
		in   Array[string]
		want []byte
	}{
		"zero values": {
			in:   Array[string]{},
			want: []byte("[]"),
		},
		"one value": {
			in:   Array[string]{"foo"},
			want: []byte(`"foo"`),
		},
		"many values": {
			in:   Array[string]{"foo", "bar"},
			want: []byte(`["foo","bar"]`),
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			got, err := json.Marshal(tt.in)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArrayUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		in   []byte
		want Array[string]
	}{
		"zero values": {
			in:   []byte("[]"),
			want: Array[string]{},
		},
		"one value": {
			in:   []byte(`"foo"`),
			want: Array[string]{"foo"},
		},
		"many values": {
			in:   []byte(`["foo","bar"]`),
			want: Array[string]{"foo", "bar"},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			var got Array[string]

			assert.NoError(t, json.Unmarshal(tt.in, &got))
			assert.Equal(t, tt.want, got)
		})
	}
}
