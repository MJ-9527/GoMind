// Package redis 全局Redis客户端
package redis

import (
	"context"

	"github.com/MJ-9527/GoMind/config"
	"github.com/MJ-9527/GoMind/pkg/logger"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func InitRedis() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.GlobalConfig.Redis.Addr,
		Password: config.GlobalConfig.Redis.Password,
		DB:       config.GlobalConfig.Redis.DB,
	})

	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	logger.Info("Redis 连接成功")
	return nil
}
