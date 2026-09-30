package audit

const (
	// User
	ActionUserCreate = "user.create"
	ActionUserUpdate = "user.update"
	ActionUserDelete = "user.delete"

	// Role
	ActionRoleCreate = "role.create"
	ActionRoleUpdate = "role.update"
	ActionRoleDelete = "role.delete"

	// Post
	ActionPostCreate  = "post.create"
	ActionPostUpdate  = "post.update"
	ActionPostDelete  = "post.delete"
	ActionPostPublish = "post.publish"
	ActionPostArchive = "post.archive"

	// Category
	ActionCategoryCreate = "category.create"
	ActionCategoryUpdate = "category.update"
	ActionCategoryDelete = "category.delete"

	// Tag
	ActionTagCreate = "tag.create"
	ActionTagUpdate = "tag.update"
	ActionTagDelete = "tag.delete"

	// Project
	ActionProjectCreate = "project.create"
	ActionProjectUpdate = "project.update"
	ActionProjectDelete = "project.delete"

	// TechStack
	ActionTechStackCreate = "techstack.create"
	ActionTechStackUpdate = "techstack.update"
	ActionTechStackDelete = "techstack.delete"
)
