package utils

import "github.com/gin-gonic/gin"

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	})
}

// Created 创建成功响应 (201)
func Created(c *gin.Context, data interface{}) {
	c.JSON(201, Response{
		Code: 201,
		Msg:  "created successfully",
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, msg string) {
	Error(c, 400, msg)
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, msg string) {
	Error(c, 401, msg)
}

// Forbidden 403 错误
func Forbidden(c *gin.Context, msg string) {
	Error(c, 403, msg)
}

// NotFound 404 错误
func NotFound(c *gin.Context, msg string) {
	Error(c, 404, msg)
}

// Conflict 409 错误
func Conflict(c *gin.Context, msg string) {
	Error(c, 409, msg)
}

// InternalServerError 500 错误
func InternalServerError(c *gin.Context, msg string) {
	Error(c, 500, msg)
}
