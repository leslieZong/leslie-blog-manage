package constant

// 系统角色名称。
//
// 系统角色属于系统核心权限角色，
// 默认不允许通过普通角色管理 API 删除。
const (
	RoleAdmin = "admin"
)

func IsSystemRole(name string) bool {
	switch name {
	case RoleAdmin:
		return true
	default:
		return false
	}
}
