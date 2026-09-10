package repository

import (
	"context"

	"leslie-blog-server/internal/modules/permission/model"
)

// PermissionRepository 定义 Permission 的数据访问能力。
//
// 注意：
//
// 这里是“接口”，不是具体实现。
//
// Service 只关心：
// “我需要查询 Permission”
//
// 而不需要关心：
// “到底是 GORM、MySQL 还是其他数据库”。
type PermissionRepository interface {
	// FindByID 根据 Permission ID 查询权限。
	FindByID(ctx context.Context, id string) (*model.Permission, error)

	// FindByName 根据权限名称查询权限。
	//
	// 例如：
	//
	// post:read
	FindByName(ctx context.Context, name string) (*model.Permission, error)

	// FindAll 查询所有权限。
	FindAll(ctx context.Context) ([]*model.Permission, error)

	// FindByResource 查询某个资源下的所有权限。
	//
	// 例如：
	//
	// resource = post
	//
	// 返回：
	//
	// post:read
	// post:create
	// post:update
	// post:delete
	FindByResource(
		ctx context.Context,
		resource string,
	) ([]*model.Permission, error)

	// FindByNames 根据权限名称列表查询权限。
	//
	// 例如：
	//
	// names = []string{"post:read", "post:create"}
	//
	// 返回：
	//
	// post:read
	// post:create
	FindByNames(ctx context.Context, names []string) ([]*model.Permission, error)
}
