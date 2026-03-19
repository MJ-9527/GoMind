// Package ai AI客户端模块
// 职责：封装大模型调用逻辑，隔离第三方SDK，提供统一的AI能力接口
package ai

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MJ-9527/GoMind/config"
	"github.com/MJ-9527/GoMind/pkg/logger"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// AIClient AI客户端实例
type AIClient struct {
	client *openai.Client
	config config.AppConfig
}

// GlobalAIClient 全局AI客户端（仅核心组件使用全局，业务层通过依赖注入）
var GlobalAIClient *AIClient

// InitAIClient 初始化AI客户端
// 依赖：需先加载config.GlobalConfig
func InitAIClient() error {
	// 1. 从环境变量获取敏感信息（优先级：系统环境变量 > .env文件）
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("未配置OPENAI_API_KEY环境变量")
	}

	// 2. 初始化OpenAI客户端配置
	clientConfig := openai.DefaultConfig(apiKey)
	// 支持自定义BaseURL（适配国产大模型）
	if baseURL := os.Getenv("OPENAI_BASE_URL"); baseURL != "" {
		clientConfig.BaseURL = baseURL
	}

	// 3. 创建客户端实例
	client := openai.NewClientWithConfig(clientConfig)
	GlobalAIClient = &AIClient{
		client: client,
		config: config.GlobalConfig,
	}

	logger.Info("AI客户端初始化成功", zap.String("model", config.GlobalConfig.AI.Model))
	return nil
}

// Chat AI对话核心方法
// 参数：
//
//	ctx: 上下文（用于超时/取消）
//	messages: 对话消息列表
//
// 返回：
//
//	string: AI回复内容
//	error: 错误信息
func (a *AIClient) Chat(ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	// 1. 设置超时（叠加配置的超时时间）
	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.config.AI.Timeout)*time.Second)
	defer cancel()

	// 2. 构建请求参数
	req := openai.ChatCompletionRequest{
		Model:    a.config.AI.Model,
		Messages: messages,
	}

	// 3. 重试机制（大厂高可用必备）
	var resp openai.ChatCompletionResponse
	var err error
	for i := 0; i < a.config.AI.MaxRetries; i++ {
		resp, err = a.client.CreateChatCompletion(ctx, req)
		if err == nil {
			break // 成功则退出重试
		}

		// 记录重试日志
		logger.Warn("AI对话请求重试",
			zap.Int("retry_times", i+1),
			zap.Error(err),
		)

		// 指数退避（避免频繁重试触发限流）
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	// 4. 最终错误处理
	if err != nil {
		logger.Error("AI对话请求失败（重试次数用尽）", zap.Error(err))
		return "", fmt.Errorf("AI请求失败: %w", err)
	}

	// 5. 校验响应
	if len(resp.Choices) == 0 {
		logger.Warn("AI对话响应无内容")
		return "", nil
	}

	return resp.Choices[0].Message.Content, nil
}
