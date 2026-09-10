package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	appErrors "leslie-blog-server/internal/errors"
	permissionRepo "leslie-blog-server/internal/modules/permission/repository"
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

	// 查询角色权限
	GetPermissions(
		ctx context.Context,
		roleID string,
	) ([]*dto.RolePermissionResponse, error)

	// 更新角色权限
	UpdatePermissions(
		ctx context.Context,
		roleID string,
		permissionNames []string,
	) error
}

// roleService 是 RoleService 的具体实现。
type roleService struct {
	repo repository.RoleRepository
	// PermissionRepository
	// 负责查询权限定义。
	permissionRepo permissionRepo.PermissionRepository

	enforcer *casbin.Enforcer
}

// NewRoleService 创建 RoleService。
func NewRoleService(
	repo repository.RoleRepository,
	// PermissionRepository
	// 负责查询权限定义。
	permissionRepo permissionRepo.PermissionRepository,
	enforcer *casbin.Enforcer,
) RoleService {

	return &roleService{
		repo:           repo,
		permissionRepo: permissionRepo,
		enforcer:       enforcer,
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

func (s *roleService) GetPermissions(
	ctx context.Context,
	roleID string,
) ([]*dto.RolePermissionResponse, error) {

	// 第一步：
	// 先确认 Role 存在。
	role, err := s.repo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// 第二步：
	// Casbin 使用 role.Name 作为 subject。
	//
	// 例如：
	//
	// role.ID   = 01KABC...
	// role.Name = editor
	//
	// Casbin：
	//
	// p, editor, post, read
	//
	// 所以这里不能传 roleID。
	policies, err := s.enforcer.GetPermissionsForUser(role.Name)
	if err != nil {
		return nil, err
	}

	// 没有权限直接返回空数组。
	if len(policies) == 0 {
		return []*dto.RolePermissionResponse{}, nil
	}

	// 第三步：
	// 将 Casbin Policy 转换成 Permission.Name。
	//
	// editor | post | read
	//
	// ↓
	//
	// post:read
	names := make([]string, 0, len(policies))

	for _, policy := range policies {

		// Casbin 当前模型：
		//
		// policy[0] = subject
		// policy[1] = object
		// policy[2] = action
		//
		// 也就是：
		//
		// [editor, post, read]

		if len(policy) < 3 {
			continue
		}

		name := policy[1] + ":" + policy[2]

		names = append(names, name)
	}

	// 第四步：
	// 根据权限名称批量查询 permissions 表。
	permissions, err := s.permissionRepo.FindByNames(ctx, names)
	if err != nil {
		return nil, err
	}

	// 第五步：
	// 将数据库 Model 转换成 API DTO。
	result := make(
		[]*dto.RolePermissionResponse,
		0,
		len(permissions),
	)

	for _, permission := range permissions {

		result = append(
			result,
			&dto.RolePermissionResponse{
				ID:          permission.ID,
				Name:        permission.Name,
				DisplayName: permission.DisplayName,
				Resource:    permission.Resource,
				Action:      permission.Action,
			},
		)
	}

	return result, nil
}

func (s *roleService) UpdatePermissions(
	ctx context.Context,

	roleID string,
	permissionNames []string,
) error {

	// 1. 查询角色。
	role, err := s.repo.FindByID(ctx, roleID)

	if err != nil {
		return err
	}

	// 2. 清洗权限名称。
	//
	// 防止：
	//
	// "post:read"
	// " post:read "
	// "post:read"
	//
	// 被当成三个不同的权限。
	uniqueNames := make([]string, 0, len(permissionNames))

	seen := make(map[string]struct{})

	for _, name := range permissionNames {

		name = strings.TrimSpace(name)

		if name == "" {
			continue
		}

		if _, exists := seen[name]; exists {
			continue
		}

		seen[name] = struct{}{}
		uniqueNames = append(uniqueNames, name)
	}

	// 3. 批量查询权限定义。
	permissions, err := s.permissionRepo.FindByNames(
		ctx,
		uniqueNames,
	)
	fmt.Println(permissions, uniqueNames)
	if err != nil {
		return err
	}

	// 4. 必须保证：
	//
	// 前端提交几个权限，
	// 数据库就必须找到几个权限。
	//
	// 否则：
	//
	// 前端提交：
	//
	// post:read
	// post:create
	// post:xxx
	//
	// 数据库只有：
	//
	// post:read
	// post:create
	//
	// post:xxx 不存在。
	//
	// 这种情况应该直接报错，
	// 而不是悄悄忽略。
	if len(permissions) != len(uniqueNames) {
		return appErrors.Wrap(
			appErrors.ErrNotFound,
			500,
			"failed to find permission",
			err,
		)
	}

	// 5. 删除角色权限。
	//
	// 这里是“整体替换”策略。
	_, err = s.enforcer.DeleteRolePolicies(
		role.Name,
	)
	if err != nil {
		return err
	}

	// 6. 添加新的权限。
	for _, permission := range permissions {

		err = s.enforcer.AddPolicyIfNotExists(
			role.Name,
			permission.Resource,
			permission.Action,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
