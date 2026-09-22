package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/techstack/dto"
	"leslie-blog-server/internal/modules/techstack/model"
	"leslie-blog-server/internal/modules/techstack/repository"
	"leslie-blog-server/internal/pkg/ulid"

	"gorm.io/gorm"
)

// TechStackService 负责 TechStack 业务逻辑。
type TechStackService interface {

	// Create 创建 TechStack。
	Create(
		ctx context.Context,
		req dto.CreateTechStackRequest,
	) (*model.TechStack, error)

	// GetByID 获取 TechStack。
	GetByID(
		ctx context.Context,
		id string,
	) (*model.TechStack, error)

	// List 获取后台 TechStack。
	List(
		ctx context.Context,
	) ([]*model.TechStack, error)

	// ListPublic 获取公开 TechStack。
	ListPublic(
		ctx context.Context,
	) ([]*model.TechStack, error)

	// GetPublicBySlug 获取公开 TechStack。
	GetPublicBySlug(
		ctx context.Context,
		slug string,
	) (*model.TechStack, error)

	// Update 更新 TechStack。
	Update(
		ctx context.Context,
		id string,
		req dto.UpdateTechStackRequest,
	) (*model.TechStack, error)

	// Delete 删除 TechStack。
	Delete(
		ctx context.Context,
		id string,
	) error
}

type techStackService struct {
	repo repository.TechStackRepository
}

func NewTechStackService(
	repo repository.TechStackRepository,
) TechStackService {

	return &techStackService{
		repo: repo,
	}
}

func (s *techStackService) Create(
	ctx context.Context,
	req dto.CreateTechStackRequest,
) (*model.TechStack, error) {

	name := strings.TrimSpace(req.Name)
	slug := strings.TrimSpace(req.Slug)

	if name == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tech stack name is required",
		)
	}

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tech stack slug is required",
		)
	}

	// 检查 name 是否已经存在。
	existing, err := s.repo.FindByName(
		ctx,
		name,
	)

	if err == nil && existing != nil {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tech stack name already exists",
		)
	}

	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, err
	}

	// 检查 slug。
	existing, err = s.repo.FindBySlug(
		ctx,
		slug,
	)

	if err == nil && existing != nil {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tech stack slug already exists",
		)
	}

	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, err
	}

	techStack := &model.TechStack{

		// 使用项目现有 ULID 工具。
		ID: ulid.New(),

		Name: name,

		Slug: slug,

		Icon: req.Icon,

		Description: req.Description,

		OfficialURL: req.OfficialURL,

		Category: req.Category,

		Sort: req.Sort,

		Status: 1,
	}

	if err := s.repo.Create(
		ctx,
		techStack,
	); err != nil {

		return nil, err
	}

	return techStack, nil
}

func (s *techStackService) GetByID(
	ctx context.Context,
	id string,
) (*model.TechStack, error) {

	techStack, err := s.repo.FindByID(
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
				"tech stack not found",
			)
		}

		return nil, err
	}

	return techStack, nil
}

func (s *techStackService) List(
	ctx context.Context,
) ([]*model.TechStack, error) {

	return s.repo.FindAll(ctx)
}

func (s *techStackService) ListPublic(
	ctx context.Context,
) ([]*model.TechStack, error) {

	return s.repo.FindPublic(ctx)
}

func (s *techStackService) GetPublicBySlug(
	ctx context.Context,
	slug string,
) (*model.TechStack, error) {

	slug = strings.TrimSpace(slug)

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tech stack slug is required",
		)
	}

	techStack, err := s.repo.FindBySlug(
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

				"tech stack not found",
			)
		}

		return nil, err
	}

	// Public API 不允许查看禁用 TechStack。
	if techStack.Status != 1 {
		return nil, appErrors.New(
			appErrors.ErrNotFound,
			http.StatusNotFound,

			"tech stack not found",
		)
	}

	return techStack, nil
}

func (s *techStackService) Delete(
	ctx context.Context,
	id string,
) error {

	techStack, err := s.repo.FindByID(
		ctx,
		id,
	)

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return appErrors.New(
				appErrors.ErrNotFound,
				http.StatusNotFound,
				"tech stack not found",
			)
		}

		return err
	}

	// 统计当前有多少 Project 正在使用。
	count, err := s.repo.CountProjects(
		ctx,
		techStack.ID,
	)

	if err != nil {
		return err
	}

	// 如果仍然被 Project 使用，则禁止删除。
	if count > 0 {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tech stack is still used by projects",
		)
	}

	// 没有任何 Project 使用，可以软删除。
	return s.repo.Delete(
		ctx,
		techStack.ID,
	)
}

func (s *techStackService) Update(
	ctx context.Context,
	id string,
	req dto.UpdateTechStackRequest,
) (*model.TechStack, error) {

	techStack, err := s.repo.FindByID(
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

				"tech stack not found",
			)
		}

		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	slug := strings.TrimSpace(req.Slug)

	if name == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,

			"tech stack name is required",
		)
	}

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,

			"tech stack slug is required",
		)
	}

	// 检查名称冲突。
	existing, err := s.repo.FindByName(
		ctx,
		name,
	)

	if err == nil &&
		existing != nil &&
		existing.ID != id {

		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,

			"tech stack name already exists",
		)
	}

	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, err
	}

	// 检查 slug 冲突。
	existing, err = s.repo.FindBySlug(
		ctx,
		slug,
	)

	if err == nil &&
		existing != nil &&
		existing.ID != id {

		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"tech stack slug already exists",
		)
	}

	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, err
	}

	techStack.Name = name
	techStack.Slug = slug
	techStack.Icon = req.Icon
	techStack.Description = req.Description
	techStack.OfficialURL = req.OfficialURL
	techStack.Category = req.Category
	techStack.Sort = req.Sort
	techStack.Status = req.Status

	if err := s.repo.Update(
		ctx,
		techStack,
	); err != nil {
		return nil, err
	}

	return techStack, nil
}
