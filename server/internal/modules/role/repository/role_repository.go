package repository

import (
	"context"

	"leslie-blog-server/internal/modules/role/model"
)

// RoleListParams 表示角色列表查询条件。
type RoleListParams struct {
	// Keyword 用于搜索角色名称和显示名称。
	Keyword string

	// Status 是可选的角色状态。
	Status *int8

	// Offset 是分页偏移量。
	Offset int

	// Limit 是分页数量。
	Limit int
}

// RoleRepository 定义角色的数据访问能力。
//
// Service 只依赖这个接口，
// 而不直接依赖 MySQL / GORM。
type RoleRepository interface {

	// 根据 ID 查询角色。
	FindByID(
		ctx context.Context,
		id string,
	) (*model.Role, error)

	// 根据角色名称查询。
	FindByName(
		ctx context.Context,
		name string,
	) (*model.Role, error)

	// 查询角色列表。
	List(
		ctx context.Context,
		params RoleListParams,
	) ([]*model.Role, int64, error)

	// 创建角色。
	Create(
		ctx context.Context,
		role *model.Role,
	) error

	// 更新角色指定字段。
	UpdateFields(
		ctx context.Context,
		id string,
		updates map[string]any,
	) error

	// 删除角色。
	Delete(
		ctx context.Context,
		role *model.Role,
	) error
}
