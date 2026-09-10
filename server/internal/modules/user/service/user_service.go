package service

import (
	"context"
	"errors"
	"strings"

	appErrors "leslie-blog-server/internal/errors"
	roleRepository "leslie-blog-server/internal/modules/role/repository"
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

	// GetRoles 查询用户当前拥有的角色。
	GetRoles(
		ctx context.Context,
		userID string,
	) ([]*dto.UserRoleResponse, error)

	// UpdateRoles 更新用户角色。
	//
	// roles 表示用户最终应该拥有的角色。
	UpdateRoles(
		ctx context.Context,
		userID string,
		roles []string,
	) error
}

type userService struct {
	repo repository.UserRepository
	// Role 数据访问。
	//
	// 用于确认角色是否真实存在，
	// 以及获取角色的展示信息。
	roleRepo roleRepository.RoleRepository
	enforcer *casbin.Enforcer
}

func NewUserService(
	repo repository.UserRepository,
	// Role 数据访问。
	//
	// 用于确认角色是否真实存在，
	// 以及获取角色的展示信息。
	roleRepo roleRepository.RoleRepository,
	enforcer *casbin.Enforcer,
) UserService {
	return &userService{
		repo:     repo,
		roleRepo: roleRepo,
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
		return nil, appErrors.ErrUsernameExists
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

// GetRoles 查询用户当前拥有的角色。
func (s *userService) GetRoles(
	ctx context.Context,
	userID string,
) ([]*dto.UserRoleResponse, error) {

	// ==================================================
	// 1. 参数校验
	// ==================================================

	if userID == "" {
		return nil, appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"user id cannot be empty",
		)
	}

	// ==================================================
	// 2. 确认用户存在
	// ==================================================
	//
	// 为什么不能直接查询 Casbin？
	//
	// 因为 Casbin 只知道：
	//
	// userID → role
	//
	// 它不知道 users 表里是否存在这个用户。
	//
	// 所以：
	//
	// URL userID
	//     ↓
	// UserRepository
	//     ↓
	// 用户存在？
	_, err := s.repo.FindByID(
		ctx,
		userID,
	)

	if err != nil {

		// 用户不存在。
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, appErrors.New(
				appErrors.ErrNotFound,
				404,
				"user not found",
			)
		}

		// 数据库等其他错误。
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to find user",
			err,
		)
	}

	// ==================================================
	// 3. 从 Casbin 获取用户角色
	// ==================================================
	//
	// Casbin 中可能存在：
	//
	// g | 01USER001 | editor
	// g | 01USER001 | reviewer
	//
	// GetRolesForUser()
	//
	// 返回：
	//
	// ["editor", "reviewer"]

	roleNames, err := s.enforcer.GetRolesForUser(
		userID,
	)

	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to get user roles",
			err,
		)
	}

	// ==================================================
	// 4. 没有角色
	// ==================================================
	//
	// 注意这里返回：
	//
	// []
	//
	// 而不是：
	//
	// null
	//
	// 对前端来说：
	//
	// []
	//
	// 更容易直接用于：
	//
	// v-for
	// selectedRoles
	// checkbox

	if len(roleNames) == 0 {
		return []*dto.UserRoleResponse{}, nil
	}

	// ==================================================
	// 5. 查询角色详细信息
	// ==================================================
	//
	// Casbin 只有：
	//
	// editor
	// reviewer
	//
	// 但是前端可能需要：
	//
	// ID
	// Name
	// DisplayName
	//
	// 所以：
	//
	// Casbin
	//   ↓
	// Role Name
	//   ↓
	// RoleRepository
	//   ↓
	// Role Model

	roles, err := s.roleRepo.FindByNames(
		ctx,
		roleNames,
	)

	if err != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to find roles",
			err,
		)
	}

	// ==================================================
	// 6. Model → DTO
	// ==================================================

	result := make(
		[]*dto.UserRoleResponse,
		0,
		len(roles),
	)

	for _, role := range roles {

		result = append(
			result,
			&dto.UserRoleResponse{
				ID:          role.ID,
				Name:        role.Name,
				DisplayName: role.DisplayName,
			},
		)
	}

	// ==================================================
	// 7. 返回
	// ==================================================

	return result, nil
}

// UpdateRoles 更新用户角色。
//
// 注意：
//
// roleNames 表示的是“最终角色集合”。
//
// 例如：
//
// 当前：
//
// editor
// reviewer
//
// 请求：
//
// ["viewer"]
//
// 最终：
//
// viewer
//
// 而不是：
//
// editor
// reviewer
// viewer
func (s *userService) UpdateRoles(
	ctx context.Context,
	userID string,
	roleNames []string,
) error {

	// ==================================================
	// 1. 参数检查
	// ==================================================

	if userID == "" {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"user id cannot be empty",
		)
	}

	// ==================================================
	// 2. 确认用户存在
	// ==================================================

	_, err := s.repo.FindByID(
		ctx,
		userID,
	)

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return appErrors.New(
				appErrors.ErrNotFound,
				404,
				"user not found",
			)
		}

		return appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to find user",
			err,
		)
	}

	// ==================================================
	// 3. 去重角色名称
	// ==================================================

	uniqueNames := make(
		[]string,
		0,
		len(roleNames),
	)

	nameSet := make(
		map[string]struct{},
		len(roleNames),
	)

	for _, name := range roleNames {

		// 去除前后空格。
		name = strings.TrimSpace(name)

		// 忽略空角色。
		if name == "" {
			continue
		}

		// 已经存在。
		if _, exists := nameSet[name]; exists {
			continue
		}

		nameSet[name] = struct{}{}

		uniqueNames = append(
			uniqueNames,
			name,
		)
	}

	// ==================================================
	// 4. 查询角色是否存在
	// ==================================================

	roles, err := s.roleRepo.FindByNames(
		ctx,
		uniqueNames,
	)

	if err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to find roles",
			err,
		)
	}

	// ==================================================
	// 5. 检查角色是否全部存在
	// ==================================================

	if len(roles) != len(uniqueNames) {
		return appErrors.New(
			appErrors.ErrInvalidParams,
			400,
			"one or more roles do not exist",
		)
	}

	// ==================================================
	// 6. 查询当前角色
	// ==================================================

	currentRoles, err := s.enforcer.GetRolesForUser(
		userID,
	)

	if err != nil {
		return appErrors.Wrap(
			appErrors.ErrInternalServer,
			500,
			"failed to get current user roles",
			err,
		)
	}

	// ==================================================
	// 7. 删除旧角色
	// ==================================================

	for _, roleName := range currentRoles {

		if err := s.enforcer.DeleteRoleForUser(
			userID,
			roleName,
		); err != nil {

			return appErrors.Wrap(
				appErrors.ErrInternalServer,
				500,
				"failed to remove current user role",
				err,
			)
		}
	}

	// ==================================================
	// 8. 添加新角色
	// ==================================================

	for _, roleName := range uniqueNames {

		if err := s.enforcer.AddRoleForUser(
			userID,
			roleName,
		); err != nil {

			return appErrors.Wrap(
				appErrors.ErrInternalServer,
				500,
				"failed to assign user role",
				err,
			)
		}
	}

	// ==================================================
	// 9. 完成
	// ==================================================

	return nil
}
