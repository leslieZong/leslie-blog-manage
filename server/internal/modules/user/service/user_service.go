package service

import (
	"context"
	"errors"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/user/model"
	"leslie-blog-server/internal/modules/user/repository"

	"gorm.io/gorm"
)

type UserService interface {
	GetByID(
		ctx context.Context,
		id string,
	) (*model.User, error)

	GetByUsername(
		ctx context.Context,
		username string,
	) (*model.User, error)

	Create(
		ctx context.Context,
		user *model.User,
	) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(
	repo repository.UserRepository,
) UserService {
	return &userService{
		repo: repo,
	}
}

// GetByID 根据 ID 查询用户。
func (s *userService) GetByID(
	ctx context.Context,
	id string,
) (*model.User, error) {

	// 参数校验。
	if id == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"user id cannot be empty",
		)
	}

	user, err := s.repo.FindByID(ctx, id)

	if err != nil {

		// 数据库没有找到用户。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.New(
				appErrors.ErrUserNotFound,
				404,
				"user not found",
			)
		}

		// 其他数据库错误。
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to find user",
			err,
		)
	}

	return user, nil
}

// GetByUsername 根据用户名查询用户。
func (s *userService) GetByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {

	if username == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"username cannot be empty",
		)
	}

	user, err := s.repo.FindByUsername(
		ctx,
		username,
	)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.New(
				appErrors.ErrUserNotFound,
				404,
				"user not found",
			)
		}

		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to find user",
			err,
		)
	}

	return user, nil
}

// Create 创建用户。
func (s *userService) Create(
	ctx context.Context,
	user *model.User,
) error {

	if user == nil {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"user cannot be nil",
		)
	}

	if user.Username == "" {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"username cannot be empty",
		)
	}

	if user.PasswordHash == "" {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"password hash cannot be empty",
		)
	}

	// 查询用户名是否已经存在。
	existingUser, err := s.repo.FindByUsername(
		ctx,
		user.Username,
	)

	if err == nil && existingUser != nil {
		return appErrors.New(
			appErrors.ErrUsernameExists,
			400,
			"username already exists",
		)
	}

	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {

		return appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to check username",
			err,
		)
	}

	// 真正创建用户。
	if err := s.repo.Create(ctx, user); err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to create user",
			err,
		)
	}

	return nil
}
