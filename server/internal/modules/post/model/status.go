package model

// PostStatus 定义文章状态。
//
// 使用 string 而不是 int，主要是为了：
// 1. 数据库可读
// 2. API 返回更直观
// 3. 后续调试方便
type PostStatus string

const (
	// PostStatusDraft 表示草稿。
	PostStatusDraft PostStatus = "draft"

	// PostStatusPublished 表示已发布。
	PostStatusPublished PostStatus = "published"

	// PostStatusArchived 表示已归档。
	PostStatusArchived PostStatus = "archived"
)

// IsValid 判断文章状态是否合法。
func (s PostStatus) IsValid() bool {
	switch s {
	case PostStatusDraft,
		PostStatusPublished,
		PostStatusArchived:
		return true

	default:
		return false
	}
}
