package model

import (
	"time"

	techstackmodel "leslie-blog-server/internal/modules/techstack/model"
)

// Project 表示 Blog 中的一个项目。
//
// Project 主要用于：
//
// 1. Admin 后台管理
// 2. Blog 首页 Projects 区域展示
// 3. Project 详情页展示
type Project struct {

	// ID 使用 ULID。
	ID string `gorm:"column:id;type:char(26);primaryKey"`

	// 项目名称。
	Name string `gorm:"column:name;size:200;not null;uniqueIndex"`

	// URL 使用的唯一标识。
	Slug string `gorm:"column:slug;size:200;not null;uniqueIndex"`

	// 项目简介。
	Description *string `gorm:"column:description;size:1000"`

	// 项目封面。
	Cover *string `gorm:"column:cover;size:500"`

	// GitHub 地址。
	GitHubURL *string `gorm:"column:github_url;size:500"`

	// 在线 Demo 地址。
	DemoURL *string `gorm:"column:demo_url;size:500"`

	// 是否精选。
	//
	// 0 = false
	// 1 = true
	Featured bool `gorm:"column:featured;not null;default:false"`

	// 是否启用。
	//
	// 0 = disabled
	// 1 = enabled
	Status int8 `gorm:"column:status;not null;default:1"`

	// 排序值。
	//
	// 数值越小，排序越靠前。
	Sort int `gorm:"column:sort;not null;default:0"`

	CreatedAt time.Time `gorm:"column:created_at"`

	UpdatedAt time.Time `gorm:"column:updated_at"`

	DeletedAt *time.Time `gorm:"column:deleted_at;index"`

	TechStacks []*techstackmodel.TechStack `gorm:"many2many:project_tech_stacks;foreignKey:ID;joinForeignKey:ProjectID;References:ID;joinReferences:TechStackID"`
}

// TableName 指定数据库表名称。
func (Project) TableName() string {
	return "projects"
}
