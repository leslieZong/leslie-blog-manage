package database

import (
	"fmt"
	"time"

	"leslie-blog-server/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMySQL 创建 MySQL 数据库连接。
//
// 参数：
//
//	cfg：MySQL 配置
//
// 返回：
//
//	*gorm.DB：GORM 数据库对象
//	error：数据库连接过程中可能发生的错误
//
// 为什么返回 error？
//
// 因为数据库连接是一个非常容易失败的操作：
//
//   - MySQL 没启动
//   - 用户名错误
//   - 密码错误
//   - 数据库不存在
//   - 端口错误
//   - 网络连接失败
//
// Go 的设计思想是：
// “错误不是异常，而是函数正常返回值的一部分。”
func NewMySQL(cfg config.MySQLConfig) (*gorm.DB, error) {

	// --------------------------------------------------
	// 第一步：构造 MySQL DSN
	// --------------------------------------------------
	//
	// DSN：
	// Data Source Name
	//
	// 可以理解成：
	// “告诉程序应该连接哪个数据库”。
	//
	// 最终类似：
	//
	// root:password@tcp(127.0.0.1:3306)/leslie_blog
	//     ?charset=utf8mb4&parseTime=True&loc=Local
	//
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)

	// --------------------------------------------------
	// 第二步：创建 GORM DB
	// --------------------------------------------------
	//
	// mysql.Open(dsn)
	//
	// 告诉 GORM：
	// “我要使用 MySQL，并且这是数据库连接信息。”
	//
	// gorm.Open() 返回：
	//
	//   *gorm.DB
	//   error
	//
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)

	// 如果 GORM 初始化失败，立即返回错误。
	if err != nil {
		return nil, fmt.Errorf("connect mysql failed: %w", err)
	}

	// --------------------------------------------------
	// 第三步：获取底层 database/sql 连接池
	// --------------------------------------------------
	//
	// GORM 底层实际上还是使用 Go 标准库：
	//
	// database/sql
	//
	// db.DB() 可以拿到底层的 *sql.DB。
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db failed: %w", err)
	}

	// --------------------------------------------------
	// 第四步：配置数据库连接池
	// --------------------------------------------------

	// 设置最大打开连接数。
	//
	// 例如：
	// 20
	//
	// 表示这个程序最多同时打开 20 个数据库连接。
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// 设置最大空闲连接数。
	//
	// 当连接暂时没有使用时，
	// 可以保留在连接池中等待下一次请求。
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// --------------------------------------------------
	// 第五步：解析连接生命周期
	// --------------------------------------------------
	//
	// 配置文件中是：
	//
	// conn_max_lifetime: "1h"
	//
	// 但是 Go 的连接池需要：
	//
	// time.Duration
	//
	// 所以需要转换。
	connMaxLifetime, err := time.ParseDuration(
		cfg.ConnMaxLifetime,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"parse mysql conn max lifetime failed: %w",
			err,
		)
	}

	// 设置连接最大生命周期。
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	// --------------------------------------------------
	// 第六步：真正测试数据库连接
	// --------------------------------------------------
	//
	// 注意：
	//
	// gorm.Open() 成功，
	// 不代表 MySQL 一定真的可以访问。
	//
	// 所以这里调用 Ping()。
	//
	// Ping 的意思：
	// “请真正尝试访问一下数据库。”
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql failed: %w", err)
	}

	// 到这里说明：
	//
	// 1. GORM 初始化成功
	// 2. 底层连接池获取成功
	// 3. 连接池配置成功
	// 4. MySQL Ping 成功
	//
	// 所以可以把 db 返回给上层。
	return db, nil
}
