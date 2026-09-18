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

	// -----------------------------------------------------
	// 开启事务
	// -----------------------------------------------------

	tx := m.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// -----------------------------------------------------
	// 执行业务逻辑
	// -----------------------------------------------------

	if err := fn(tx); err != nil {

		// -------------------------------------------------
		// 业务失败 → 回滚
		// -------------------------------------------------

		_ = tx.Rollback()

		return err
	}

	// -----------------------------------------------------
	// 所有操作成功 → 提交
	// -----------------------------------------------------

	return tx.Commit().Error
}
