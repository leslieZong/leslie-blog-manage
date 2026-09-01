// Package model 数据库实体与 DTO 定义
package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 通用主键 + 时间戳
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// SoftBaseModel 带软删除
type SoftBaseModel struct {
	BaseModel
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
