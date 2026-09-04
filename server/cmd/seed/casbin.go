package main

import (
	"fmt"

	"leslie-blog-server/internal/pkg/casbin"
)

// seedCasbin 初始化默认权限。
func seedCasbin(
	enforcer *casbin.Enforcer,
	adminUserID string,
) error {

	permissions := [][2]string{

		// User
		{"user", "read"},
		{"user", "create"},
		{"user", "update"},
		{"user", "delete"},

		// Post
		{"post", "read"},
		{"post", "create"},
		{"post", "update"},
		{"post", "delete"},
		{"post", "publish"},

		// Category
		{"category", "read"},
		{"category", "create"},
		{"category", "update"},
		{"category", "delete"},

		// Tag
		{"tag", "read"},
		{"tag", "create"},
		{"tag", "update"},
		{"tag", "delete"},

		// Dashboard
		{"dashboard", "read"},
	}

	// 给 admin 添加权限。
	for _, permission := range permissions {

		if err := enforcer.AddPolicyIfNotExists(
			"admin",
			permission[0],
			permission[1],
		); err != nil {
			return fmt.Errorf(
				"add admin policy failed: %w",
				err,
			)
		}
	}

	// 给 admin 用户分配 admin 角色。
	if err := enforcer.AddRoleForUserIfNotExists(
		adminUserID,
		"admin",
	); err != nil {
		return fmt.Errorf(
			"assign admin role failed: %w",
			err,
		)
	}

	return nil
}
