package seeder

import (
	"context"
	"errors"
	"leslie-blog-server/internal/modules/role/model"
	"leslie-blog-server/internal/modules/role/repository"
	"leslie-blog-server/internal/pkg/ulid"
	"strings"

	"gorm.io/gorm"
)

// SeedRole 表示一个系统初始化角色。
//
// Seeder 不需要直接关心 HTTP、Handler。
// 它只需要描述：
// “系统应该有哪些默认角色”。
type SeedRole struct {
	// Name 是角色的程序内部名称。
	//
	// 例如：
	// admin
	// editor
	// viewer
	Name string

	// DisplayName 是给管理员看的名称。
	//
	// 例如：
	// 系统管理员
	// 内容编辑
	// 访客
	DisplayName string

	// Description 是角色说明。
	Description string
}

var defaultRoles = []SeedRole{
	{
		Name:        "admin",
		DisplayName: "系统管理员",
		Description: "拥有系统全部管理权限",
	},
	{
		Name:        "editor",
		DisplayName: "内容编辑",
		Description: "负责文章及内容相关管理",
	},
	{
		Name:        "viewer",
		DisplayName: "只读用户",
		Description: "只能查看系统内容",
	},
}

// Seed
// 初始化系统默认角色。
func Seed(
	ctx context.Context,
	repo repository.RoleRepository,
) error {

	for _, item := range defaultRoles {

		// 先查询角色是否已经存在。
		role, err := repo.FindByName(
			ctx,
			item.Name,
		)

		if err == nil && role != nil {
			// 已经存在。
			// Seeder 应该直接跳过，而不是重复创建。
			continue
		}

		// 这里的 NotFound 判断，
		// 必须替换成你项目当前已有的错误判断方式。
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 创建新的 Role Model。
		role = &model.Role{
			ID:          ulid.New(),
			Name:        item.Name,
			DisplayName: item.DisplayName,
			Description: strings.TrimSpace(item.Description),
			Status:      1,
		}

		if err := repo.Create(ctx, role); err != nil {
			return err
		}
	}

	return nil
}
