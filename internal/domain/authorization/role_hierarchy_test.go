package authorization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRoleLevel_KnownRoles(t *testing.T) {
	tests := []struct {
		role     string
		expected int
	}{
		{"USER", 1},
		{"DRIVER", 2},
		{"ADMIN", 3},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			level := GetRoleLevel(tt.role)
			assert.Equal(t, tt.expected, level)
		})
	}
}

func TestGetRoleLevel_UnknownRole(t *testing.T) {
	level := GetRoleLevel("SUPERADMIN")
	assert.Equal(t, 0, level)
}

func TestGetRoleLevel_EmptyString(t *testing.T) {
	level := GetRoleLevel("")
	assert.Equal(t, 0, level)
}

func TestGetRoleLevel_CaseSensitive(t *testing.T) {
	level := GetRoleLevel("user")
	assert.Equal(t, 0, level, "lowercase 'user' should not match 'USER'")
}

func TestRoleHierarchy_Ordering(t *testing.T) {
	assert.Greater(t, GetRoleLevel("ADMIN"), GetRoleLevel("DRIVER"),
		"ADMIN should have higher level than DRIVER")
	assert.Greater(t, GetRoleLevel("DRIVER"), GetRoleLevel("USER"),
		"DRIVER should have higher level than USER")
	assert.Greater(t, GetRoleLevel("ADMIN"), GetRoleLevel("USER"),
		"ADMIN should have higher level than USER")
}

func TestRoleHierarchy_AllKnownRolesPresent(t *testing.T) {
	expectedRoles := []string{"USER", "DRIVER", "ADMIN"}
	for _, role := range expectedRoles {
		assert.Contains(t, RoleHierarchy, role)
		assert.Greater(t, RoleHierarchy[role], 0,
			"role %s should have positive level", role)
	}
}

func TestRoleHierarchy_ExactlyThreeRoles(t *testing.T) {
	assert.Len(t, RoleHierarchy, 3, "should have exactly 3 roles defined")
}
