package repository

import "leslie-blog-server/internal/pkg/pagination"

// PostListQuery 是 Repository 层的文章查询条件。
//
// 注意：
//
// 它不是 DTO。
//
// DTO 属于 HTTP 层。
// Query 属于 Repository 查询层。
type PostListQuery struct {

	// 分页参数。
	pagination.Params

	// 分类 ID。
	CategoryID string

	// 文章状态。
	Status string

	// 标签 ID。
	TagID string
}
