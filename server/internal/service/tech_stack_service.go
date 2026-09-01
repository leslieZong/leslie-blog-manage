package service

import (
	"errors"

	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/repository"
	"gorm.io/gorm"
)

type TechStackService struct{ repo *repository.TechStackRepo }

func NewTechStackService(repo *repository.TechStackRepo) *TechStackService {
	return &TechStackService{repo: repo}
}

func (s *TechStackService) List(keyword string) ([]model.TechStack, error) {
	return s.repo.List(keyword)
}

func (s *TechStackService) Get(id uint) (*model.TechStack, error) { return s.repo.FindByID(id) }

func (s *TechStackService) Create(req model.TechStackReq) error {
	t := model.TechStack{
		Name:        req.Name,
		Icon:        req.Icon,
		Category:    req.Category,
		Level:       req.Level,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
	return s.repo.Create(&t)
}

func (s *TechStackService) Update(id uint, req model.TechStackReq) error {
	t, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("技术不存在")
		}
		return err
	}
	t.Name = req.Name
	t.Icon = req.Icon
	t.Category = req.Category
	t.Level = req.Level
	t.Description = req.Description
	t.SortOrder = req.SortOrder
	return s.repo.Update(t)
}

func (s *TechStackService) Delete(id uint) error { return s.repo.Delete(id) }
