package service

import (
	"context"
	"errors"

	"leslie-blog-server/internal/modules/user/model"
	"leslie-blog-server/internal/modules/user/repository"

	"gorm.io/gorm"
)

// UserService 定义 User 业务能力。
//
// Repository 解决：
//
// “数据怎么存取？”
//
// Service 解决：
//
// “业务应该怎么运行？”
type UserService interface {

	// 根据 ID 获取用户。
	GetByID(
		ctx context.Context,
		id string,
	) (*model.User, error)

	// 根据用户名获取用户。
	GetByUsername(
		ctx context.Context,
		username string,
	) (*model.User, error)

	// 创建用户。
	Create(
		ctx context.Context,
		user *model.User,
	) error
}

// userService 是 UserService 的具体实现。
type userService struct {

	// repo 是 UserRepository。
	//
	// 注意：
	//
	// Service 依赖的是：
	//
	// UserRepository interface
	//
	// 而不是：
	//
	// mysqlUserRepository
	//
	// 这就是面向接口编程。
	repo repository.UserRepository
}

// NewUserService 创建 UserService。
//
// 这里同样使用依赖注入。
func NewUserService(
	repo repository.UserRepository,
) UserService {

	return &userService{
		repo: repo,
	}
}

// GetByID 根据 ID 获取用户。
func (s *userService) GetByID(
	ctx context.Context,
	id string,
) (*model.User, error) {

	// 这里可以增加业务规则。
	//
	// 例如：
	//
	// id 不能为空。
	if id == "" {
		return nil, errors.New("user id cannot be empty")
	}

	// Repository 只负责真正查询数据库。
	return s.repo.FindByID(ctx, id)
}

// GetByUsername 根据用户名获取用户。
func (s *userService) GetByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {

	// 业务参数检查。
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	// 调用 Repository 查询数据库。
	return s.repo.FindByUsername(ctx, username)
}

// Create 创建用户。
func (s *userService) Create(
	ctx context.Context,
	user *model.User,
) error {

	// 基础业务校验。
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if user.Username == "" {
		return errors.New("username cannot be empty")
	}

	if user.PasswordHash == "" {
		return errors.New("password hash cannot be empty")
	}

	// 检查用户名是否已经存在。
	existingUser, err := s.repo.FindByUsername(
		ctx,
		user.Username,
	)

	if err == nil && existingUser != nil {
		return errors.New("username already exists")
	}

	// 如果不是“数据不存在”，
	// 而是数据库真的发生错误，
	// 就继续返回错误。
	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {

		return err
	}

	// 真正创建用户。
	return s.repo.Create(ctx, user)
}
