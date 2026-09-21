package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/project/dto"
	"leslie-blog-server/internal/modules/project/model"
	"leslie-blog-server/internal/modules/project/repository"
	"leslie-blog-server/internal/pkg/ulid"

	"gorm.io/gorm"
)

type ProjectService interface {
	Create(
		ctx context.Context,
		req dto.CreateProjectRequest,
	) (*model.Project, error)

	GetByID(
		ctx context.Context,
		id string,
	) (*model.Project, error)

	GetBySlug(
		ctx context.Context,
		slug string,
	) (*model.Project, error)

	ListPage(
		ctx context.Context,
		query repository.ProjectListQuery,
	) ([]*model.Project, int64, error)

	Update(
		ctx context.Context,
		id string,
		req dto.UpdateProjectRequest,
	) (*model.Project, error)

	Delete(
		ctx context.Context,
		id string,
	) error
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(
	repo repository.ProjectRepository,
) ProjectService {

	return &projectService{
		repo: repo,
	}
}

func (s *projectService) Create(
	ctx context.Context,
	req dto.CreateProjectRequest,
) (*model.Project, error) {

	name := strings.TrimSpace(req.Name)
	slug := strings.TrimSpace(req.Slug)

	// ---------------------------------------------
	// 1. 基础参数
	// ---------------------------------------------

	if name == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"project name is required",
		)
	}

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"project slug is required",
		)
	}

	// ---------------------------------------------
	// 2. Name 唯一性
	// ---------------------------------------------

	if err := s.validateNameUnique(
		ctx,
		name,
		"",
	); err != nil {
		return nil, err
	}

	// ---------------------------------------------
	// 3. Slug 唯一性
	// ---------------------------------------------

	if err := s.validateSlugUnique(
		ctx,
		slug,
		"",
	); err != nil {
		return nil, err
	}

	// ---------------------------------------------
	// 4. 创建模型
	// ---------------------------------------------

	project := &model.Project{
		ID:          ulid.New(),
		Name:        name,
		Slug:        slug,
		Description: req.Description,
		Cover:       req.Cover,
		GitHubURL:   req.GitHubURL,
		DemoURL:     req.DemoURL,
		Featured:    req.Featured,

		// 新 Project 默认启用。
		Status: 1,

		Sort: req.Sort,
	}

	// ---------------------------------------------
	// 5. 持久化
	// ---------------------------------------------

	if err := s.repo.Create(
		ctx,
		project,
	); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) validateNameUnique(
	ctx context.Context,
	name string,
	currentID string,
) error {

	project, err := s.repo.FindByName(
		ctx,
		name,
	)

	// 没找到。
	//
	// 说明当前 name 没有冲突。
	if errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		return nil
	}

	// 其他数据库错误。
	if err != nil {
		return err
	}

	// 找到了，但是是当前 Project 自己。
	if project.ID == currentID {
		return nil
	}

	// 找到了其他 Project。
	return appErrors.New(
		appErrors.ErrInvalidParams,
		http.StatusBadRequest,
		"project name already exists",
	)
}

func (s *projectService) validateSlugUnique(
	ctx context.Context,
	slug string,
	currentID string,
) error {

	project, err := s.repo.FindBySlug(
		ctx,
		slug,
	)

	if errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		return nil
	}

	if err != nil {
		return err
	}

	if project.ID == currentID {
		return nil
	}

	return appErrors.New(
		appErrors.ErrInvalidParams,
		http.StatusBadRequest,
		"project slug already exists",
	)
}

func (s *projectService) Update(
	ctx context.Context,
	id string,
	req dto.UpdateProjectRequest,
) (*model.Project, error) {

	// ---------------------------------------------
	// 1. 查询 Project
	// ---------------------------------------------

	project, err := s.repo.FindByID(
		ctx,
		id,
	)

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, appErrors.New(
				appErrors.ErrNotFound,
				http.StatusNotFound,
				"project not found",
			)
		}

		return nil, err
	}

	// ---------------------------------------------
	// 2. 基础参数
	// ---------------------------------------------

	name := strings.TrimSpace(req.Name)
	slug := strings.TrimSpace(req.Slug)

	if name == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"project name is required",
		)
	}

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"project slug is required",
		)
	}

	// ---------------------------------------------
	// 3. Status 校验
	// ---------------------------------------------

	if req.Status != 0 && req.Status != 1 {

		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"invalid project status",
		)
	}

	// ---------------------------------------------
	// 4. 唯一性检查
	// ---------------------------------------------

	if err := s.validateNameUnique(
		ctx,
		name,
		id,
	); err != nil {
		return nil, err
	}

	if err := s.validateSlugUnique(
		ctx,
		slug,
		id,
	); err != nil {
		return nil, err
	}

	// ---------------------------------------------
	// 5. 修改模型
	// ---------------------------------------------

	project.Name = name
	project.Slug = slug
	project.Description = req.Description
	project.Cover = req.Cover
	project.GitHubURL = req.GitHubURL
	project.DemoURL = req.DemoURL
	project.Featured = req.Featured
	project.Status = req.Status
	project.Sort = req.Sort

	// ---------------------------------------------
	// 6. 保存
	// ---------------------------------------------

	if err := s.repo.Update(
		ctx,
		project,
	); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) GetByID(
	ctx context.Context,
	id string,
) (*model.Project, error) {

	project, err := s.repo.FindByID(
		ctx,
		id,
	)

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, appErrors.New(
				appErrors.ErrNotFound,
				http.StatusNotFound,
				"project not found",
			)
		}

		return nil, err
	}

	return project, nil
}

func (s *projectService) GetBySlug(
	ctx context.Context,
	slug string,
) (*model.Project, error) {

	slug = strings.TrimSpace(slug)

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"project slug is required",
		)
	}

	project, err := s.repo.FindBySlug(
		ctx,
		slug,
	)

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, appErrors.New(
				appErrors.ErrNotFound,
				http.StatusNotFound,
				"project not found",
			)
		}

		return nil, err
	}

	return project, nil
}

func (s *projectService) ListPage(
	ctx context.Context,
	query repository.ProjectListQuery,
) ([]*model.Project, int64, error) {

	list, total, err := s.repo.FindPage(
		ctx,
		query,
	)

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (s *projectService) Delete(
	ctx context.Context,
	id string,
) error {

	// 先确认 Project 存在。
	_, err := s.GetByID(
		ctx,
		id,
	)

	if err != nil {
		return err
	}

	// 执行 Soft Delete。
	return s.repo.Delete(
		ctx,
		id,
	)
}
