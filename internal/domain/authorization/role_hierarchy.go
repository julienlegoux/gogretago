package authorization

// RoleHierarchy maps role names to numeric privilege levels.
// Higher values indicate greater privileges.
var RoleHierarchy = map[string]int{
	"USER":   1,
	"DRIVER": 2,
	"ADMIN":  3,
}

// GetRoleLevel returns the privilege level for a role (0 if unknown)
func GetRoleLevel(role string) int {
	if level, ok := RoleHierarchy[role]; ok {
		return level
	}
	return 0
}
