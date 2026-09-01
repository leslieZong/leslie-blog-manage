package model

// TechStack 技术栈
type TechStack struct {
	BaseModel
	Name        string `gorm:"size:64;not null" json:"name"`
	Icon        string `gorm:"size:512" json:"icon"`
	Category    string `gorm:"size:64" json:"category"`
	Level       int    `gorm:"default:0" json:"level"` // 0-100
	Description string `gorm:"size:255" json:"description"`
	SortOrder   int    `gorm:"default:0" json:"sortOrder"`
}

func (TechStack) TableName() string { return "tech_stack" }
