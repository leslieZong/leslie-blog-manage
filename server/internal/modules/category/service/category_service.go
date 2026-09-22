package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/category/model"
	"leslie-blog-server/internal/modules/category/repository"

	"leslie-blog-server/internal/pkg/cache"
	"leslie-blog-server/internal/pkg/ulid"

	"gorm.io/gorm"
)

// CategoryService
//
// Category 业务逻辑层。
type CategoryService interface {

	// 创建分类。
	Create(
		ctx context.Context,
		name string,
		slug string,
		description *string,
		sort int,
	) (*model.Category, error)

	// 根据 ID 获取分类。
	GetByID(
		ctx context.Context,
		id string,
	) (*model.Category, error)

	// 查询全部分类。
	List(
		ctx context.Context,
	) ([]*model.Category, error)

	// 修改分类。
	Update(
		ctx context.Context,
		id string,
		name string,
		slug string,
		description *string,
		sort int,
		status int8,
	) (*model.Category, error)

	// 删除分类。
	Delete(
		ctx context.Context,
		id string,
	) error

	ListPage(
		ctx context.Context,
		query repository.CategoryListQuery,
	) ([]*model.Category, int64, error)

	// Public
	ListPublic(
		ctx context.Context,
	) ([]*model.Category, error)

	GetPublicBySlug(
		ctx context.Context,
		slug string,
	) (*model.Category, error)
}

// categoryService
//
// CategoryService 的具体实现。
type categoryService struct {
	repo  repository.CategoryRepository
	cache cache.Cache
}

// NewCategoryService
//
// 创建 Category Service。
func NewCategoryService(
	repo repository.CategoryRepository,
	cache cache.Cache,
) CategoryService {
	return &categoryService{
		repo:  repo,
		cache: cache,
	}
}

// Create
//
// 创建分类。
func (s *categoryService) Create(
	ctx context.Context,
	name string,
	slug string,
	description *string,
	sort int,
) (*model.Category, error) {

	// --------------------------------------------------------
	// 第一步：清理输入
	// --------------------------------------------------------

	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)

	// --------------------------------------------------------
	// 第二步：参数校验
	// --------------------------------------------------------

	if name == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category name cannot be empty",
		)
	}

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category slug cannot be empty",
		)
	}

	// --------------------------------------------------------
	// 第三步：检查 name 是否重复
	// --------------------------------------------------------

	existing, err := s.repo.FindByName(ctx, name)

	if err == nil && existing != nil {
		return nil, appErrors.New(
			appErrors.ErrCategoryNameExists,
			http.StatusBadRequest,
			"category name already exists",
		)
	}

	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to check category name",
			err,
		)
	}

	// --------------------------------------------------------
	// 第四步：检查 slug 是否重复
	// --------------------------------------------------------

	existing, err = s.repo.FindBySlug(ctx, slug)

	if err == nil && existing != nil {
		return nil, appErrors.New(
			appErrors.ErrCategorySlugExists,
			http.StatusBadRequest,
			"category slug already exists",
		)
	}

	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to check category slug",
			err,
		)
	}

	// --------------------------------------------------------
	// 第五步：创建 Model
	// --------------------------------------------------------

	category := &model.Category{
		ID:          ulid.New(),
		Name:        name,
		Slug:        slug,
		Description: description,
		Sort:        sort,

		// 新分类默认启用。
		Status: 1,
	}

	// --------------------------------------------------------
	// 第六步：写入数据库
	// --------------------------------------------------------

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to create category",
			err,
		)
	}
	// 删除 Home Cache。
	if err := cache.InvalidateHome(
		ctx,
		s.cache,
	); err != nil {
		// 记录日志
	}

	return category, nil
}

// GetByID
//
// 根据 ID 查询分类。
func (s *categoryService) GetByID(
	ctx context.Context,
	id string,
) (*model.Category, error) {

	id = strings.TrimSpace(id)

	if id == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category id cannot be empty",
		)
	}

	category, err := s.repo.FindByID(ctx, id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.New(
			appErrors.ErrNotFound,
			http.StatusNotFound,
			"category not found",
		)
	}

	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to query category",
			err,
		)
	}

	return category, nil
}

// List
//
// 查询所有未删除分类。
func (s *categoryService) List(
	ctx context.Context,
) ([]*model.Category, error) {

	categories, err := s.repo.FindAll(ctx)

	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to query categories",
			err,
		)
	}

	return categories, nil
}

