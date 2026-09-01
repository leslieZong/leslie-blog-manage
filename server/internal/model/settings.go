package model

// Settings 站点设置（单行，id=1）
type Settings struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	SiteName      string `gorm:"column:site_name" json:"siteName"`
	Logo          string `gorm:"size:512" json:"logo"`
	Description   string `gorm:"size:500" json:"description"`
	Keywords      string `gorm:"size:255" json:"keywords"`
	Author        string `gorm:"size:64" json:"author"`
	ICP           string `gorm:"size:128" json:"icp"`
	SocialGithub  string `gorm:"column:social_github;size:255" json:"socialGithub,omitempty"`
	SocialEmail   string `gorm:"column:social_email;size:255" json:"socialEmail,omitempty"`
	SocialTwitter string `gorm:"column:social_twitter;size:255" json:"socialTwitter,omitempty"`
}

func (Settings) TableName() string { return "settings" }
