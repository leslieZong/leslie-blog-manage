package model

import (
	"time"

	"gorm.io/datatypes"
)

// Post 文章（语言无关字段）。多语言内容存于 PostTranslation
type Post struct {
	BaseModel
	Slug            string     `gorm:"size:128;uniqueIndex;not null" json:"slug"`
	Cover           string     `gorm:"size:512" json:"cover"`
	CategoryID      uint       `gorm:"index;default:0" json:"categoryId"`
	Status          int8       `gorm:"index;default:0" json:"status"` // 0 草稿 1 已发布 2 已下线
	IsTop           bool       `gorm:"default:false" json:"isTop"`
	ViewCount       uint       `gorm:"default:0" json:"viewCount"`
	CommentCount    uint       `gorm:"default:0" json:"commentCount"`
	SEOTitle        string     `gorm:"column:seo_title;size:255" json:"seoTitle,omitempty"`
	SEODescription  string     `gorm:"column:seo_description;size:500" json:"seoDescription,omitempty"`
	SEOKeywords     string     `gorm:"column:seo_keywords;size:255" json:"seoKeywords,omitempty"`
	PublishedAt     *time.Time `json:"publishedAt,omitempty"`

	// 关联（Preload 填充）
	Category     Category           `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
	Tags         []PostTag          `gorm:"foreignKey:PostID;references:ID" json:"tags,omitempty"`
	Translations []PostTranslation  `gorm:"foreignKey:PostID;references:ID" json:"translations,omitempty"`
	// 当前 locale 的翻译（非 DB，按 locale 填充）
	Translation *PostTranslation `gorm:"-" json:"translation,omitempty"`
}

func (Post) TableName() string { return "posts" }

// PostTranslation 文章多语言翻译
type PostTranslation struct {
	BaseModel
	PostID  uint   `gorm:"uniqueIndex:uk_post_translations;not null" json:"postId"`
	Locale  string `gorm:"size:16;uniqueIndex:uk_post_translations;not null" json:"locale"` // zh-CN / en-US
	Title   string `gorm:"size:255;not null" json:"title"`
	Summary string `gorm:"size:500" json:"summary"`
	Content string `gorm:"type:longtext;not null" json:"content"`
}

func (PostTranslation) TableName() string { return "post_translations" }

// PostTag 文章标签
type PostTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"uniqueIndex:uk_post_tags;not null" json:"postId"`
	Name      string    `gorm:"size:64;uniqueIndex:uk_post_tags;not null" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func (PostTag) TableName() string { return "post_tags" }

// TagJSON 用于 projects.tech_stack 等 JSON 列的便捷类型别名
type TagJSON = datatypes.JSON
