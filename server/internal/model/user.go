package model

// User 后台用户
type User struct {
	BaseModel
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:128;not null" json:"-"`
	Nickname string `gorm:"size:64" json:"nickname"`
	Email    string `gorm:"size:128" json:"email"`
	Avatar   string `gorm:"size:512" json:"avatar"`
	Role     string `gorm:"size:32;default:admin" json:"role"`
}

// TableName 显式表名
func (User) TableName() string { return "users" }
