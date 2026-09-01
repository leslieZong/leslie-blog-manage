package model

// Comment 评论
type Comment struct {
	BaseModel
	PostID   uint   `gorm:"index;not null" json:"postId"`
	Author   string `gorm:"size:64;not null" json:"author"`
	Email    string `gorm:"size:128" json:"email"`
	Avatar   string `gorm:"size:512" json:"avatar"`
	Content  string `gorm:"type:text;not null" json:"content"`
	ParentID uint   `gorm:"index;default:0" json:"parentId"`
	Status   int8   `gorm:"index;default:0" json:"status"` // 0 待审 1 通过 2 拒绝
	IP       string `gorm:"size:64" json:"ip"`

	// 关联文章标题（查询时填充）
	PostTitle string `gorm:"-" json:"postTitle,omitempty"`
}

func (Comment) TableName() string { return "comments" }
