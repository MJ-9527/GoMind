// Package main 服务器启动入口
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/MJ-9527/GoMind/config"
	"github.com/MJ-9527/GoMind/pkg/logger"
	"go.uber.org/zap"
)

// 程序版本
const version = "0.0.1"

func main() {
	// ========== 1. 加载配置 ==========
	if err := config.LoadConfig("./config/config.yaml"); err != nil {
		// 启动失败直接Fatal,打印错误并退出
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// ========== 2. 初始化日志 ==========
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	// ========== 3. 启动日志（标准化）==========
	logger.Info("服务启动",
		zap.String("version", version),
		zap.String("env", config.GlobalConfig.Server.Mode),
		zap.String("host", config.GlobalConfig.Server.Host),
		zap.Int("port", config.GlobalConfig.Server.Port),
	)

	// =========== 4. 监听退出信号（优雅关闭）==========
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// ========== 5. 退出日志（标准化） ==========
	logger.Info("服务开始优雅退出")
	// 后续可添加：关闭数据库连接、清理资源等逻辑
	logger.Info("服务已退出")
}
