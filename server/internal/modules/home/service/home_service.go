package service

import (
	"context"
	"encoding/json"
	"time"

	categorydto "leslie-blog-server/internal/modules/category/dto"
	categorymodel "leslie-blog-server/internal/modules/category/model"
	categoryservice "leslie-blog-server/internal/modules/category/service"
	"leslie-blog-server/internal/modules/home/dto"
	postdto "leslie-blog-server/internal/modules/post/dto"
	postmodel "leslie-blog-server/internal/modules/post/model"
	postrepository "leslie-blog-server/internal/modules/post/repository"
	postservice "leslie-blog-server/internal/modules/post/service"
	projectdto "leslie-blog-server/internal/modules/project/dto"
	projectmodel "leslie-blog-server/internal/modules/project/model"
	projectrepository "leslie-blog-server/internal/modules/project/repository"
	projectservice "leslie-blog-server/internal/modules/project/service"
	techstackdto "leslie-blog-server/internal/modules/techstack/dto"
	techstackmodel "leslie-blog-server/internal/modules/techstack/model"
	techstackservice "leslie-blog-server/internal/modules/techstack/service"
	"leslie-blog-server/internal/pkg/cache"
	"leslie-blog-server/internal/pkg/pagination"

	"golang.org/x/sync/errgroup"
)

// HomeService
//
// 负责组装 Blog 首页数据。
type HomeService interface {

	// GetHome 获取首页聚合数据。
	GetHome(
		ctx context.Context,
	) (*dto.HomeResponse, error)
}

type homeService struct {
	postService      postservice.PostService
	categoryService  categoryservice.CategoryService
	projectService   projectservice.ProjectService
	techStackService techstackservice.TechStackService
	cache            cache.Cache
}

func NewHomeService(
	postService postservice.PostService,
	categoryService categoryservice.CategoryService,
	projectService projectservice.ProjectService,
	techStackService techstackservice.TechStackService,
	cache cache.Cache,
) HomeService {

	return &homeService{
		postService: postService,

		categoryService: categoryService,

		projectService: projectService,

		techStackService: techStackService,
		cache:            cache,
	}
}

func (s *homeService) GetHome(
	ctx context.Context,
) (*dto.HomeResponse, error) {

	// -----------------------------------------
	// 1. 尝试读取缓存
	// -----------------------------------------

	cached, err := s.cache.Get(
		ctx,
		cache.KeyHome,
	)

	if err == nil {

		var result dto.HomeResponse

		if err := json.Unmarshal(
			[]byte(cached),
			&result,
		); err == nil {

			return &result, nil
		}

		// JSON 解析失败时，
		// 不直接让整个 Home API 失败。
		//
		// 因为缓存本身只是性能优化。
	}

	// 2. 缓存不存在 / 缓存异常
	//    从数据库重新加载。
	result, err := s.loadHome(ctx)

	if err != nil {
		return nil, err
	}

	// 3. 写入缓存
	data, err := json.Marshal(result)

	if err == nil {

		// 缓存写入失败不能影响主流程。
		_ = s.cache.Set(
			ctx,
			cache.KeyHome,
			string(data),
			5*time.Minute,
		)
	}

	// 4. 返回数据库结果
	return result, nil
}

func (s *homeService) loadHome(
	ctx context.Context,
) (*dto.HomeResponse, error) {

	var (
		featuredPosts []*postmodel.Post

		latestPosts []*postmodel.Post

		projects []*projectmodel.Project

		categories []*categorymodel.Category

		techStacks []*techstackmodel.TechStack
	)

	g, ctx := errgroup.WithContext(ctx)

	// ---------------------------------------------
	// Featured Posts
	// ---------------------------------------------

	g.Go(func() error {

		query := postrepository.PostListQuery{
			Params: pagination.NewParams(
				1,
				3,
			),
		}

		result, _, err :=
			s.postService.ListFeatured(
				ctx,
				query,
			)

		if err != nil {
			return err
		}

		featuredPosts = result

		return nil
	})

	// ---------------------------------------------
	// Latest Posts
	// ---------------------------------------------

	g.Go(func() error {

		query := postrepository.PostListQuery{
			Params: pagination.NewParams(
				1,
				6,
			),
		}

		result, _, err :=
			s.postService.ListPublished(
				ctx,
				query,
			)

		if err != nil {
			return err
		}

		latestPosts = result

		return nil
	})

	// ---------------------------------------------
	// Projects
	// ---------------------------------------------

	g.Go(func() error {
		featured := true
		query := projectrepository.ProjectListQuery{
			Params: pagination.NewParams(
				1,
				6,
			),
			Featured: &featured,
		}

		result, _, err :=
			s.projectService.ListPublic(
				ctx,
				query,
			)

		if err != nil {
			return err
		}

		projects = result

		return nil
	})

	// ---------------------------------------------
	// Categories
	// ---------------------------------------------

	g.Go(func() error {

		result, err :=
			s.categoryService.ListPublic(
				ctx,
			)

		if err != nil {
			return err
		}

		categories = result

		return nil
	})

	// ---------------------------------------------
	// TechStacks
	// ---------------------------------------------

	g.Go(func() error {

		result, err :=
			s.techStackService.ListPublic(
				ctx,
			)

		if err != nil {
			return err
		}

		techStacks = result

		return nil
	})

	// ---------------------------------------------
	// 等待所有任务结束
	// ---------------------------------------------

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// ---------------------------------------------
	// 所有查询成功
	// 开始组装响应
	// ---------------------------------------------

	return &dto.HomeResponse{
		FeaturedPosts: postdto.FromPublicModelList(
			featuredPosts,
		),

		LatestPosts: postdto.FromPublicModelList(
			latestPosts,
		),

		Categories: categorydto.FromSimpleModelList(
			categories,
		),

		Projects: projectdto.FromPublicModels(
			projects,
		),

		TechStacks: techstackdto.FromPublicModelList(
			techStacks,
		),
	}, nil
}
