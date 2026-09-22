package model

import "time"

// TechStack 表示一个技术栈。
//
// 例如：
//
// Vue
// TypeScript
// Go
// MySQL
// Docker
type TechStack struct {

	// 唯一 ID。
	ID string `gorm:"column:id;type:char(26);primaryKey"`

	// 技术名称。
	Name string `gorm:"column:name;size:100;not null;uniqueIndex"`

	// URL 友好的唯一标识。
	Slug string `gorm:"column:slug;size:100;not null;uniqueIndex"`

	// 图标地址。
	Icon *string `gorm:"column:icon;size:500"`

	// 技术描述。
	Description *string `gorm:"column:description;size:255"`

	// 官方网站。
	OfficialURL *string `gorm:"column:official_url;size:500"`

	// 技术分类。
	//
	// 例如：
	//
	// frontend
	// backend
	// database
	// devops
	Category *string `gorm:"column:category;size:50"`

	// 排序。
	Sort int `gorm:"column:sort;not null;default:0"`

	// 状态：
	//
	// 1 = 启用
	// 0 = 禁用
	Status int8 `gorm:"column:status;not null;default:1"`

	CreatedAt time.Time `gorm:"column:created_at"`

	UpdatedAt time.Time `gorm:"column:updated_at"`

	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

// TableName 指定数据库表名。
func (TechStack) TableName() string {
	return "tech_stacks"
}
