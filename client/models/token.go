package models

import (
	"slices"

	"github.com/google/uuid"
)

type Token struct {
	AgentID    uuid.UUID           `json:"agentId"`
	AgentRoles TokenAgentRoleSlice `json:"agentRoles"`
}

type TokenAgentRoleSlice []TokenAgentRole

func (roles TokenAgentRoleSlice) ContainsGlobalRole() bool {
	globalRoles := []TokenAgentRole{
		TokenAgentRoleDeveloper,
		TokenAgentRoleService,
		TokenAgentRoleStaff,
	}

	for _, role := range roles {
		if slices.Contains(globalRoles, role) {
			return true
		}
	}

	return false
}

type TokenAgentRole string

const (
	TokenAgentRoleUnknown   TokenAgentRole = ""
	TokenAgentRoleAgent     TokenAgentRole = "agent"
	TokenAgentRoleDeveloper TokenAgentRole = "developer"
	TokenAgentRoleService   TokenAgentRole = "service"
	TokenAgentRoleStaff     TokenAgentRole = "staff"
)
