package database

import (
	"context"

	"gorm.io/gorm"
)

// TransactionManager 负责执行数据库事务。
//
// 它把 GORM 的事务细节隐藏起来。
//
// Service 只需要知道：
//
// “我要执行一个事务”
//
// 而不需要知道：
//
// Begin()
// Commit()
// Rollback()
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager 创建事务管理器。
func NewTransactionManager(
	db *gorm.DB,
) *TransactionManager {

	return &TransactionManager{
		db: db,
	}
}

// WithTransaction 在事务中执行 fn。
//
// 业务流程：
//
// BEGIN
//
//	↓
//
// fn()
//
//	↓
//
// 成功？
// ├── 是 → COMMIT
// └── 否 → ROLLBACK
func (m *TransactionManager) WithTransaction(
	ctx context.Context,
	fn func(tx *gorm.DB) error,
) error {
	return m.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {

			return fn(tx)
		})
}
