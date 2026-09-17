package auth

// IsAdminRole 判断角色列表中是否包含 admin。
func IsAdminRole(roles []string) bool {

	for _, role := range roles {
		if role == "admin" {
			return true
		}
	}

	return false
}
