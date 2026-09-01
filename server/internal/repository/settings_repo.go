package repository

import (
	"leslie-blog-server/internal/model"

	"gorm.io/gorm"
)

type SettingsRepo struct{ db *gorm.DB }

func NewSettingsRepo(db *gorm.DB) *SettingsRepo { return &SettingsRepo{db: db} }

// Get 取单行设置（id=1），不存在则插入默认行
func (r *SettingsRepo) Get() (*model.Settings, error) {
	var s model.Settings
	err := r.db.First(&s, 1).Error
	if err == gorm.ErrRecordNotFound {
		s = model.Settings{ID: 1, SiteName: "Leslie Blog"}
		if e := r.db.Create(&s).Error; e != nil {
			return nil, e
		}
		return &s, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SettingsRepo) Update(s *model.Settings) error {
	s.ID = 1
	return r.db.Save(s).Error
}
