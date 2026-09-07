package service

import (
	"context"
	"errors"

	appErrors "leslie-blog-server/internal/errors"
	"leslie-blog-server/internal/modules/user/dto"
	"leslie-blog-server/internal/modules/user/model"
	"leslie-blog-server/internal/modules/user/repository"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/password"
	"leslie-blog-server/internal/pkg/ulid"

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

	List(
		ctx context.Context,
		params repository.UserListParams,
	) ([]*model.User, int64, error)

	CreateFromRequest(
		ctx context.Context,
		req *dto.CreateUserRequest,
	) (*model.User, error)

	Update(
		ctx context.Context,
		id string,
		req *dto.UpdateUserRequest,
	) (*model.User, error)

	Delete(
		ctx context.Context,
		operatorUserID string,
		targetUserID string,
	) error
}

type userService struct {
	repo     repository.UserRepository
	enforcer *casbin.Enforcer
}

func NewUserService(
	repo repository.UserRepository,
	enforcer *casbin.Enforcer,
) UserService {
	return &userService{
		repo:     repo,
		enforcer: enforcer,
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
			return nil, appErrors.ErrUserNotFound
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
			return nil, appErrors.ErrUserNotFound
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
		return appErrors.ErrUsernameExists
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

func (s *userService) List(
	ctx context.Context,
	params repository.UserListParams,
) ([]*model.User, int64, error) {

	// Service 层可以对分页参数进行业务约束。
	if params.Offset < 0 {
		params.Offset = 0
	}

	if params.Limit <= 0 {
		params.Limit = 20
	}

	if params.Limit > 100 {
		params.Limit = 100
	}

	return s.repo.List(
		ctx,
		params,
	)
}

func (s *userService) CreateFromRequest(
	ctx context.Context,
	req *dto.CreateUserRequest,
) (*model.User, error) {

	// =========================================================
	// 第一步：防止 nil
	// =========================================================

	if req == nil {
		return nil, errors.New(
			"create user request cannot be nil",
		)
	}

	// =========================================================
	// 第二步：检查用户名是否已经存在
	// =========================================================

	existingUser, err := s.repo.FindByUsername(
		ctx,
		req.Username,
	)

	if err == nil && existingUser != nil {
		return nil, errors.New(
			"username already exists",
		)
	}

	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, err
	}

	// =========================================================
	// 第三步：密码 Hash
	// =========================================================
	//
	// 绝对不能：
	//
	// PasswordHash: req.Password
	//
	// 必须：
	//
	// Plain Password
	//       ↓
	// bcrypt
	//       ↓
	// PasswordHash
	// =========================================================

	passwordHash, err := password.Hash(
		req.Password,
	)

	if err != nil {
		return nil, err
	}

	// =========================================================
	// 第四步：创建 User Model
	// =========================================================

	user := &model.User{
		ID:           ulid.New(),
		Username:     req.Username,
		PasswordHash: passwordHash,
		Email:        req.Email,
		DisplayName:  req.DisplayName,
		AvatarURL:    req.AvatarURL,
		Status:       1,
	}

	// =========================================================
	// 第五步：保存数据库
	// =========================================================

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// =========================================================
	// 第六步：返回创建后的 User
	// =========================================================

	return user, nil
}

func (s *userService) Update(
	ctx context.Context,
	id string,
	req *dto.UpdateUserRequest,
) (*model.User, error) {

	// =========================================================
	// 1. 参数基础检查
	// =========================================================

	if id == "" {
		return nil, errors.New(
			"user id cannot be empty",
		)
	}

	if req == nil {
		return nil, errors.New(
			"update user request cannot be nil",
		)
	}

	// =========================================================
	// 2. 查询用户
	// =========================================================

	user, err := s.repo.FindByID(
		ctx,
		id,
	)

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, appErrors.ErrUserNotFound
		}

		return nil, err
	}

	// =========================================================
	// 3. 构造允许更新的字段
	// =========================================================

	updates := map[string]any{
		"email":        req.Email,
		"display_name": req.DisplayName,
		"avatar_url":   req.AvatarURL,
		"status":       req.Status,
	}

	// =========================================================
	// 4. 执行数据库更新
	// =========================================================

	if err := s.repo.UpdateFields(
		ctx,
		id,
		updates,
	); err != nil {
		return nil, err
	}

	// =========================================================
	// 5. 再查询一次
	//
	// 为什么？
	//
	// 因为数据库中的 updated_at 可能已经发生变化。
	//
	// 我们重新查询可以拿到最新完整数据。
	// =========================================================

	user, err = s.repo.FindByID(
		ctx,
		id,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(
	ctx context.Context,
	operatorUserID string,
	targetUserID string,
) error {

	if operatorUserID == "" {
		return errors.New("operator user id cannot be empty")
	}

	if targetUserID == "" {
		return errors.New("target user id cannot be empty")
	}

	if operatorUserID == targetUserID {
		return appErrors.ErrCannotDeleteSelf
	}

	user, err := s.repo.FindByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.ErrUserNotFound
		}

		return err
	}

	if err := s.repo.Delete(ctx, user); err != nil {
		return err
	}

	if err := s.enforcer.DeleteRolesForUser(targetUserID); err != nil {
		return err
	}

	return nil
}
