package model

// Category 分类
type Category struct {
	BaseModel
	Name        string `gorm:"size:64;not null" json:"name"`
	Slug        string `gorm:"size:64;uniqueIndex;not null" json:"slug"`
	Description string `gorm:"size:255" json:"description"`
	ParentID    uint   `gorm:"index;default:0" json:"parentId"`
	SortOrder   int    `gorm:"default:0" json:"sortOrder"`
	// 关联文章数（查询时填充，非 DB 字段）
	PostCount int64 `gorm:"-" json:"postCount,omitempty"`
}

func (Category) TableName() string { return "categories" }