func (s *categoryService) validateNameUnique(
	ctx context.Context,
	name string,
	currentID string,
) error {

	category, err := s.repo.FindByName(
		ctx,
		name,
	)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	if category.ID != currentID {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category name already exists",
		)
	}

	return nil
}

func (s *categoryService) validateSlugUnique(
	ctx context.Context,
	slug string,
	currentID string,
) error {

	category, err := s.repo.FindBySlug(
		ctx,
		slug,
	)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	if category.ID != currentID {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category slug already exists",
		)
	}

	return nil
}

// Update
//
// 修改分类。
func (s *categoryService) Update(
	ctx context.Context,
	id string,
	name string,
	slug string,
	description *string,
	sort int,
	status int8,
) (*model.Category, error) {

	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)

	if id == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category id cannot be empty",
		)
	}

	if name == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category name cannot be empty",
		)
	}

	if slug == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category slug cannot be empty",
		)
	}

	if status != 0 && status != 1 {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"invalid category status",
		)
	}

	// --------------------------------------------------------
	// 查询原分类
	// --------------------------------------------------------

	category, err := s.repo.FindByID(ctx, id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.New(
			appErrors.ErrNotFound,
			http.StatusNotFound,
			"category not found",
		)
	}

	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to query category",
			err,
		)
	}

	// --------------------------------------------------------
	// 检查新的 name 是否和其他分类冲突
	// --------------------------------------------------------

	if err := s.validateNameUnique(ctx, name, id); err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// 检查新的 slug 是否和其他分类冲突
	// --------------------------------------------------------

	if err := s.validateSlugUnique(ctx, slug, id); err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// 更新允许修改的字段
	// --------------------------------------------------------

	category.Name = name
	category.Slug = slug
	category.Description = description
	category.Sort = sort
	category.Status = status

	// --------------------------------------------------------
	// 保存
	// --------------------------------------------------------

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to update category",
			err,
		)
	}
	// 删除 Home Cache。
	if err := cache.InvalidateHome(
		ctx,
		s.cache,
	); err != nil {
		// 记录日志
	}

	return category, nil
}

// Delete
//
// 删除分类。
//
// 业务规则：
//
// 1. 分类必须存在
// 2. 分类下面不能还有文章
// 3. 使用软删除
func (s *categoryService) Delete(
	ctx context.Context,
	id string,
) error {

	id = strings.TrimSpace(id)

	if id == "" {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			http.StatusBadRequest,
			"category id cannot be empty",
		)
	}

	// --------------------------------------------------------
	// 查询分类
	// --------------------------------------------------------

	_, err := s.repo.FindByID(ctx, id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return appErrors.New(
			appErrors.ErrNotFound,
			http.StatusNotFound,
			"category not found",
		)
	}

	if err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to query category",
			err,
		)
	}

	// --------------------------------------------------------
	// 查询这个分类下面是否还有文章
	// --------------------------------------------------------

	count, err := s.repo.CountPosts(ctx, id)

	if err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to count category posts",
			err,
		)
	}

	// --------------------------------------------------------
	// 有文章则禁止删除
	// --------------------------------------------------------

	if count > 0 {
		return appErrors.New(
			appErrors.ErrCategoryHasPosts,
			http.StatusBadRequest,
			"category contains posts and cannot be deleted",
		)
	}

	// --------------------------------------------------------
	// 执行软删除
	// --------------------------------------------------------

	if err := s.repo.Delete(ctx, id); err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternalServer,
			http.StatusInternalServerError,
			"failed to delete category",
			err,
		)
	}
	// 删除 Home Cache。
	if err := cache.InvalidateHome(
		ctx,
		s.cache,
	); err != nil {
		// 记录日志
	}

	return nil
}

func (s *categoryService) ListPage(
	ctx context.Context,
	query repository.CategoryListQuery,
) ([]*model.Category, int64, error) {

	list, total, err := s.repo.FindPage(
		ctx,
		query,
	)

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (s *categoryService) ListPublic(
	ctx context.Context,
) ([]*model.Category, error) {

	categories, err := s.repo.FindPublicAll(
		ctx,
	)

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *categoryService) GetPublicBySlug(
	ctx context.Context,
	slug string,
) (*model.Category, error) {

	category, err := s.repo.FindPublicBySlug(
		ctx,
		slug,
	)

	if err != nil {
		return nil, err
	}

	return category, nil
}
