package model

import (
	"time"

	"gorm.io/gorm"
)

// Post 表示博客文章。
//
// Model 的职责：
//
// 1. 描述数据库 posts 表
// 2. 描述文章在 Go 中的数据结构
//
// 注意：
// Model 不负责业务逻辑。
// 不应该在 Model 里面写：
// - 创建文章
// - 发布文章
// - 权限判断
// - 数据库查询
type Post struct {
	// ID 是文章唯一标识。
	//
	// 使用 string 保存 ULID。
	ID string `gorm:"column:id;primaryKey"`

	// Title 是文章标题。
	Title string `gorm:"column:title;size:200;not null"`

	// Slug 是文章 URL 标识。
	//
	// 数据库中要求唯一。
	Slug string `gorm:"column:slug;size:200;uniqueIndex;not null"`

	// Summary 是文章摘要。
	//
	// 数据库允许 NULL，所以使用 *string。
	Summary *string `gorm:"column:summary;size:500"`

	// Content 是 Markdown 原文。
	//
	// 数据库使用 LONGTEXT。
	Content string `gorm:"column:content;type:longtext;not null"`

	// Cover 是文章封面地址。
	//
	// 没有封面时可以为 NULL。
	Cover *string `gorm:"column:cover;size:500"`

	// Status 是文章状态。
	//
	// draft
	// published
	// archived
	Status string `gorm:"column:status;size:20;not null;default:draft"`

	// AuthorID 表示文章作者。
	//
	// 对应 users.id。
	AuthorID string `gorm:"column:author_id;size:26;not null"`

	// PublishedAt 表示文章第一次发布的时间。
	//
	// 草稿状态下可以为 NULL。
	PublishedAt *time.Time `gorm:"column:published_at"`

	// ViewCount 是阅读次数。
	ViewCount uint64 `gorm:"column:view_count;not null;default:0"`

	// CreatedAt 是创建时间。
	CreatedAt time.Time `gorm:"column:created_at"`

	// UpdatedAt 是最后修改时间。
	UpdatedAt time.Time `gorm:"column:updated_at"`

	// DeletedAt 是软删除时间。
	//
	// NULL 表示没有被删除。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 指定 Model 对应的数据库表。
func (Post) TableName() string {
	return "posts"
}
