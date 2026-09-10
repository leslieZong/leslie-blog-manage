package service

import (
	"context"
	"errors"

	appErrors "leslie-blog-server/internal/errors"
	constant "leslie-blog-server/internal/modules/role/const"
	"leslie-blog-server/internal/modules/role/dto"
	"leslie-blog-server/internal/modules/role/model"
	"leslie-blog-server/internal/modules/role/repository"
	"leslie-blog-server/internal/pkg/casbin"
	"leslie-blog-server/internal/pkg/ulid"

	"gorm.io/gorm"
)

// RoleService 定义角色业务能力。
type RoleService interface {

	// 查询单个角色。
	GetByID(
		ctx context.Context,
		id string,
	) (*model.Role, error)

	// 查询角色列表。
	List(
		ctx context.Context,
		params repository.RoleListParams,
	) ([]*model.Role, int64, error)

	// 创建角色。
	Create(
		ctx context.Context,
		req *dto.CreateRoleRequest,
	) (*model.Role, error)

	// 更新角色。
	Update(
		ctx context.Context,
		id string,
		req *dto.UpdateRoleRequest,
	) (*model.Role, error)

	// 删除角色。
	Delete(
		ctx context.Context,
		id string,
	) error
}

// roleService 是 RoleService 的具体实现。
type roleService struct {
	repo     repository.RoleRepository
	enforcer *casbin.Enforcer
}

// NewRoleService 创建 RoleService。
func NewRoleService(
	repo repository.RoleRepository,
	enforcer *casbin.Enforcer,
) RoleService {

	return &roleService{
		repo:     repo,
		enforcer: enforcer,
	}
}

func (s *roleService) GetByID(
	ctx context.Context,
	id string,
) (*model.Role, error) {

	if id == "" {
		return nil, errors.New(
			"role id cannot be empty",
		)
	}

	role, err := s.repo.FindByID(ctx, id)

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, appErrors.Wrap(
				appErrors.ErrNotFound,
				500,
				"failed to find role",
				err,
			)
		}

		return nil, err
	}

	return role, nil
}

func (s *roleService) Create(
	ctx context.Context,
	req *dto.CreateRoleRequest,
) (*model.Role, error) {

	if req == nil {
		return nil, errors.New(
			"create role request cannot be nil",
		)
	}

	// 先检查角色名称是否已经存在。
	existingRole, err := s.repo.FindByName(
		ctx,
		req.Name,
	)

	if err == nil && existingRole != nil {
		return nil, appErrors.Wrap(
			appErrors.ErrConflict,
			40901,
			"role name already exists",
			err,
		)
	}

	// FindByName 没找到数据是正常情况。
	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, err
	}

	// 创建角色 Model。
	role := &model.Role{
		ID:          ulid.New(),
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Status:      1,
	}

	// 保存到数据库。
	if err := s.repo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleService) Update(
	ctx context.Context,
	id string,
	req *dto.UpdateRoleRequest,
) (*model.Role, error) {

	if id == "" {
		return nil, errors.New(
			"role id cannot be empty",
		)
	}

	if req == nil {
		return nil, errors.New(
			"update role request cannot be nil",
		)
	}

	// 确认角色存在。
	_, err := s.repo.FindByID(ctx, id)

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, appErrors.Wrap(
				appErrors.ErrNotFound,
				500,
				"failed to find role",
				err,
			)
		}

		return nil, err
	}

	updates := map[string]any{
		"display_name": req.DisplayName,
		"description":  req.Description,
		"status":       req.Status,
	}

	if err := s.repo.UpdateFields(
		ctx,
		id,
		updates,
	); err != nil {
		return nil, err
	}

	role, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleService) Delete(
	ctx context.Context,
	id string,
) error {
	// 1. 查询角色
	// 2. 判断角色是否存在
	// 3. 判断是否系统角色
	// 4. 判断是否正在被用户使用
	// 5. 删除角色
	// 6. 清理 Casbin

	if id == "" {
		return errors.New(
			"role id cannot be empty",
		)
	}

	role, err := s.repo.FindByID(ctx, id)

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return appErrors.Wrap(
				appErrors.ErrNotFound,
				500,
				"failed to find role",
				err,
			)

		}

		return err
	}

	// admin 是系统内置超级管理员角色。
	//
	// 第一阶段不允许删除。
	if constant.IsSystemRole(role.Name) {
		return errors.New(
			"cannot delete system role",
		)
	}

	// 判断是否正在被用户使用。
	users, err := s.enforcer.GetUsersForRole(role.Name)
	if err != nil {
		return err
	}

	if len(users) > 0 {
		return errors.New(
			"role is in use",
		)
	}
	// 删除角色对应的 Casbin 权限
	if _, err := s.enforcer.DeleteRolePolicies(role.Name); err != nil {
		return err
	}
	// 删除角色
	if err := s.repo.Delete(ctx, role); err != nil {
		return err
	}

	return nil
}

func (s *roleService) List(
	ctx context.Context,
	params repository.RoleListParams,
) ([]*model.Role, int64, error) {

	if params.Offset < 0 {
		params.Offset = 0
	}

	if params.Limit <= 0 {
		params.Limit = 20
	}

	return s.repo.List(ctx, params)
}
