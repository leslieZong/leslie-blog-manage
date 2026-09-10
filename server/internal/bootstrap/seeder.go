package bootstrap

import (
	"context"
	permissionRepository "leslie-blog-server/internal/modules/permission/repository"
	permissionSeeder "leslie-blog-server/internal/modules/permission/seeder"
	roleRepository "leslie-blog-server/internal/modules/role/repository"
	roleSeeder "leslie-blog-server/internal/modules/role/seeder"
	"leslie-blog-server/internal/pkg/casbin"

	"gorm.io/gorm"
)

func SeedDatabase(
	ctx context.Context,
	db *gorm.DB,
	enforcer *casbin.Enforcer,
) error {

	// 创建 Repository。
	roleRepo := roleRepository.NewRoleRepository(db)

	permissionRepo :=
		permissionRepository.NewPermissionRepository(db)

	// 第一步：初始化 Permission。
	if err := permissionSeeder.Seed(
		ctx,
		permissionRepo,
	); err != nil {
		return err
	}

	// 第二步：初始化 Role。
	if err := roleSeeder.Seed(
		ctx,
		roleRepo,
	); err != nil {
		return err
	}

	// 第三步：
	// 初始化 Role → Permission。
	if err := roleSeeder.SeedRolePermissions(
		ctx,
		roleRepo,
		permissionRepo,
		enforcer,
	); err != nil {
		return err
	}

	return nil
}
