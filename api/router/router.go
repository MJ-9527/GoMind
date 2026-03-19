// Package router 路由管理模块
// 职责：统一注册所有HTTP路由，划分路由分组，挂载中间件
package router

import (
	"github.com/MJ-9527/GoMind/api/handler"
	"github.com/MJ-9527/GoMind/api/middleware"
	"github.com/MJ-9527/GoMind/config"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// 1.根据环境设置Gin模式（大厂规范：生产环境禁用Debug）
	if config.GlobalConfig.Server.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 2. 创建Gin引擎(禁用默认Recovery,后续自定义)
	r := gin.New()

	// 3. 挂载全局中间键（日志+异常回复，顺序：日志早前）
	r.Use(middleware.LogMiddleware()) //自定义日志中间键
	r.Use(gin.Recovery())             //异常回复，避免服务崩溃

	// 4. 健康检查接口
	r.GET("/health", handler.HealthCheckHandler)

	// 5. 业务接口分组（版本化管理，便于迭代）
	v1 := r.Group("/api/v1")
	{
		// ========= Agent 路由 =========
		agent := v1.Group("/agent")
		{
			agent.POST("chat", handler.AgentChatHandler) // Agent对话接口
		}
	}

	return r
}
