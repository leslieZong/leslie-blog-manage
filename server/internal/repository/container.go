package repository

import "gorm.io/gorm"

// Container 仓储容器，集中持有各域仓储实例
type Container struct {
	User      *UserRepo
	Category  *CategoryRepo
	Post      *PostRepo
	Project   *ProjectRepo
	TechStack *TechStackRepo
	Media     *MediaRepo
	Comment   *CommentRepo
	Settings  *SettingsRepo
}

// withDB 公共父结构
type withDB struct{ db *gorm.DB }

// helpers
func paginate(db *gorm.DB, page, pageSize int) *gorm.DB {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return db.Offset((page - 1) * pageSize).Limit(pageSize)
}
