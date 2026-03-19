// Package handler 接口处理器模块
// 职责：实现HTTP接口的具体逻辑，处理请求参数，返回响应
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ========== 统一返回格式定义 ==========

// Response 接口统一返回结构体
// 注意：
// 1. JSON字段名用小写+下划线（前端友好）
// 2. 所有字段首字母大写（Gin序列化需要）
// 3. Data为interface{}类型，支持返回任意结构数据
type Response struct {
	Code    int         `json:"code"`    // 业务错误码（0=成功，非0=失败）
	Message string      `json:"message"` // 提示信息（用户可感知）
	Data    interface{} `json:"data"`    // 业务数据（成功时返回，失败时为null）
}

// ========== 全局错误码定义（统一管理，便于前端适配） ==========
const (
	CodeSuccess       = 0   // 成功
	CodeServerError   = 500 // 服务器内部错误
	CodeInvalidParams = 400 // 参数错误（格式/必填/长度等）
	CodeNotFound      = 404 // 资源不存在
	CodeUnauthorized  = 401 // 未授权（后续扩展认证用）
	CodeForbidden     = 403 // 权限不足（后续扩展权限用）
)

// ========== 响应工具函数（简化重复代码，统一响应风格） ==========

// Success 成功响应
// 参数：
//
//	c: Gin上下文
//	data: 要返回的业务数据（任意类型）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Fail 失败响应
// 参数：
//
//	c: Gin上下文
//	code: 业务错误码（使用上面定义的常量）
//	message: 错误提示信息
func Fail(c *gin.Context, code int, message string) {
	// 注意：HTTP状态码统一返回200，业务错误用code区分（前端更易处理）
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ServerError 服务器内部错误快捷响应
// 参数：
//
//	c: Gin上下文
//	err: 原始错误（用于日志，不返回给前端）
func ServerError(c *gin.Context, err error) {
	// 前端只返回通用提示，具体错误日志内部记录
	Fail(c, CodeServerError, "服务器内部错误，请稍后重试")
}

// InvalidParams 非法参数快捷响应
// 参数：
//
//	c: Gin上下文
//	message: 具体的参数错误提示
func InvalidParams(c *gin.Context, message string) {
	Fail(c, CodeInvalidParams, "参数错误："+message)
}
