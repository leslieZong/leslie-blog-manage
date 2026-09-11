package dto

import (
	"leslie-blog-server/internal/modules/post/model"
	"time"
)

// CreatePostRequest 创建文章请求。
//
// 这个结构体描述：
// “前端创建文章时，允许提交哪些字段”。
type CreatePostRequest struct {

	// Title 文章标题。
	Title string `json:"title"`

	// Slug 文章 URL 标识。
	//
	// 例如：
	// https://example.com/posts/go-learning
	//
	// slug 就是：
	// go-learning
	Slug string `json:"slug"`

	// Summary 文章摘要。
	//
	// 使用 *string 是为了区分：
	//
	// 1. 没有传
	// 2. 传了空字符串
	//
	// 当前阶段不需要过度纠结，
	// 后续 Update 时这个区别会更加明显。
	Summary *string `json:"summary"`

	// Content 文章正文。
	Content string `json:"content"`

	// Cover 文章封面地址。
	Cover *string `json:"cover"`
}

// UpdatePostRequest 更新文章请求。
//
// 注意：
// Update 和 Create 的字段目前基本一致，
// 但是这里仍然建议定义成独立 DTO。
//
// 不建议直接：
// type UpdatePostRequest = CreatePostRequest
//
// 因为未来两个接口的业务字段很可能不同。
type UpdatePostRequest struct {

	// Title 文章标题。
	Title string `json:"title"`

	// Slug 文章 URL 标识。
	Slug string `json:"slug"`

	// Summary 文章摘要。
	Summary *string `json:"summary"`

	// Content 文章正文。
	Content string `json:"content"`

	// Cover 文章封面。
	Cover *string `json:"cover"`
}

// PostResponse 文章详情响应。
//
// 注意：
// 这是 API 对外暴露的数据结构，
// 不是数据库 Model。
type PostResponse struct {

	// ID 文章 ID。
	ID string `json:"id"`

	// Title 文章标题。
	Title string `json:"title"`

	// Slug 文章 slug。
	Slug string `json:"slug"`

	// Summary 文章摘要。
	Summary *string `json:"summary"`

	// Content 文章正文。
	Content string `json:"content"`

	// Cover 文章封面。
	Cover *string `json:"cover"`

	// Status 当前文章状态。
	//
	// draft
	// published
	// archived
	Status string `json:"status"`

	// AuthorID 作者 ID。
	AuthorID string `json:"author_id"`

	// PublishedAt 发布时间。
	PublishedAt *time.Time `json:"published_at"`

	// ViewCount 阅读量。
	ViewCount uint64 `json:"view_count"`

	// CreatedAt 创建时间。
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 更新时间。
	UpdatedAt time.Time `json:"updated_at"`
}

// FromModel 将数据库 Model 转换成 API Response DTO。
func FromModel(post *model.Post) *PostResponse {

	if post == nil {
		return nil
	}

	return &PostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Summary:     post.Summary,
		Content:     post.Content,
		Cover:       post.Cover,
		Status:      post.Status,
		AuthorID:    post.AuthorID,
		PublishedAt: post.PublishedAt,
		ViewCount:   post.ViewCount,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}

// PostListItem 文章列表项。
//
// 列表通常不需要返回完整 Content，
// 避免一次查询大量正文。
type PostListItem struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Summary     *string    `json:"summary"`
	Cover       *string    `json:"cover"`
	Status      string     `json:"status"`
	AuthorID    string     `json:"author_id"`
	PublishedAt *time.Time `json:"published_at"`
	ViewCount   uint64     `json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// FromModelList 将文章 Model 列表转换成 API 列表 DTO。
func FromModelList(posts []*model.Post) []*PostListItem {

	result := make([]*PostListItem, 0, len(posts))

	for _, post := range posts {

		if post == nil {
			continue
		}

		result = append(result, &PostListItem{
			ID:          post.ID,
			Title:       post.Title,
			Slug:        post.Slug,
			Summary:     post.Summary,
			Cover:       post.Cover,
			Status:      post.Status,
			AuthorID:    post.AuthorID,
			PublishedAt: post.PublishedAt,
			ViewCount:   post.ViewCount,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
		})
	}

	return result
}
