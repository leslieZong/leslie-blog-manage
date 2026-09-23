package router

import (
	"leslie-blog-server/internal/health"
	"leslie-blog-server/internal/modules/auth"
	authHandler "leslie-blog-server/internal/modules/auth/handler"
	"leslie-blog-server/internal/modules/category"
	categoryHandler "leslie-blog-server/internal/modules/category/handler"
	"leslie-blog-server/internal/modules/home"
	homeHandler "leslie-blog-server/internal/modules/home/handler"
	"leslie-blog-server/internal/modules/post"
	postHandler "leslie-blog-server/internal/modules/post/handler"
	"leslie-blog-server/internal/modules/project"
	projectHandler "leslie-blog-server/internal/modules/project/handler"
	"leslie-blog-server/internal/modules/role"
	roleHandler "leslie-blog-server/internal/modules/role/handler"
	"leslie-blog-server/internal/modules/tag"
	tagHandler "leslie-blog-server/internal/modules/tag/handler"
	"leslie-blog-server/internal/modules/techstack"
	techstackHandler "leslie-blog-server/internal/modules/techstack/handler"
	"leslie-blog-server/internal/modules/user"
	userHandler "leslie-blog-server/internal/modules/user/handler"
	"leslie-blog-server/internal/pkg/casbin"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine

	userHandler *userHandler.UserHandler

	authHandler           *authHandler.AuthHandler
	roleHandler           *roleHandler.RoleHandler
	postHandler           *postHandler.PostHandler
	publicPostHandler     *postHandler.PublicPostHandler
	categoryHandler       *categoryHandler.CategoryHandler
	publicCategoryHandler *categoryHandler.PublicCategoryHandler
	tagHandler            *tagHandler.TagHandler
	projectHandler        *projectHandler.ProjectHandler
	projectPublicHandler  *projectHandler.ProjectPublicHandler
	techstackHandler      *techstackHandler.TechStackHandler
	homeHandler           *homeHandler.HomeHandler
	healthHandler         *health.Handler
	readyHandler          *health.ReadyHandler

	// jwtMiddleware 是 JWT 认证中间件。
	//
	// Server 创建好 Middleware 后，
	// 注入到 Router。
	jwtMiddleware gin.HandlerFunc
	enforcer      *casbin.Enforcer
}

// New 创建 Router。
func New(
	engine *gin.Engine,
	userHandler *userHandler.UserHandler,
	authHandler *authHandler.AuthHandler,
	roleHandler *roleHandler.RoleHandler,
	postHandler *postHandler.PostHandler,
	publicPostHandler *postHandler.PublicPostHandler,
	categoryHandler *categoryHandler.CategoryHandler,
	publicCategoryHandler *categoryHandler.PublicCategoryHandler,
	tagHandler *tagHandler.TagHandler,
	projectHandler *projectHandler.ProjectHandler,
	projectPublicHandler *projectHandler.ProjectPublicHandler,
	techstackHandler *techstackHandler.TechStackHandler,
	homeHandler *homeHandler.HomeHandler,
	healthHandler *health.Handler,
	readyHandler *health.ReadyHandler,
	jwtMiddleware gin.HandlerFunc,
	enforcer *casbin.Enforcer,
) *Router {

	return &Router{
		engine:                engine,
		userHandler:           userHandler,
		authHandler:           authHandler,
		roleHandler:           roleHandler,
		postHandler:           postHandler,
		publicPostHandler:     publicPostHandler,
		categoryHandler:       categoryHandler,
		publicCategoryHandler: publicCategoryHandler,
		tagHandler:            tagHandler,
		projectHandler:        projectHandler,
		projectPublicHandler:  projectPublicHandler,
		techstackHandler:      techstackHandler,
		homeHandler:           homeHandler,
		healthHandler:         healthHandler,
		readyHandler:          readyHandler,
		jwtMiddleware:         jwtMiddleware,
		enforcer:              enforcer,
	}
}

// Register 注册整个项目的路由。
func (r *Router) Register() {

	// ==================================================
	// 健康检查
	// ==================================================
	r.engine.GET(
		"/health",
		r.healthHandler.Health,
	)

	r.engine.GET(
		"/ready",
		r.readyHandler.Ready,
	)
	// ==================================================
	// Public API
	// ==================================================

	v1 := r.engine.Group("/api/v1")
	post.RegisterPublicRoutes(
		v1,
		r.publicPostHandler,
	)
	category.RegisterPublicRoutes(
		v1,
		r.publicCategoryHandler,
	)
	project.RegisterPublicRoutes(
		v1,
		r.projectPublicHandler,
	)
	techstack.RegisterPublicRoutes(
		v1,
		r.techstackHandler,
	)
	home.RegisterPublicRoutes(
		v1,
		r.homeHandler,
	)

	// 当前还没有公共 API。
	//
	// 所以这里先保留。
	_ = v1

	// ==================================================
	// Admin API
	// ==================================================

	admin := r.engine.Group("/api/admin/v1")

	// --------------------------------------------------
	// 登录接口
	// --------------------------------------------------
	//
	// 登录之前没有 JWT，
	// 所以这里不能使用 JWT Middleware。
	auth.RegisterRoutes(
		admin,
		r.authHandler,
		r.jwtMiddleware,
	)

	// --------------------------------------------------
	// 需要登录的接口
	// --------------------------------------------------
	//
	// 创建一个新的 Router Group。
	protected := admin.Group("")

	// 给 protected Group 添加 JWT Middleware。
	protected.Use(r.jwtMiddleware)

	// 所有注册到 protected 的接口，
	// 都必须先通过 JWT 验证。
	user.RegisterRoutes(
		protected,
		r.userHandler,
		r.enforcer,
	)

	role.RegisterRoutes(
		protected,
		r.roleHandler,
		r.enforcer,
	)

	post.RegisterRoutes(
		protected,
		r.postHandler,
		r.enforcer,
	)

	category.RegisterRoutes(
		protected,
		r.categoryHandler,
		r.enforcer,
	)
	tag.RegisterRoutes(
		protected,
		r.tagHandler,
		r.enforcer,
	)
	project.RegisterRoutes(
		protected,
		r.projectHandler,
		r.enforcer,
	)
	techstack.RegisterRoutes(
		protected,
		r.techstackHandler,
		r.enforcer,
	)

}
