package model

import "gorm.io/datatypes"

// Project 项目
type Project struct {
	BaseModel
	Name        string         `gorm:"size:128;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Cover       string         `gorm:"size:512" json:"cover"`
	DemoURL     string         `gorm:"column:demo_url;size:512" json:"demoUrl"`
	RepoURL     string         `gorm:"column:repo_url;size:512" json:"repoUrl"`
	TechStack   datatypes.JSON `gorm:"type:json" json:"techStack,omitempty"`
	Status      int8           `gorm:"default:0" json:"status"` // 0 进行中 1 已完成 2 已归档
	SortOrder   int            `gorm:"default:0" json:"sortOrder"`
}

func (Project) TableName() string { return "projects" }
