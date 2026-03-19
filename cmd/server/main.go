// Package main 服务器启动入口
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MJ-9527/GoMind/api/router"
	"github.com/MJ-9527/GoMind/config"
	"github.com/MJ-9527/GoMind/pkg/ai"
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
	defer logger.Logger.Sync()

	// ========== 初始化AI客户端 ==========
	if err := ai.InitAIClient(); err != nil {
		logger.Fatal("初始化AI客户端失败", zap.Error(err))
	}

	// ========= 3. 初始化路由 ==========
	r := router.InitRouter()

	// ========= 4. 启动HTTP服务 =========
	addr := fmt.Sprintf("%s:%d", config.GlobalConfig.Server.Host, config.GlobalConfig.Server.Port)
	logger.Info("HTTP服务启动中", zap.String("addr", addr))

	// 创建HTTP服务器实例（便于优雅关闭）
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 异步启动服务（避免阻塞信号监听）
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP服务启动失败", zap.Error(err))
		}
	}()

	// ========== 5. 监听退出信号 ==========
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// ========== 6. 优雅关闭服务 ==========
	logger.Info("开始优雅关闭HTTP服务")
	if err := server.Shutdown(nil); err != nil {
		logger.Error("HTTP服务关闭失败", zap.Error(err))
	}
	logger.Info("HTTP服务已成功关闭")
}
