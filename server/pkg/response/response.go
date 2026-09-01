// Package response 统一 API 返回封装
//
// 约定：{ code, message, data }，code === 0 为成功；401 表示未授权（前端据此跳登录）。
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeOK        = 0
	CodeBadRequest = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeInternal     = 500
)

// Body 统一返回结构
type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Page 分页结构
type Page struct {
	List     interface{} `json:"list"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
}

// OK 成功返回 data
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: CodeOK, Message: "success", Data: data})
}

// OKPage 成功返回分页数据
func OKPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	OK(c, Page{List: list, Total: total, Page: page, PageSize: pageSize})
}

// OKMsg 仅返回成功消息
func OKMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Body{Code: CodeOK, Message: msg})
}

// Fail 返回失败，HTTP 200 + 错误 code（前端按 res.code 判断）
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Body{Code: code, Message: msg})
}

// FailWith 返回失败 + data
func FailWith(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: code, Message: msg, Data: data})
}

// Unauthorized 401（带 HTTP 状态码，便于网关 / 拦截器识别）
func Unauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, Body{Code: CodeUnauthorized, Message: msg})
}

// Forbidden 403
func Forbidden(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusForbidden, Body{Code: CodeForbidden, Message: msg})
}

// Error 系统内部错误（HTTP 500）
func Error(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, Body{
		Code: CodeInternal, Message: err.Error(),
	})
}

// BadRequest 参数错误
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Body{Code: CodeBadRequest, Message: msg})
}
