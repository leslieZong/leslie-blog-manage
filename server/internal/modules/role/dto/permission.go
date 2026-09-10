package dto

// RolePermissionResponse
// 表示“角色拥有的一个权限”。
//
// 注意：
// 这里不是 Casbin Policy 的直接返回结果。
// 我们会把 Casbin 的：
//
// editor | post | read
//
// 转换成：
//
// post:read
//
// 然后再去 permissions 表查询完整权限信息。
type RolePermissionResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
}

// UpdateRolePermissionsRequest
// 用于整体更新某个角色的权限。
type UpdateRolePermissionsRequest struct {

	// Permissions 保存权限名称。
	//
	// 例如：
	//
	// [
	//     "post:read",
	//     "post:create",
	//     "post:update"
	// ]
	//
	// 注意：
	// 这里传的是 Permission.Name，
	// 而不是 Permission.ID。
	Permissions []string `json:"permissions"`
}
