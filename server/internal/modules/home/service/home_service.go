package service

import (
	"context"

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
}

func NewHomeService(
	postService postservice.PostService,
	categoryService categoryservice.CategoryService,
	projectService projectservice.ProjectService,
	techStackService techstackservice.TechStackService,
) HomeService {

	return &homeService{
		postService: postService,

		categoryService: categoryService,

		projectService: projectService,

		techStackService: techStackService,
	}
}

func (s *homeService) GetHome(
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
