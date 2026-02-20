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
