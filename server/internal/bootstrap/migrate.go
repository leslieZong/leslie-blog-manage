package bootstrap

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	sqlmigrations "leslie-blog-server/migrations"
)

// RunMigrations 负责执行尚未应用的数据库迁移。
//
// db 是 GORM 底层使用的数据库连接。
// 迁移工具和业务 Repository 使用同一个数据库。
func RunMigrations(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("run migrations: database connection is nil")
	}

	// 1. 创建 MySQL 迁移驱动。
	//
	// WithInstance 复用项目现有的数据库连接，
	// 不需要再次读取配置或新建数据库连接。
	driver, err := mysqlmigrate.WithInstance(
		db,
		&mysqlmigrate.Config{
			MigrationsTable: "schema_migrations",
		},
	)
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}

	// 2. 从 Go 嵌入的文件系统读取 SQL。
	source, err := iofs.New(sqlmigrations.Files, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	// 3. 创建迁移执行器。
	m, err := migrate.NewWithInstance(
		"iofs",
		source,
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("create migration runner: %w", err)
	}

	// 4. 执行所有尚未应用的 up 迁移。
	err = m.Up()

	// 5. 正常情况下，没有新迁移需要执行时会返回 ErrNoChange。
	// 这不是错误，说明数据库已经是最新版本。
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		_, dirtyErr := m.Close()
		if dirtyErr != nil {
			return fmt.Errorf(
				"migration failed: %w; close migration runner: %v",
				err,
				dirtyErr,
			)
		}

		return fmt.Errorf("apply database migrations: %w", err)
	}

	// 6. 释放迁移工具持有的资源。
	sourceErr, databaseErr := m.Close()
	if sourceErr != nil {
		return fmt.Errorf("close migration source: %w", sourceErr)
	}
	if databaseErr != nil {
		return fmt.Errorf("close migration database driver: %w", databaseErr)
	}

	return nil
}
