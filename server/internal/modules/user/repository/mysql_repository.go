package repository

import (
	"context"
	"errors"

	"leslie-blog-server/internal/modules/user/model"

	"gorm.io/gorm"
)

// mysqlUserRepository 是 UserRepository 的 MySQL 实现。
//
// 注意：
//
// 这个 struct 没有导出，
// 因为我们不希望模块外部直接依赖具体实现。
//
// 外部只需要知道：
//
// # UserRepository
//
// 而不需要知道：
//
// mysqlUserRepository
type mysqlUserRepository struct {

	// db 是 GORM 数据库对象。
	//
	// 它来自：
	//
	// internal/database
	//
	// 最终连接到 MySQL。
	db *gorm.DB
}

// NewUserRepository 创建 UserRepository。
//
// 这里就是我们前面讲过的：
//
// # Dependency Injection
//
// “依赖注入”。
//
// 我们不是在 Repository 内部自己创建数据库：
//
// ❌ gorm.Open(...)
//
// 而是从外面把 db 传进来：
//
// ✅ NewUserRepository(db)
func NewUserRepository(db *gorm.DB) UserRepository {

	return &mysqlUserRepository{
		db: db,
	}
}

// FindByID 根据用户 ID 查询用户。
func (r *mysqlUserRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.User, error) {

	// 创建 User 对象，用来接收数据库查询结果。
	var user model.User

	// 使用 GORM 查询。
	//
	// WithContext(ctx)
	//
	// 把当前请求的 Context 传给数据库操作。
	//
	// 如果请求被取消，
	// 数据库操作也可以感知这个取消。
	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&user).
		Error

	// 如果查询发生错误。
	if err != nil {

		// 如果是“数据不存在”，
		// 直接把 gorm.ErrRecordNotFound 返回出去。
		//
		// Service 层可以根据这个错误判断：
		//
		// 用户不存在。
		return nil, err
	}

	// 查询成功。
	return &user, nil
}

// FindByUsername 根据用户名查询用户。
func (r *mysqlUserRepository) FindByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {

	// 创建 User 对象接收查询结果。
	var user model.User

	// 根据 username 查询。
	err := r.db.
		WithContext(ctx).
		Where("username = ?", username).
		First(&user).
		Error

	// 返回查询结果。
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create 创建用户。
func (r *mysqlUserRepository) Create(
	ctx context.Context,
	user *model.User,
) error {

	// GORM Create 会生成 INSERT SQL。
	//
	// 类似：
	//
	// INSERT INTO users (...)
	// VALUES (...)
	//
	// 具体 SQL 由 GORM 生成。
	return r.db.
		WithContext(ctx).
		Create(user).
		Error
}

// Update 更新用户。
func (r *mysqlUserRepository) Update(
	ctx context.Context,
	user *model.User,
) error {

	// Save 会根据 User 的主键进行保存。
	//
	// 当前 User.ID 就是数据库中的主键。
	return r.db.
		WithContext(ctx).
		Save(user).
		Error
}

// IsNotFound 判断 Repository 返回的错误
// 是否代表“数据不存在”。
//
// 这个函数暂时属于 Repository 层的辅助能力。
func IsNotFound(err error) bool {

	return errors.Is(
		err,
		gorm.ErrRecordNotFound,
	)
}
