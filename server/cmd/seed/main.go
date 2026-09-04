package main

import (
	"log"

	"leslie-blog-server/internal/config"
	"leslie-blog-server/internal/database"
	"leslie-blog-server/internal/modules/user/model"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/password"
	"leslie-blog-server/internal/pkg/ulid"

	"gorm.io/gorm"
)

func main() {

	// 1. 加载配置。
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf(
			"load config failed: %v",
			err,
		)
	}

	// 2. 连接 MySQL。
	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf(
			"connect mysql failed: %v",
			err,
		)
	}

	// 3. 初始化 Casbin。
	enforcer, err := casbin.New(
		db,
		"./configs/casbin_model.conf",
	)

	if err != nil {
		log.Fatalf(
			"create casbin failed: %v",
			err,
		)
	}

	// 4. 创建 Admin 用户。
	adminUser, err := seedAdmin(db)
	if err != nil {
		log.Fatalf(
			"seed admin failed: %v",
			err,
		)
	}

	// 5. 初始化 Admin 权限。
	if err := seedCasbin(
		enforcer,
		adminUser.ID,
	); err != nil {
		log.Fatalf(
			"seed casbin failed: %v",
			err,
		)
	}

	log.Println("seed completed")
}

// seedAdmin 创建或获取 admin 用户。
func seedAdmin(
	db *gorm.DB,
) (*model.User, error) {

	username := "admin"
	plainPassword := "123456"

	var existingUser model.User

	err := db.
		Where("username = ?", username).
		First(&existingUser).
		Error

	if err == nil {
		log.Printf(
			"user %q already exists, skip",
			username,
		)

		return &existingUser, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	passwordHash, err := password.Hash(
		plainPassword,
	)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           ulid.New(),
		Username:     username,
		PasswordHash: passwordHash,
		DisplayName:  "Administrator",
		AvatarURL:    "",
		Status:       1,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	log.Printf(
		"created admin user: %s",
		user.Username,
	)

	return user, nil
}
