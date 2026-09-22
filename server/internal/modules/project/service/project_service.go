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
	techstackRepository "leslie-blog-server/internal/modules/techstack/repository"
	"leslie-blog-server/internal/pkg/database"
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

	// ListPublic 查询公开 Project。
	//
	// Public API 只能返回：
	//
	// status = 1
	//
	// 客户端不能通过 query 参数覆盖 status。
	ListPublic(
		ctx context.Context,
		query repository.ProjectListQuery,
	) ([]*model.Project, int64, error)

	// GetPublicByID 根据 ID 查询公开 Project。
	GetPublicByID(
		ctx context.Context,
		id string,
	) (*model.Project, error)

	// GetPublicBySlug 根据 slug 查询公开 Project。
	GetPublicBySlug(
		ctx context.Context,
		slug string,
	) (*model.Project, error)
}

type projectService struct {
	repo          repository.ProjectRepository
	techstackRepo techstackRepository.TechStackRepository

	transactionManager *database.TransactionManager
}

func NewProjectService(
	repo repository.ProjectRepository,
	techstackRepo techstackRepository.TechStackRepository,
	transactionManager *database.TransactionManager,
) ProjectService {

	return &projectService{
		repo:               repo,
		techstackRepo:      techstackRepo,
		transactionManager: transactionManager,
	}
}

func normalizeIDs(ids []string) []string {

	result := make(
		[]string,
		0,
		len(ids),
	)

	seen := make(
		map[string]struct{},
		len(ids),
	)

	for _, id := range ids {

		id = strings.TrimSpace(id)

		if id == "" {
			continue
		}

		if _, exists := seen[id]; exists {
			continue
		}

		seen[id] = struct{}{}

		result = append(
			result,
			id,
		)
	}

	return result
}

func (s *projectService) validateTechStacks(
	ctx context.Context,
	ids []string,
) error {

	ids = normalizeIDs(ids)

	if len(ids) == 0 {
		return nil
	}

	techStacks, err := s.techstackRepo.FindByIDs(
		ctx,
		ids,
	)

	if err != nil {
		return err
	}

	if len(techStacks) != len(ids) {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"one or more tech stacks are invalid",
		)
	}

	return nil
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

	// 先处理 TechStack ID。
	techStackIDs := normalizeIDs(
		req.TechStackIDs,
	)

	// 验证 TechStack 是否全部存在并且启用。
	if err := s.validateTechStacks(
		ctx,
		techStackIDs,
	); err != nil {
		return nil, err
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
	// 5.1 保存 Project
	// ---------------------------------------------

	err := s.transactionManager.WithTransaction(
		ctx,
		func(tx *gorm.DB) error {
			transactionProjectRepo :=
				repository.NewProjectRepository(tx)

			transactionRelationRepo :=
				repository.NewProjectTechStackRepository(tx)

			if err := transactionProjectRepo.Create(
				ctx,
				project,
			); err != nil {
				return err
			}

			if err := transactionRelationRepo.ReplaceTechStacks(
				ctx,
				project.ID,
				req.TechStackIDs,
			); err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	project, err = s.repo.FindByID(
		ctx,
		project.ID,
	)
	if err != nil {
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

	techStackIDs := normalizeIDs(
		req.TechStackIDs,
	)

	if err := s.validateTechStacks(
		ctx,
		techStackIDs,
	); err != nil {
		return nil, err
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
	// 6. 开始事务
	// ---------------------------------------------

	err = s.transactionManager.WithTransaction(
		ctx,
		func(tx *gorm.DB) error {

			projectRepo :=
				repository.NewProjectRepository(tx)

			projectTechStackRepo :=
				repository.NewProjectTechStackRepository(tx)

			// 更新 Project。
			if err := projectRepo.Update(
				ctx,
				project,
			); err != nil {
				return err
			}

			// 替换 TechStack。
			if err := projectTechStackRepo.ReplaceTechStacks(
				ctx,
				project.ID,
				techStackIDs,
			); err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}
	project, err = s.repo.FindByID(
		ctx,
		project.ID,
	)
	if err != nil {
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

func (s *projectService) ListPublic(
	ctx context.Context,
	query repository.ProjectListQuery,
) ([]*model.Project, int64, error) {

	// Public API 不允许客户端控制 status。
	//
	// 无论客户端传：
	//
	// ?status=0
	// ?status=1
	//
	// 这里最终都强制：
	//
	// status = 1
	status := int8(1)

	query.Status = &status

	list, total, err := s.repo.FindPage(
		ctx,
		query,
	)

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (s *projectService) GetPublicByID(
	ctx context.Context,
	id string,
) (*model.Project, error) {

	project, err := s.repo.FindPublicByID(
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

func (s *projectService) GetPublicBySlug(
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

	project, err := s.repo.FindPublicBySlug(
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
