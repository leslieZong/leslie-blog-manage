package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 表示后台系统中的一个角色。
//
// 例如：
//
// admin
// editor
// viewer
//
// 注意：
// Role 本身只描述“角色这个业务实体”。
// 具体有哪些权限，由 Casbin 管理。
type Role struct {
	// ID 是角色唯一标识。
	//
	// Leslie Blog 统一使用 ULID。
	ID string `gorm:"column:id;type:char(26);primaryKey"`

	// Name 是角色内部名称。
	//
	// 例如：
	//
	// admin
	// editor
	// viewer
	//
	// 程序内部使用这个值。
	Name string `gorm:"column:name;type:varchar(50);uniqueIndex;not null"`

	// DisplayName 是给管理员看的名称。
	//
	// 例如：
	//
	// Name:
	// editor
	//
	// DisplayName:
	// 内容编辑
	DisplayName string `gorm:"column:display_name;type:varchar(100);not null"`

	// Description 是角色描述。
	Description string `gorm:"column:description;type:varchar(255);not null"`

	// Status：
	//
	// 1 = 启用
	// 0 = 禁用
	Status int8 `gorm:"column:status;not null;default:1"`

	// 创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`

	// 更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`

	// DeletedAt 用于 GORM 软删除。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 指定 Role 对应的数据库表。
func (Role) TableName() string {
	return "roles"
}
