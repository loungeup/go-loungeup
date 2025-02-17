package urlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeJoinPath(t *testing.T) {
	tests := map[string]struct {
		base     string
		elements []string
		want     string
	}{
		"simple": {
			base:     "https://example.com",
			elements: []string{"path1", "path2"},
			want:     "https://example.com/path1/path2",
		},
		"only base": {
			base:     "https://example.com",
			elements: []string{},
			want:     "https://example.com/",
		},
		"empty base": {
			base:     "",
			elements: []string{"path1", "path2"},
			want:     "path1/path2",
		},
		"invalid base": {
			base:     "1234://",
			elements: []string{"path1", "path2"},
			want:     "path1/path2",
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, SafeJoinPath(tt.base, tt.elements...))
		})
	}
}
