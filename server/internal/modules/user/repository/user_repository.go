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
	FindByID(
		ctx context.Context,
		id string,
	) (*model.User, error)

	FindByUsername(
		ctx context.Context,
		username string,
	) (*model.User, error)

	List(
		ctx context.Context,
		params UserListParams,
	) ([]*model.User, int64, error)

	Create(
		ctx context.Context,
		user *model.User,
	) error

	Update(
		ctx context.Context,
		user *model.User,
	) error

	UpdateFields(
		ctx context.Context,
		id string,
		updates map[string]any,
	) error

	Delete(
		ctx context.Context,
		user *model.User,
	) error
}
