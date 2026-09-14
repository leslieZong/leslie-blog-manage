package dto

import (
	"time"

	"leslie-blog-server/internal/modules/post/model"
)

// PublicPostResponse
//
// 这是博客前台使用的文章详情响应。
//
// 注意：
// Public API 不应该直接复用 Admin 的 PostResponse。
//
// 原因：
// Admin 和 Public 是两个不同的数据边界。
type PublicPostResponse struct {
	// ID
	ID string `json:"id"`

	// 文章标题
	Title string `json:"title"`

	// URL 中使用的文章标识
	Slug string `json:"slug"`

	// 文章摘要
	Summary *string `json:"summary"`

	// Markdown / HTML / Rich Text 内容
	Content string `json:"content"`

	// 封面
	Cover *string `json:"cover"`

	// 发布时间
	PublishedAt *time.Time `json:"publishedAt"`

	// 阅读数量
	ViewCount uint64 `json:"viewCount"`
}

// PublicPostListItem
//
// Public API 的文章列表数据。
//
// 列表接口没有必要返回完整 Content。
type PublicPostListItem struct {
	ID string `json:"id"`

	Title string `json:"title"`

	Slug string `json:"slug"`

	Summary *string `json:"summary"`

	Cover *string `json:"cover"`

	PublishedAt *time.Time `json:"publishedAt"`

	ViewCount uint64 `json:"viewCount"`
}

// FromPublicModel
//
// 将数据库 Model 转换成 Public DTO。
func FromPublicModel(post *model.Post) *PublicPostResponse {
	if post == nil {
		return nil
	}

	return &PublicPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Summary:     post.Summary,
		Content:     post.Content,
		Cover:       post.Cover,
		PublishedAt: post.PublishedAt,
		ViewCount:   post.ViewCount,
	}
}

// FromPublicModelList
//
// 将文章 Model 列表转换成 Public DTO 列表。
func FromPublicModelList(posts []*model.Post) []*PublicPostListItem {
	result := make([]*PublicPostListItem, 0, len(posts))

	for _, post := range posts {
		if post == nil {
			continue
		}

		result = append(result, &PublicPostListItem{
			ID:          post.ID,
			Title:       post.Title,
			Slug:        post.Slug,
			Summary:     post.Summary,
			Cover:       post.Cover,
			PublishedAt: post.PublishedAt,
			ViewCount:   post.ViewCount,
		})
	}

	return result
}
