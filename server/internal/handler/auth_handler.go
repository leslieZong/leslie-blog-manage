package handler

import (
	"leslie-blog-server/internal/model"
	"leslie-blog-server/internal/service"
	"leslie-blog-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{ svc *service.AuthService }

func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

// Login POST /api/v1/admin/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	resp, err := h.svc.Login(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, resp)
}

// Logout POST /api/v1/admin/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// 无状态 JWT：前端删除 token 即可
	response.OKMsg(c, "已退出")
}

// Me GET /api/v1/admin/auth/user
func (h *AuthHandler) Me(c *gin.Context) {
	uid := currentUserID(c)
	info, err := h.svc.CurrentUser(uid)
	if err != nil {
		response.Fail(c, response.CodeNotFound, err.Error())
		return
	}
	response.OK(c, info)
}
