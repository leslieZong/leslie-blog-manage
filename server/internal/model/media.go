package model

// Media 媒体文件
type Media struct {
	BaseModel
	Name     string `gorm:"size:255;not null" json:"name"`
	URL      string `gorm:"column:url;size:512;not null" json:"url"`
	Type     string `gorm:"size:32;index" json:"type"` // image / video / file
	Size     int64  `gorm:"default:0" json:"size"`
	MimeType string `gorm:"size:128" json:"mimeType"`
}

func (Media) TableName() string { return "media" }
