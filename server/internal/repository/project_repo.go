package repository

import (
	"leslie-blog-server/internal/model"

	"gorm.io/gorm"
)

type ProjectRepo struct{ db *gorm.DB }

func NewProjectRepo(db *gorm.DB) *ProjectRepo { return &ProjectRepo{db: db} }

func (r *ProjectRepo) List(keyword string, status *int8, page, pageSize int) ([]model.Project, int64, error) {
	var total int64
	q := r.db.Model(&model.Project{})
	if keyword != "" {
		q = q.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Project
	if err := paginate(q.Order("sort_order ASC, id DESC"), page, pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// ListAll 不分页列表（公开侧）
func (r *ProjectRepo) ListAll() ([]model.Project, error) {
	var list []model.Project
	if err := r.db.Order("sort_order ASC, id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ProjectRepo) FindByID(id uint) (*model.Project, error) {
	var p model.Project
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepo) Create(p *model.Project) error { return r.db.Create(p).Error }
func (r *ProjectRepo) Update(p *model.Project) error { return r.db.Save(p).Error }
func (r *ProjectRepo) Delete(id uint) error          { return r.db.Delete(&model.Project{}, id).Error }
