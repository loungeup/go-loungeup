package uuidutil

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLast(t *testing.T) {
	tests := map[string]struct {
		in   uuid.UUIDs
		want uuid.UUID
	}{
		"empty": {
			in:   uuid.UUIDs{},
			want: uuid.Nil,
		},
		"simple": {
			in: uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
			want: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, Last(tt.in))
		})
	}
}
