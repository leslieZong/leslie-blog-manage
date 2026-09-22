package repository

import (
	"context"
	"leslie-blog-server/internal/modules/project/model"
	"strings"

	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(
	db *gorm.DB,
) ProjectRepository {

	return &projectRepository{
		db: db,
	}
}

func (r *projectRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.Project, error) {

	var project model.Project

	err := r.db.
		WithContext(ctx).
		Preload("TechStacks").
		Where("id = ?", id).
		First(&project).
		Error

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*model.Project, error) {

	var project model.Project

	err := r.db.
		WithContext(ctx).
		Preload("TechStacks").
		Where("slug = ?", slug).
		First(&project).
		Error

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) FindByName(
	ctx context.Context,
	name string,
) (*model.Project, error) {

	var project model.Project

	err := r.db.
		WithContext(ctx).
		Preload("TechStacks").
		Where("name = ?", name).
		First(&project).
		Error

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) FindPage(
	ctx context.Context,
	query ProjectListQuery,
) ([]*model.Project, int64, error) {

	var (
		projects []*model.Project
		total    int64
	)

	// 基础查询。
	db := r.db.
		WithContext(ctx).
		Model(&model.Project{}).
		Preload("TechStacks")

	// --------------------------------------------------
	// Keyword
	// --------------------------------------------------

	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {

		like := "%" + keyword + "%"

		db = db.Where(
			"name LIKE ? OR slug LIKE ?",
			like,
			like,
		)
	}

	// --------------------------------------------------
	// Featured
	// --------------------------------------------------

	if query.Featured != nil {

		db = db.Where(
			"featured = ?",
			*query.Featured,
		)
	}

	// --------------------------------------------------
	// Status
	// --------------------------------------------------

	if query.Status != nil {

		db = db.Where(
			"status = ?",
			*query.Status,
		)
	}

	// --------------------------------------------------
	// Count
	// --------------------------------------------------

	if err := db.
		Count(&total).
		Error; err != nil {

		return nil, 0, err
	}

	// --------------------------------------------------
	// List
	// --------------------------------------------------

	if err := db.
		Order("sort ASC").
		Order("created_at DESC").
		Limit(query.PageSize).
		Offset(query.Offset()).
		Find(&projects).
		Error; err != nil {

		return nil, 0, err
	}

	return projects, total, nil
}

func (r *projectRepository) Create(
	ctx context.Context,
	project *model.Project,
) error {

	return r.db.
		WithContext(ctx).
		Create(project).
		Error
}

func (r *projectRepository) Update(
	ctx context.Context,
	project *model.Project,
) error {

	return r.db.
		WithContext(ctx).
		Save(project).
		Error
}

func (r *projectRepository) Delete(
	ctx context.Context,
	id string,
) error {

	return r.db.
		WithContext(ctx).
		Delete(
			&model.Project{},
			"id = ?",
			id,
		).
		Error
}

func (r *projectRepository) FindPublicByID(
	ctx context.Context,
	id string,
) (*model.Project, error) {

	var project model.Project

	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Where("status = ?", 1).
		First(&project).
		Error

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) FindPublicBySlug(
	ctx context.Context,
	slug string,
) (*model.Project, error) {

	var project model.Project

	err := r.db.
		WithContext(ctx).
		Where("slug = ?", slug).
		Where("status = ?", 1).
		First(&project).
		Error

	if err != nil {
		return nil, err
	}

	return &project, nil
}
