package dto

import (
	"time"

	"leslie-blog-server/internal/modules/role/model"
)

// RoleResponse 是返回给前端的角色数据。
//
// 注意：
// Model 是数据库结构。
// DTO 是 API 结构。
//
// 两者不要混用。
type RoleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	Description string    `json:"description"`
	Status      int8      `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// FromRole 将数据库 Model 转换成 API DTO。
func FromRole(role *model.Role) *RoleResponse {
	if role == nil {
		return nil
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		Status:      role.Status,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// CreateRoleRequest 是创建角色时的请求参数。
type CreateRoleRequest struct {
	// Name 是程序内部使用的角色名称。
	//
	// 例如：
	// admin
	// editor
	// viewer
	Name string `json:"name" validate:"required,min=2,max=50"`

	// DisplayName 是管理员看到的名称。
	DisplayName string `json:"displayName" validate:"required,max=100"`

	// Description 是角色说明。
	Description string `json:"description" validate:"omitempty,max=255"`
}

// UpdateRoleRequest 是修改角色时的请求参数。
type UpdateRoleRequest struct {
	DisplayName string `json:"displayName" validate:"required,max=100"`

	Description string `json:"description" validate:"omitempty,max=255"`

	Status int8 `json:"status" validate:"oneof=0 1"`
}
