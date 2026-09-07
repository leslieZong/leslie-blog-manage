package repository

import (
	"context"

	"leslie-blog-server/internal/modules/user/model"
)

// UserListParams 是用户列表查询条件。
type UserListParams struct {
	Keyword string
	Status  *int8
	Offset  int
	Limit   int
}

// UserRepository 定义用户数据访问能力。
type UserRepository interface {

	// 根据 ID 查询用户。
	FindByID(
		ctx context.Context,
		id string,
	) (*model.User, error)

	// 根据用户名查询用户。
	FindByUsername(
		ctx context.Context,
		username string,
	) (*model.User, error)

	// 查询用户列表。
	List(
		ctx context.Context,
		params UserListParams,
	) ([]*model.User, int64, error)

	// 创建用户。
	Create(
		ctx context.Context,
		user *model.User,
	) error

	// 更新用户。
	Update(
		ctx context.Context,
		user *model.User,
	) error

	// 删除用户。
	Delete(
		ctx context.Context,
		user *model.User,
	) error
}
