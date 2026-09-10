package permission

const (
	// User 权限
	UserRead       = "user:read"
	UserCreate     = "user:create"
	UserUpdate     = "user:update"
	UserDelete     = "user:delete"
	UserAssignRole = "user:assign-role"

	// Role 权限
	RoleRead   = "role:read"
	RoleCreate = "role:create"
	RoleUpdate = "role:update"
	RoleDelete = "role:delete"

	// Post 权限
	PostRead    = "post:read"
	PostCreate  = "post:create"
	PostUpdate  = "post:update"
	PostDelete  = "post:delete"
	PostPublish = "post:publish"

	// Category 权限
	CategoryRead   = "category:read"
	CategoryCreate = "category:create"
	CategoryUpdate = "category:update"
	CategoryDelete = "category:delete"

	// Tag 权限
	TagRead   = "tag:read"
	TagCreate = "tag:create"
	TagUpdate = "tag:update"
	TagDelete = "tag:delete"

	// Project 权限
	ProjectRead   = "project:read"
	ProjectCreate = "project:create"
	ProjectUpdate = "project:update"
	ProjectDelete = "project:delete"
)
