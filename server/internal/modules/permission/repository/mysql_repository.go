package repository

import (
	"context"
	"leslie-blog-server/internal/modules/permission/model"

	"gorm.io/gorm"
)

// gormPermissionRepository 是 PermissionRepository 的 GORM 实现。
//
// 注意：
//
// 这个结构体没有导出，说明：
// Repository 的具体实现只应该通过构造函数暴露。
type gormPermissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository 创建 PermissionRepository。
//
// 参数：
//
// db
// ↓
// GORM 数据库连接
//
// 返回：
//
// # PermissionRepository
//
// Service 不需要知道具体实现是 gormPermissionRepository。
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &gormPermissionRepository{
		db: db,
	}
}

// FindByID 根据 ID 查询 Permission。
func (r *gormPermissionRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.Permission, error) {

	var permission model.Permission

	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&permission).
		Error

	if err != nil {
		return nil, err
	}

	return &permission, nil
}

// FindByName 根据权限名称查询 Permission。
//
// 例如：
//
// post:read
func (r *gormPermissionRepository) FindByName(
	ctx context.Context,
	name string,
) (*model.Permission, error) {

	var permission model.Permission

	err := r.db.
		WithContext(ctx).
		Where("name = ?", name).
		First(&permission).
		Error

	if err != nil {
		return nil, err
	}

	return &permission, nil
}

// FindAll 查询所有未被软删除的 Permission。
func (r *gormPermissionRepository) FindAll(
	ctx context.Context,
) ([]*model.Permission, error) {

	var permissions []*model.Permission

	err := r.db.
		WithContext(ctx).
		Order("resource ASC, action ASC").
		Find(&permissions).
		Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// FindByResource 查询指定资源下的所有 Permission。
//
// 例如：
//
// FindByResource(ctx, "post")
//
// 返回：
//
// post:read
// post:create
// post:update
// post:delete
// post:publish
func (r *gormPermissionRepository) FindByResource(
	ctx context.Context,
	resource string,
) ([]*model.Permission, error) {

	var permissions []*model.Permission

	err := r.db.
		WithContext(ctx).
		Where("resource = ?", resource).
		Order("action ASC").
		Find(&permissions).
		Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *gormPermissionRepository) FindByNames(
	ctx context.Context,
	names []string,
) ([]*model.Permission, error) {

	var permissions []*model.Permission

	err := r.db.
		WithContext(ctx).
		Where("name IN ?", names).
		Order("resource ASC, action ASC").
		Find(&permissions).
		Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// Create 创建权限。
func (r *gormPermissionRepository) Create(
	ctx context.Context,
	permission *model.Permission,
) error {
	return r.db.WithContext(ctx).Create(permission).Error
}
