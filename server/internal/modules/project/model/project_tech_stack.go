package model

// ProjectTechStack 表示 Project 与 TechStack 的关系。
//
// 它不仅表示：
//
// Project 使用了哪个 TechStack。
//
// 还保存：
//
// sort
//
// 因此它本身也是一个有业务意义的数据模型。
type ProjectTechStack struct {
	ProjectID string `gorm:"column:project_id;type:char(26);primaryKey"`

	TechStackID string `gorm:"column:tech_stack_id;type:char(26);primaryKey"`

	Sort int `gorm:"column:sort;not null;default:0"`
}

func (ProjectTechStack) TableName() string {
	return "project_tech_stacks"
}
