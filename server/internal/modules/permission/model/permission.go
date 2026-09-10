package model

import (
	"time"

	"gorm.io/gorm"
)

// Permission 表示系统中的一个权限定义。
//
// 注意：
//
// Permission 本身只描述：
// “系统存在什么权限”
//
// 它不负责描述：
// “哪个角色拥有这个权限”
//
// 角色与权限之间的关系由 Casbin 管理。
type Permission struct {
	ID          string  `gorm:"column:id;primaryKey"`
	Name        string  `gorm:"column:name;uniqueIndex;size:100;not null"`
	DisplayName string  `gorm:"column:display_name;size:100;not null"`
	Resource    string  `gorm:"column:resource;size:50;not null"`
	Action      string  `gorm:"column:action;size:50;not null"`
	Description *string `gorm:"column:description;size:255"`
	Status      int8    `gorm:"column:status;not null;default:1"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 指定数据库表名。
func (Permission) TableName() string {
	return "permissions"
}
