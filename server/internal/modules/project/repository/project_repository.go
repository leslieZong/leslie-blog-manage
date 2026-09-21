package repository

import (
	"context"

	"leslie-blog-server/internal/modules/project/model"
)

type ProjectRepository interface {

	// 根据 ID 查询。
	FindByID(
		ctx context.Context,
		id string,
	) (*model.Project, error)

	// 根据 slug 查询。
	FindBySlug(
		ctx context.Context,
		slug string,
	) (*model.Project, error)

	// 根据名称查询。
	FindByName(
		ctx context.Context,
		name string,
	) (*model.Project, error)

	// 分页查询。
	FindPage(
		ctx context.Context,
		query ProjectListQuery,
	) ([]*model.Project, int64, error)

	// 创建。
	Create(
		ctx context.Context,
		project *model.Project,
	) error

	// 修改。
	Update(
		ctx context.Context,
		project *model.Project,
	) error

	// 删除。
	Delete(
		ctx context.Context,
		id string,
	) error
	FindPublicByID(
		ctx context.Context,
		id string,
	) (*model.Project, error)

	FindPublicBySlug(
		ctx context.Context,
		slug string,
	) (*model.Project, error)
}
