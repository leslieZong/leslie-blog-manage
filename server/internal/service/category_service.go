package service

import (
	"errors"
	"strings"

	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/repository"
	"gorm.io/gorm"
)

type CategoryService struct{ repo *repository.CategoryRepo }

func NewCategoryService(repo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List(keyword string) ([]model.Category, error) {
	return s.repo.List(keyword)
}

func (s *CategoryService) Get(id uint) (*model.Category, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) Create(req model.CategoryReq) error {
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = slugify(req.Name)
	}
	c := model.Category{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
	}
	return s.repo.Create(&c)
}

func (s *CategoryService) Update(id uint, req model.CategoryReq) error {
	c, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分类不存在")
		}
		return err
	}
	if s2 := strings.TrimSpace(req.Slug); s2 != "" {
		c.Slug = s2
	}
	c.Name = req.Name
	c.Description = req.Description
	c.ParentID = req.ParentID
	c.SortOrder = req.SortOrder
	return s.repo.Update(c)
}

func (s *CategoryService) Delete(id uint) error { return s.repo.Delete(id) }
