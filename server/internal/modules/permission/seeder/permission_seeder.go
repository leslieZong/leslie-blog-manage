package seeder

import (
	"context"
	"errors"

	"leslie-blog-server/internal/modules/permission/model"
	permissionRepository "leslie-blog-server/internal/modules/permission/repository"
	"leslie-blog-server/internal/pkg/ulid"

	"gorm.io/gorm"
)

// SeedPermission 表示一个系统初始化权限。
type SeedPermission struct {
	Name        string
	DisplayName string
	Resource    string
	Action      string
	Description string
}

// systemPermissions 是 Leslie Blog 系统的权限定义。
//
// 注意：
//
// 这里定义的是：
// “系统有哪些合法权限”
//
// 例如：
//
// post:create
//
// 并不是：
// “谁拥有 post:create”
//
// “谁拥有权限”由 Casbin Policy 管理。
var systemPermissions = []SeedPermission{
	{
		Name:        "user:read",
		DisplayName: "查看用户",
		Resource:    "user",
		Action:      "read",
		Description: "查看后台用户",
	},
	{
		Name:        "user:create",
		DisplayName: "创建用户",
		Resource:    "user",
		Action:      "create",
		Description: "创建后台用户",
	},
	{
		Name:        "user:update",
		DisplayName: "修改用户",
		Resource:    "user",
		Action:      "update",
		Description: "修改后台用户",
	},
	{
		Name:        "user:delete",
		DisplayName: "删除用户",
		Resource:    "user",
		Action:      "delete",
		Description: "删除后台用户",
	},
	{
		Name:        "user:assign-role",
		DisplayName: "分配角色",
		Resource:    "user",
		Action:      "assign-role",
		Description: "为用户分配角色",
	},

	{
		Name:        "role:read",
		DisplayName: "查看角色",
		Resource:    "role",
		Action:      "read",
		Description: "查看系统角色",
	},
	{
		Name:        "role:create",
		DisplayName: "创建角色",
		Resource:    "role",
		Action:      "create",
		Description: "创建系统角色",
	},
	{
		Name:        "role:update",
		DisplayName: "修改角色",
		Resource:    "role",
		Action:      "update",
		Description: "修改系统角色",
	},
	{
		Name:        "role:delete",
		DisplayName: "删除角色",
		Resource:    "role",
		Action:      "delete",
		Description: "删除系统角色",
	},

	{
		Name:        "post:read",
		DisplayName: "查看文章",
		Resource:    "post",
		Action:      "read",
		Description: "查看文章",
	},
	{
		Name:        "post:create",
		DisplayName: "创建文章",
		Resource:    "post",
		Action:      "create",
		Description: "创建文章",
	},
	{
		Name:        "post:update",
		DisplayName: "修改文章",
		Resource:    "post",
		Action:      "update",
		Description: "修改文章",
	},
	{
		Name:        "post:delete",
		DisplayName: "删除文章",
		Resource:    "post",
		Action:      "delete",
		Description: "删除文章",
	},
	{
		Name:        "post:publish",
		DisplayName: "发布文章",
		Resource:    "post",
		Action:      "publish",
		Description: "发布文章",
	},

	{
		Name:        "category:read",
		DisplayName: "查看分类",
		Resource:    "category",
		Action:      "read",
		Description: "查看文章分类",
	},
	{
		Name:        "category:create",
		DisplayName: "创建分类",
		Resource:    "category",
		Action:      "create",
		Description: "创建文章分类",
	},
	{
		Name:        "category:update",
		DisplayName: "修改分类",
		Resource:    "category",
		Action:      "update",
		Description: "修改文章分类",
	},
	{
		Name:        "category:delete",
		DisplayName: "删除分类",
		Resource:    "category",
		Action:      "delete",
		Description: "删除文章分类",
	},

	{
		Name:        "tag:read",
		DisplayName: "查看标签",
		Resource:    "tag",
		Action:      "read",
		Description: "查看文章标签",
	},
	{
		Name:        "tag:create",
		DisplayName: "创建标签",
		Resource:    "tag",
		Action:      "create",
		Description: "创建文章标签",
	},
	{
		Name:        "tag:update",
		DisplayName: "修改标签",
		Resource:    "tag",
		Action:      "update",
		Description: "修改文章标签",
	},
	{
		Name:        "tag:delete",
		DisplayName: "删除标签",
		Resource:    "tag",
		Action:      "delete",
		Description: "删除文章标签",
	},

	{
		Name:        "project:read",
		DisplayName: "查看项目",
		Resource:    "project",
		Action:      "read",
		Description: "查看项目",
	},
	{
		Name:        "project:create",
		DisplayName: "创建项目",
		Resource:    "project",
		Action:      "create",
		Description: "创建项目",
	},
	{
		Name:        "project:update",
		DisplayName: "修改项目",
		Resource:    "project",
		Action:      "update",
		Description: "修改项目",
	},
	{
		Name:        "project:delete",
		DisplayName: "删除项目",
		Resource:    "project",
		Action:      "delete",
		Description: "删除项目",
	},
}

// Seed 初始化系统 Permission。
//
// 这个方法必须保证幂等性：
//
// 第一次启动：
// 创建所有不存在的权限。
//
// 第二次启动：
// 已存在的权限直接跳过。
func Seed(
	ctx context.Context,
	repo permissionRepository.PermissionRepository,
) error {

	// 遍历所有系统预定义权限。
	for _, item := range systemPermissions {

		// 第一步：
		// 根据 Name 查询权限是否已经存在。
		permission, err := repo.FindByName(
			ctx,
			item.Name,
		)

		// 查询成功，说明数据库中已经存在。
		if err == nil && permission != nil {

			// 已存在，直接跳过。
			continue
		}

		// 如果错误不是“记录不存在”，
		// 说明数据库查询本身发生了异常。
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 第二步：
		// 数据库不存在该 Permission，
		// 创建新的 Permission。
		permission = &model.Permission{
			ID:          ulid.New(),
			Name:        item.Name,
			DisplayName: item.DisplayName,
			Resource:    item.Resource,
			Action:      item.Action,
			Description: &item.Description,
		}

		// 第三步：
		// 写入数据库。
		if err := repo.Create(ctx, permission); err != nil {
			return err
		}
	}

	return nil
}
