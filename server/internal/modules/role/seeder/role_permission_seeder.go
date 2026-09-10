package seeder

import (
	"context"
	"errors"
	permissionRepository "leslie-blog-server/internal/modules/permission/repository"
	"leslie-blog-server/internal/modules/role/repository"
	"leslie-blog-server/internal/pkg/casbin"
)

// defaultRolePermissions
//
// key:
//
//	角色名称
//
// value:
//
//	这个角色默认拥有的权限名称。
var defaultRolePermissions = map[string][]string{

	// admin 拥有全部权限。
	"admin": {
		"user:read",
		"user:create",
		"user:update",
		"user:delete",
		"user:assign-role",

		"role:read",
		"role:create",
		"role:update",
		"role:delete",

		"post:read",
		"post:create",
		"post:update",
		"post:delete",
		"post:publish",

		"category:read",
		"category:create",
		"category:update",
		"category:delete",

		"tag:read",
		"tag:create",
		"tag:update",
		"tag:delete",

		"project:read",
		"project:create",
		"project:update",
		"project:delete",
	},

	// editor 负责内容管理。
	"editor": {
		"post:read",
		"post:create",
		"post:update",
		"post:publish",

		"category:read",

		"tag:read",

		"project:read",
	},

	// viewer 只拥有查看权限。
	"viewer": {
		"post:read",
		"category:read",
		"tag:read",
		"project:read",
	},
}

func SeedRolePermissions(
	ctx context.Context,
	roleRepo repository.RoleRepository,
	permissionRepo permissionRepository.PermissionRepository,
	enforcer *casbin.Enforcer,
) error {

	for roleName, permissionNames := range defaultRolePermissions {

		// 1. 查询 Role。
		role, err := roleRepo.FindByName(
			ctx,
			roleName,
		)

		if err != nil {
			return err
		}

		// 2. 批量查询 Permission。
		permissions, err := permissionRepo.FindByNames(
			ctx,
			permissionNames,
		)

		if err != nil {
			return err
		}

		// 3. 必须保证配置中的每一个权限
		//    都已经存在于 permissions 表。
		if len(permissions) != len(permissionNames) {
			return errors.New("permission not found")
		}

		// 4. 将 Permission 转换成 Casbin Policy。
		for _, permission := range permissions {

			err := enforcer.AddPolicyIfNotExists(
				role.Name,
				permission.Resource,
				permission.Action,
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
