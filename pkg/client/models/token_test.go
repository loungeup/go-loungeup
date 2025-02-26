package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenAgentRoleSliceContainsGlobalRole(t *testing.T) {
	tests := map[string]struct {
		roles TokenAgentRoleSlice
		want  bool
	}{
		"empty": {
			roles: TokenAgentRoleSlice{},
			want:  false,
		},
		"with global role": {
			roles: TokenAgentRoleSlice{TokenAgentRoleDeveloper},
			want:  true,
		},
		"without global role": {
			roles: TokenAgentRoleSlice{TokenAgentRoleUnknown, TokenAgentRoleAgent, "foo"},
			want:  false,
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.roles.ContainsGlobalRole())
		})
	}
}
