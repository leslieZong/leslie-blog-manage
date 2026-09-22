package dto

import (
	categorydto "leslie-blog-server/internal/modules/category/dto"
	postdto "leslie-blog-server/internal/modules/post/dto"
	projectdto "leslie-blog-server/internal/modules/project/dto"
	techstackdto "leslie-blog-server/internal/modules/techstack/dto"
)

// HomeResponse
//
// Leslie Blog 首页所需要的全部数据。
//
// 注意：
// 这个 DTO 不是数据库 Model。
//
// 它只是一个 API 返回结构。
type HomeResponse struct {

	// FeaturedPosts
	//
	// 首页精选文章。
	FeaturedPosts []*postdto.PublicPostListItem `json:"featuredPosts"`

	// LatestPosts
	//
	// 首页最新文章。
	LatestPosts []*postdto.PublicPostListItem `json:"latestPosts"`

	// Categories
	//
	// 首页 Topics / Categories 数据。
	Categories []*categorydto.SimpleCategoryResponse `json:"categories"`

	// Projects
	//
	// 首页 Projects 数据。
	Projects []*projectdto.PublicProjectResponse `json:"projects"`

	// TechStacks
	//
	// 首页技术栈数据。
	TechStacks []*techstackdto.PublicTechStackResponse `json:"techStacks"`
}
