package audit

// Event 表示一次业务审计事件。
//
// 它描述的是：
//
// 谁不需要业务层自己填写
// 请求信息不需要业务层自己填写
//
// 业务层只需要告诉 Audit：
//
// 1. 做了什么
// 2. 操作什么资源
// 3. 操作哪个资源
type Event struct {

	// Action 表示具体业务动作。
	//
	// 例如：
	//
	// post.create
	// post.update
	// post.publish
	// post.delete
	Action string

	// Resource 表示资源类型。
	//
	// 例如：
	//
	// post
	// user
	// role
	// project
	Resource string

	// ResourceID 表示具体操作的资源 ID。
	//
	// 例如：
	//
	// post.ID
	// user.ID
	// role.ID
	ResourceID string
}
