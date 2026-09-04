package model

import (
	"time"

	"gorm.io/gorm"
)

// User 表示后台系统中的用户。
//
// 它与数据库中的 users 表对应。
//
// 注意：
//
// Model 的主要职责是描述“数据是什么样子”，
// 而不是处理 HTTP 请求，也不是处理具体业务逻辑。
type User struct {

	// --------------------------------------------------
	// ID
	// --------------------------------------------------
	//
	// 对应数据库：
	//
	// id CHAR(26)
	//
	// 这里使用 string。
	//
	// 因为我们后面计划使用 ULID 作为主键。
	ID string `gorm:"column:id;type:char(26);primaryKey"`

	// --------------------------------------------------
	// Username
	// --------------------------------------------------
	//
	// 对应：
	//
	// username VARCHAR(50)
	//
	// 用户登录时使用。
	Username string `gorm:"column:username;type:varchar(50);uniqueIndex;not null"`

	// --------------------------------------------------
	// Email
	// --------------------------------------------------
	//
	// 对应：
	//
	// email VARCHAR(255)
	//
	// 可以用于用户邮箱。
	Email *string `gorm:"column:email;type:varchar(255);uniqueIndex"`

	// --------------------------------------------------
	// PasswordHash
	// --------------------------------------------------
	//
	// 对应：
	//
	// password_hash VARCHAR(255)
	//
	// 注意：
	//
	// 这里保存的是密码哈希，
	// 绝对不是用户的明文密码。
	PasswordHash string `gorm:"column:password_hash;type:varchar(255);not null"`

	// --------------------------------------------------
	// DisplayName
	// --------------------------------------------------
	//
	// 用户在管理后台显示的名称。
	DisplayName string `gorm:"column:display_name;type:varchar(100);not null"`

	// --------------------------------------------------
	// AvatarURL
	// --------------------------------------------------
	//
	// 用户头像地址。
	AvatarURL string `gorm:"column:avatar_url;type:varchar(500);not null"`

	// --------------------------------------------------
	// Status
	// --------------------------------------------------
	//
	// 1 = 正常
	// 0 = 禁用
	//
	// 登录时需要检查用户是否被禁用。
	Status int8 `gorm:"column:status;not null;default:1"`

	// --------------------------------------------------
	// LastLoginAt
	// --------------------------------------------------
	//
	// 最后登录时间。
	//
	// 使用指针是因为：
	//
	// 新创建的用户可能从来没有登录过。
	//
	// 所以它可以是：
	//
	// nil
	//
	// 或者：
	//
	// 2026-09-02 18:00:00
	LastLoginAt *time.Time `gorm:"column:last_login_at"`

	// --------------------------------------------------
	// CreatedAt
	// --------------------------------------------------
	//
	// 创建时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`

	// --------------------------------------------------
	// UpdatedAt
	// --------------------------------------------------
	//
	// 最后更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`

	// --------------------------------------------------
	// DeletedAt
	// --------------------------------------------------
	//
	// GORM 的软删除字段。
	//
	// 如果 DeletedAt 有值，
	// 表示这个用户已经被软删除。
	//
	// GORM 查询时默认会过滤掉已经软删除的数据。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 指定 User 对应的数据库表名称。
//
// GORM 默认会根据 struct 名称推断表名：
//
// User
// ↓
// users
//
// 虽然默认情况下其实已经可以得到 users，
// 但这里显式写出来有两个好处：
//
// 1. 新手更容易理解 Model 与表的关系
// 2. 如果以后表名发生变化，只需要修改这里
func (User) TableName() string {
	return "users"
}
