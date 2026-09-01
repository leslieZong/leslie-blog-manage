package service

import (
	"errors"

	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/repository"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProjectService struct{ repo *repository.ProjectRepo }

func NewProjectService(repo *repository.ProjectRepo) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) List(keyword string, status *int8, page, pageSize int) ([]model.Project, int64, error) {
	return s.repo.List(keyword, status, page, pageSize)
}

func (s *ProjectService) ListAll() ([]model.Project, error) { return s.repo.ListAll() }

func (s *ProjectService) Get(id uint) (*model.Project, error) { return s.repo.FindByID(id) }

func (s *ProjectService) Create(req model.ProjectReq) error {
	p := model.Project{
		Name:        req.Name,
		Description: req.Description,
		Cover:       req.Cover,
		DemoURL:     req.DemoURL,
		RepoURL:     req.RepoURL,
		TechStack:   toJSON(req.TechStack),
		Status:      req.Status,
		SortOrder:   req.SortOrder,
	}
	return s.repo.Create(&p)
}

func (s *ProjectService) Update(id uint, req model.ProjectReq) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
		}
		return err
	}
	p.Name = req.Name
	p.Description = req.Description
	p.Cover = req.Cover
	p.DemoURL = req.DemoURL
	p.RepoURL = req.RepoURL
	p.TechStack = toJSON(req.TechStack)
	p.Status = req.Status
	p.SortOrder = req.SortOrder
	return s.repo.Update(p)
}

func (s *ProjectService) Delete(id uint) error { return s.repo.Delete(id) }

// toJSON 把 []string 转为 datatypes.JSON；空则 nil
func toJSON(ss []string) datatypes.JSON {
	if len(ss) == 0 {
		return nil
	}
	b, _ := datatypes.JSON.Marshal(ss)
	return b
}
