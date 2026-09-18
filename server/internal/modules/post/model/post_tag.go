package model

// PostTag 表示文章和标签之间的关联。
//
// 对应数据库：
//
// post_tags
//
// 这张表没有自己的 ID。
//
// 它使用：
//
// post_id + tag_id
//
// 作为联合主键。
type PostTag struct {

	// PostID 文章 ID。
	PostID string `gorm:"column:post_id;type:char(26);primaryKey"`

	// TagID 标签 ID。
	TagID string `gorm:"column:tag_id;type:char(26);primaryKey"`
}

// TableName 指定数据库表名。
func (PostTag) TableName() string {
	return "post_tags"
}
