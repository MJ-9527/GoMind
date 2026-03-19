package middleware

import (
	"time"

	"github.com/MJ-9527/GoMind/pkg/logger"
	"github.com/gin-gonic/gin"
)

// LogMiddleware 自定义请求日志中间件
// 功能：记录每个请求的方法、路径、耗时、状态码等信息
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 记录请求开始的时间
		startTime := time.Now()

		// 2. 处理请求（调用后续中间件/处理器）
		c.Next()

		// 3. 记录请求结束后的信息
		latency := time.Since(startTime) //耗时
		logger.Info("HTTP 请求日志",
			logger.String("method", c.Request.Method),          // 请求方法
			logger.String("path", c.Request.URL.Path),          // 请求路径
			logger.Int("status", c.Writer.Status()),            // 响应状态码
			logger.Duration("latency", latency),                // 耗时
			logger.String("client_ip", c.ClientIP()),           // 客户端IP
			logger.String("user_agent", c.Request.UserAgent()), // 客户端UA
		)
	}
}
