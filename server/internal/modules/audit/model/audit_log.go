package model

import "time"

type AuditLog struct {
	ID string `gorm:"column:id;primaryKey"`

	UserID *string `gorm:"column:user_id"`

	Action string `gorm:"column:action"`

	Resource string `gorm:"column:resource"`

	ResourceID *string `gorm:"column:resource_id"`

	RequestID *string `gorm:"column:request_id"`

	IP *string `gorm:"column:ip"`

	UserAgent *string `gorm:"column:user_agent"`

	Result string `gorm:"column:result"`

	ErrorMessage *string `gorm:"column:error_message"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
