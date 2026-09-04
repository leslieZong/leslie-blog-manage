package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Enforcer 是整个项目使用的 Casbin 权限执行器。
//
// Casbin 的核心能力都通过 Enforcer 完成。
//
// 例如：
//
// Enforce()
// AddPolicy()
// AddRoleForUser()
// DeletePolicy()
// ...
type Enforcer struct {
	// engine 是真正的 Casbin Enforcer。
	engine *casbin.Enforcer
}

// New 创建 Casbin Enforcer。
//
// db 是 MySQL GORM 连接。
// modelPath 是 Casbin Model 文件路径。
func New(
	db *gorm.DB,
	modelPath string,
) (*Enforcer, error) {

	// ==================================================
	// 1. 创建 GORM Adapter
	// ==================================================
	//
	// Adapter 的作用：
	//
	// Casbin
	//    ↓
	// Adapter
	//    ↓
	// GORM
	//    ↓
	// MySQL
	//
	// Casbin 不需要直接操作 MySQL。
	// 它通过 Adapter 保存 Policy。
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf(
			"create casbin adapter failed: %w",
			err,
		)
	}

	// ==================================================
	// 2. 根据 Model 创建 Casbin Enforcer
	// ==================================================

	engine, err := casbin.NewEnforcer(
		modelPath,
		adapter,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create casbin enforcer failed: %w",
			err,
		)
	}

	// ==================================================
	// 3. 从数据库加载 Policy
	// ==================================================

	if err := engine.LoadPolicy(); err != nil {
		return nil, fmt.Errorf(
			"load casbin policy failed: %w",
			err,
		)
	}

	// ==================================================
	// 4. 返回我们的 Enforcer
	// ==================================================

	return &Enforcer{
		engine: engine,
	}, nil
}

// Enforce 检查当前用户是否拥有指定权限。
//
// subject：用户 / 角色
// object：资源
// action：动作
//
// 返回：
//
// true
//
//	允许
//
// false
//
//	拒绝
func (e *Enforcer) Enforce(
	subject string,
	object string,
	action string,
) (bool, error) {

	return e.engine.Enforce(
		subject,
		object,
		action,
	)
}

// AddRoleForUser 给用户分配角色。
//
// 例如：
//
// AddRoleForUser("01K...", "admin")
//
// 表示：
//
// 用户 01K...
//
//	↓
//
// admin
func (e *Enforcer) AddRoleForUser(
	userID string,
	role string,
) error {

	_, err := e.engine.AddRoleForUser(
		userID,
		role,
	)

	return err
}

// DeleteRoleForUser 删除用户角色。
func (e *Enforcer) DeleteRoleForUser(
	userID string,
	role string,
) error {

	_, err := e.engine.DeleteRoleForUser(
		userID,
		role,
	)

	return err
}

// GetRolesForUser 查询用户角色。
func (e *Enforcer) GetRolesForUser(
	userID string,
) ([]string, error) {

	return e.engine.GetRolesForUser(userID)
}

// AddPolicy 添加权限策略。
//
// 例如：
//
// AddPolicy("admin", "post", "delete")
//
// 表示 admin 可以删除文章。
func (e *Enforcer) AddPolicy(
	role string,
	object string,
	action string,
) error {

	_, err := e.engine.AddPolicy(
		role,
		object,
		action,
	)

	return err
}

// DeletePolicy 删除权限策略。
func (e *Enforcer) DeletePolicy(
	role string,
	object string,
	action string,
) error {

	_, err := e.engine.RemovePolicy(
		role,
		object,
		action,
	)

	return err
}

// AddPolicyIfNotExists 添加权限策略。
//
// 如果权限已经存在，则什么都不做。
func (e *Enforcer) AddPolicyIfNotExists(
	role string,
	object string,
	action string,
) error {

	exists, err := e.engine.HasPolicy(
		role,
		object,
		action,
	)

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = e.engine.AddPolicy(
		role,
		object,
		action,
	)

	return err
}

// AddRoleForUserIfNotExists 给用户分配角色。
//
// 如果用户已经拥有该角色，则什么都不做。
func (e *Enforcer) AddRoleForUserIfNotExists(
	userID string,
	role string,
) error {

	hasRole, err := e.engine.HasRoleForUser(
		userID,
		role,
	)

	if err != nil {
		return err
	}

	if hasRole {
		return nil
	}

	_, err = e.engine.AddRoleForUser(
		userID,
		role,
	)

	return err
}
