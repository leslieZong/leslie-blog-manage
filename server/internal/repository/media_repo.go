package repository

import (
	"leslie-blog-server/internal/model"

	"gorm.io/gorm"
)

type MediaRepo struct{ db *gorm.DB }

func NewMediaRepo(db *gorm.DB) *MediaRepo { return &MediaRepo{db: db} }

func (r *MediaRepo) List(keyword, mtype string, page, pageSize int) ([]model.Media, int64, error) {
	var total int64
	q := r.db.Model(&model.Media{})
	if keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	if mtype != "" {
		q = q.Where("type = ?", mtype)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Media
	if err := paginate(q.Order("id DESC"), page, pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *MediaRepo) Create(m *model.Media) error { return r.db.Create(m).Error }
func (r *MediaRepo) Delete(id uint) error          { return r.db.Delete(&model.Media{}, id).Error }
