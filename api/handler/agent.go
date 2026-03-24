package handler

import (
	"context"

	"github.com/MJ-9527/GoMind/pkg/ai"
	"github.com/MJ-9527/GoMind/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// ========== 请求参数结构体（大厂规范：参数校验） ==========

// AgentChatRequest Agent对话请求参数
type AgentChatRequest struct {
	UserInput string `json:"user_input" binding:"required,min=1,max=1000"` // 用户输入（必填，长度1-1000）
	SessionID string `json:"session_id" binding:"required,min=1"`          // 会话ID（必填，区分不同用户）
}

// ========== 核心处理器 ==========

// AgentChatHandler AI Agent对话接口处理器
func AgentChatHandler(c *gin.Context) {
	// 1. 参数绑定与校验（Gin内置校验，大厂必备）
	var req AgentChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Agent对话接口参数校验失败",
			zap.String("session_id", req.SessionID),
			zap.Error(err),
		)
		Fail(c, CodeInvalidParams, "参数错误："+err.Error())
		return
	}

	// 2. 记录请求日志（上下文信息完整）
	logger.Info("接收Agent对话请求",
		zap.String("session_id", req.SessionID),
		zap.String("user_input", req.UserInput),
	)

	// 3. 构建大模型对话消息
	messages := []openai.ChatCompletionMessage{
		// System指令：定义Agent角色（可配置化）
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你是一个专业的Go开发助手Agent，回答简洁、准确、实用，只讲干货。",
		},
	}

	// 4. 加载历史会话并拼接本次用户输入
	history := ai.GetMessages(req.SessionID)
	if len(history) > 0 {
		messages = append(messages, history...)
	}

	currentUserMsg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.UserInput,
	}
	messages = append(messages, currentUserMsg)

	// 5. 调用AI客户端
	respContent, err := ai.GlobalAIClient.Chat(context.Background(), messages)
	if err != nil {
		logger.Error("Agent对话接口调用AI失败",
			zap.String("session_id", req.SessionID),
			zap.Error(err),
		)
		ServerError(c, err)
		return
	}

	// 6. 保存本轮对话到会话（供后续请求读取）
	ai.SaveMessage(req.SessionID, currentUserMsg)
	ai.SaveMessage(req.SessionID, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: respContent,
	})

	// 7. 构造响应数据
	respData := map[string]interface{}{
		"session_id": req.SessionID,
		"reply":      respContent,
	}

	// 8. 返回成功响应
	logger.Info("Agent对话接口响应成功", zap.String("session_id", req.SessionID))
	Success(c, respData)
}
