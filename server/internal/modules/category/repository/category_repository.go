package repository

import (
	"context"

	"leslie-blog-server/internal/modules/category/model"
)

// CategoryRepository
//
// Repository 负责 Category 数据访问。
//
// Service 不应该直接操作 GORM。
// 所以 Service 通过 Repository 查询数据库。
type CategoryRepository interface {

	// 根据 ID 查询分类。
	FindByID(
		ctx context.Context,
		id string,
	) (*model.Category, error)

	// 根据名称查询分类。
	FindByName(
		ctx context.Context,
		name string,
	) (*model.Category, error)

	// 根据 slug 查询分类。
	FindBySlug(
		ctx context.Context,
		slug string,
	) (*model.Category, error)

	// 查询全部分类。
	FindAll(
		ctx context.Context,
	) ([]*model.Category, error)

	// 创建分类。
	Create(
		ctx context.Context,
		category *model.Category,
	) error

	// 更新分类。
	Update(
		ctx context.Context,
		category *model.Category,
	) error

	// 删除分类。
	Delete(
		ctx context.Context,
		id string,
	) error

	// 查询分类下面是否存在文章。
	//
	// 删除分类之前需要使用。
	CountPosts(
		ctx context.Context,
		categoryID string,
	) (int64, error)
}
