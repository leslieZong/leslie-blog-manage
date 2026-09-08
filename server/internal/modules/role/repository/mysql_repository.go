package repository

import (
	"context"

	"leslie-blog-server/internal/modules/role/model"

	"gorm.io/gorm"
)

// mysqlRoleRepository 是 RoleRepository 的 MySQL 实现。
type mysqlRoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色 Repository。
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &mysqlRoleRepository{
		db: db,
	}
}

// FindByID 根据角色 ID 查询角色。
func (r *mysqlRoleRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.Role, error) {

	var role model.Role

	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&role).
		Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// FindByName 根据角色名称查询。
func (r *mysqlRoleRepository) FindByName(
	ctx context.Context,
	name string,
) (*model.Role, error) {

	var role model.Role

	err := r.db.
		WithContext(ctx).
		Where("name = ?", name).
		First(&role).
		Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// List 查询角色列表。
func (r *mysqlRoleRepository) List(
	ctx context.Context,
	params RoleListParams,
) ([]*model.Role, int64, error) {

	query := r.db.
		WithContext(ctx).
		Model(&model.Role{})

	// 如果有关键词，
	// 搜索 name 和 display_name。
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"

		query = query.Where(
			"name LIKE ? OR display_name LIKE ?",
			keyword,
			keyword,
		)
	}

	// 如果传递了状态，
	// 则增加状态过滤。
	if params.Status != nil {
		query = query.Where(
			"status = ?",
			*params.Status,
		)
	}

	// 查询总数量。
	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询当前页数据。
	var roles []*model.Role

	err := query.
		Order("created_at DESC").
		Offset(params.Offset).
		Limit(params.Limit).
		Find(&roles).
		Error

	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// Create 创建角色。
func (r *mysqlRoleRepository) Create(
	ctx context.Context,
	role *model.Role,
) error {

	return r.db.
		WithContext(ctx).
		Create(role).
		Error
}

// UpdateFields 更新角色指定字段。
func (r *mysqlRoleRepository) UpdateFields(
	ctx context.Context,
	id string,
	updates map[string]any,
) error {

	return r.db.
		WithContext(ctx).
		Model(&model.Role{}).
		Where("id = ?", id).
		Updates(updates).
		Error
}

// Delete 执行角色软删除。
func (r *mysqlRoleRepository) Delete(
	ctx context.Context,
	role *model.Role,
) error {

	return r.db.
		WithContext(ctx).
		Delete(role).
		Error
}
