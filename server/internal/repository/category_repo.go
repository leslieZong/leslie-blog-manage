package repository

import (
	"leslie-blog-server/internal/model"

	"gorm.io/gorm"
)

type CategoryRepo struct{ db *gorm.DB }

func NewCategoryRepo(db *gorm.DB) *CategoryRepo { return &CategoryRepo{db: db} }

// List 分类列表（含每个分类文章数）
func (r *CategoryRepo) List(keyword string) ([]model.Category, error) {
	var cats []model.Category
	q := r.db.Order("sort_order ASC, id ASC")
	if keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	if err := q.Find(&cats).Error; err != nil {
		return nil, err
	}
	// 填充文章数
	for i := range cats {
		var n int64
		r.db.Model(&model.Post{}).Where("category_id = ?", cats[i].ID).Count(&n)
		cats[i].PostCount = n
	}
	return cats, nil
}

func (r *CategoryRepo) FindByID(id uint) (*model.Category, error) {
	var c model.Category
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepo) Create(c *model.Category) error {
	return r.db.Create(c).Error
}

func (r *CategoryRepo) Update(c *model.Category) error {
	return r.db.Save(c).Error
}

func (r *CategoryRepo) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}
