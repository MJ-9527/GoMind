// Package ai 会话管理：存储/读取对话历史
package ai

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MJ-9527/GoMind/pkg/logger"
	"github.com/MJ-9527/GoMind/pkg/redis"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

const sessionExpire = time.Hour * 24

// SaveMessage 保存单条消息
func SaveMessage(sessionID string, msg openai.ChatCompletionMessage) {
	key := "agent:session:" + sessionID

	// 读取已有列表
	list, err := redis.Client.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		logger.Error("读取会话失败", zap.Error(err))
		return
	}

	// 反序列化
	var messages []openai.ChatCompletionMessage
	for _, item := range list {
		var m openai.ChatCompletionMessage
		if json.Unmarshal([]byte(item), &m) == nil {
			messages = append(messages, m)
		}
	}

	// 限制最大长度（防爆炸）
	if len(messages) >= 20 {
		redis.Client.LPop(context.Background(), key)
	}

	// 新增
	bs, _ := json.Marshal(msg)
	redis.Client.RPush(context.Background(), key, bs)
	redis.Client.Expire(context.Background(), key, sessionExpire)
}

// GetMessages 获取会话全部消息
func GetMessages(sessionID string) []openai.ChatCompletionMessage {
	key := "agent:session:" + sessionID
	list, err := redis.Client.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		logger.Error("读取会话失败", zap.Error(err))
		return nil
	}

	var messages []openai.ChatCompletionMessage
	for _, item := range list {
		var m openai.ChatCompletionMessage
		if json.Unmarshal([]byte(item), &m) == nil {
			messages = append(messages, m)
		}
	}
	return messages
}
