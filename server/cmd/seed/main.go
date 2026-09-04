package main

import (
	"log"

	"leslie-blog-server/internal/config"
	"leslie-blog-server/internal/database"
	"leslie-blog-server/internal/modules/user/model"
	"leslie-blog-server/internal/pkg/password"
	"leslie-blog-server/internal/pkg/ulid"

	"gorm.io/gorm"
)

func main() {

	// ----------------------------------------
	// 1. 加载项目配置
	// ----------------------------------------
	//
	// 和 API Server 一样，
	// Seed 也需要知道：
	//
	// MySQL 在哪里？
	// 数据库叫什么？
	// 用户名是什么？
	// 密码是什么？
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf(
			"load config failed: %v",
			err,
		)
	}

	// ----------------------------------------
	// 2. 创建数据库连接
	// ----------------------------------------
	db, err := database.NewMySQL(
		cfg.MySQL,
	)

	if err != nil {
		log.Fatalf(
			"connect mysql failed: %v",
			err,
		)
	}

	// ----------------------------------------
	// 3. 获取底层 *sql.DB
	// ----------------------------------------
	//
	// GORM 是对 database/sql 的封装。
	//
	// 程序结束时，
	// 我们应该关闭数据库连接池。
	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf(
			"get sql db failed: %v",
			err,
		)
	}

	defer sqlDB.Close()

	// ----------------------------------------
	// 4. 创建管理员
	// ----------------------------------------
	if err := seedAdmin(db); err != nil {
		log.Fatalf(
			"seed admin failed: %v",
			err,
		)
	}

	log.Println("seed completed")
}

// seedAdmin 创建默认管理员账号。
func seedAdmin(db *gorm.DB) error {

	// ----------------------------------------
	// 1. 定义初始管理员账号
	// ----------------------------------------
	//
	// 开发阶段我们暂时使用：
	//
	// username: admin
	// password: 123456
	//
	// 注意：
	//
	// 真实生产环境不能使用这种弱密码。
	username := "admin"
	plainPassword := "123456"

	// ----------------------------------------
	// 2. 检查管理员是否已经存在
	// ----------------------------------------
	var existingUser model.User

	err := db.
		Where("username = ?", username).
		First(&existingUser).
		Error

	if err == nil {

		// 用户已经存在。
		//
		// Seed 应该尽可能做到幂等。
		//
		// 什么叫幂等？
		//
		// 第一次执行：
		//
		// 创建 admin
		//
		// 第二次执行：
		//
		// 不应该再创建一个 admin。
		log.Printf(
			"user %q already exists, skip",
			username,
		)

		return nil
	}

	// 如果不是“记录不存在”，
	// 那就说明数据库查询真正发生了错误。
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// ----------------------------------------
	// 3. 对管理员密码进行 Hash
	// ----------------------------------------
	passwordHash, err := password.Hash(
		plainPassword,
	)

	if err != nil {
		return err
	}

	// ----------------------------------------
	// 4. 创建 User Model
	// ----------------------------------------
	user := &model.User{
		// 使用 ULID 作为用户 ID。
		ID: ulid.New(),

		// 登录用户名。
		Username: username,

		// 保存 bcrypt Hash，
		// 绝对不能保存 plainPassword。
		PasswordHash: passwordHash,

		// 管理员显示名称。
		DisplayName: "Administrator",

		// 暂时没有头像。
		AvatarURL: "",

		// 1 表示正常状态。
		Status: 1,
	}

	// ----------------------------------------
	// 5. 写入数据库
	// ----------------------------------------
	if err := db.Create(user).Error; err != nil {
		return err
	}

	log.Printf(
		"created admin user: %s",
		user.Username,
	)

	return nil
}
