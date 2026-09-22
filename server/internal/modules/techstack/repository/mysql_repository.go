package repository

import (
	"context"
	"leslie-blog-server/internal/modules/techstack/model"

	"gorm.io/gorm"
)

type techStackRepository struct {
	db *gorm.DB
}

// NewTechStackRepository 创建 TechStack Repository。
func NewTechStackRepository(
	db *gorm.DB,
) TechStackRepository {

	return &techStackRepository{
		db: db,
	}
}

func (r *techStackRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.TechStack, error) {

	var techStack model.TechStack

	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&techStack).
		Error

	if err != nil {
		return nil, err
	}

	return &techStack, nil
}

func (r *techStackRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*model.TechStack, error) {

	var techStack model.TechStack

	err := r.db.
		WithContext(ctx).
		Where("slug = ?", slug).
		First(&techStack).
		Error

	if err != nil {
		return nil, err
	}

	return &techStack, nil
}

func (r *techStackRepository) FindByName(
	ctx context.Context,
	name string,
) (*model.TechStack, error) {

	var techStack model.TechStack

	err := r.db.
		WithContext(ctx).
		Where("name = ?", name).
		First(&techStack).
		Error

	if err != nil {
		return nil, err
	}

	return &techStack, nil
}

func (r *techStackRepository) FindAll(
	ctx context.Context,
) ([]*model.TechStack, error) {

	var techStacks []*model.TechStack

	err := r.db.
		WithContext(ctx).
		Order("sort ASC").
		Order("created_at DESC").
		Find(&techStacks).
		Error

	if err != nil {
		return nil, err
	}

	return techStacks, nil
}

func (r *techStackRepository) FindPublic(
	ctx context.Context,
) ([]*model.TechStack, error) {

	var techStacks []*model.TechStack

	err := r.db.
		WithContext(ctx).
		Where("status = ?", 1).
		Order("sort ASC").
		Order("created_at DESC").
		Find(&techStacks).
		Error

	if err != nil {
		return nil, err
	}

	return techStacks, nil
}

func (r *techStackRepository) Create(
	ctx context.Context,
	techStack *model.TechStack,
) error {

	return r.db.
		WithContext(ctx).
		Create(techStack).
		Error
}

func (r *techStackRepository) Update(
	ctx context.Context,
	techStack *model.TechStack,
) error {

	return r.db.
		WithContext(ctx).
		Save(techStack).
		Error
}

func (r *techStackRepository) Delete(
	ctx context.Context,
	id string,
) error {

	return r.db.
		WithContext(ctx).
		Delete(
			&model.TechStack{},
			"id = ?",
			id,
		).
		Error
}

func (r *techStackRepository) CountProjects(
	ctx context.Context,
	techStackID string,
) (int64, error) {

	var count int64

	err := r.db.
		WithContext(ctx).
		Table("project_tech_stacks").
		Where(
			"tech_stack_id = ?",
			techStackID,
		).
		Count(&count).
		Error

	return count, err
}

func (r *techStackRepository) FindByIDs(
	ctx context.Context,
	ids []string,
) ([]*model.TechStack, error) {

	if len(ids) == 0 {
		return []*model.TechStack{}, nil
	}

	var techStacks []*model.TechStack

	err := r.db.
		WithContext(ctx).
		Where("id IN ?", ids).
		Where("status = ?", 1).
		Find(&techStacks).
		Error

	if err != nil {
		return nil, err
	}

	return techStacks, nil
}
