package repository

import (
	"leslie-blog-server/internal/model"

	"gorm.io/gorm"
)

type TechStackRepo struct{ db *gorm.DB }

func NewTechStackRepo(db *gorm.DB) *TechStackRepo { return &TechStackRepo{db: db} }

func (r *TechStackRepo) List(keyword string) ([]model.TechStack, error) {
	var list []model.TechStack
	q := r.db.Order("sort_order ASC, id ASC")
	if keyword != "" {
		q = q.Where("name LIKE ? OR category LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *TechStackRepo) FindByID(id uint) (*model.TechStack, error) {
	var t model.TechStack
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TechStackRepo) Create(t *model.TechStack) error { return r.db.Create(t).Error }
func (r *TechStackRepo) Update(t *model.TechStack) error { return r.db.Save(t).Error }
func (r *TechStackRepo) Delete(id uint) error            { return r.db.Delete(&model.TechStack{}, id).Error }
