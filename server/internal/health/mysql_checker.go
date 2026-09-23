package health

import (
	"context"

	"gorm.io/gorm"
)

// MySQLChecker 用于检查 MySQL 是否可以正常访问。
type MySQLChecker struct {
	db *gorm.DB
}

// NewMySQLChecker 创建 MySQL Checker。
func NewMySQLChecker(
	db *gorm.DB,
) *MySQLChecker {

	return &MySQLChecker{
		db: db,
	}
}

// Check 执行一次非常轻量的数据库检查。
func (c *MySQLChecker) Check(
	ctx context.Context,
) error {

	// GORM 底层使用 database/sql。
	//
	// DB() 可以获取原生 *sql.DB，
	// 然后使用 PingContext 检查数据库连接。
	sqlDB, err := c.db.DB()

	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}
