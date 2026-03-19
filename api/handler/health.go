package handler

import (
	"github.com/MJ-9527/GoMind/config"
	"github.com/MJ-9527/GoMind/pkg/logger"
	"github.com/gin-gonic/gin"
)

// HealthCheckHandler  健康检查接口处理器
// 功能：返回服务状态、版本、环境等信息，用于监控平台检测服务可用性
func HealthCheckHandler(c *gin.Context) {
	// 1. 构造联考检查数据
	healthData := map[string]interface{}{
		"status":  "up",                            //服务状态
		"version": "0.0.1",                         //版本号
		"env":     config.GlobalConfig.Server.Mode, //运行环境
		"port":    config.GlobalConfig.Server.Port, //监听端口
	}

	// 2.记录健康检查日志（排查监控告警问题)
	logger.Info("健康检查请求",
		logger.String("client_ip", c.ClientIP()),
		logger.Any("healthData", healthData),
	)

	Success(c, healthData)
}
