// Package redis 全局 Redis 客户端及常用工具函数
package redis

import (
"context"
"fmt"
"time"

"github.com/MJ-9527/GoMind/config"
"github.com/MJ-9527/GoMind/pkg/logger"
"github.com/redis/go-redis/v9"
"go.uber.org/zap"
)

var Client *redis.Client

// InitRedis 初始化 Redis 客户端
func InitRedis() error {
Client = redis.NewClient(&redis.Options{
Addr:     config.GlobalConfig.Redis.Addr,
Password: config.GlobalConfig.Redis.Password,
DB:       config.GlobalConfig.Redis.DB,
PoolSize: 100, // 连接池大小
})

_, err := Client.Ping(context.Background()).Result()
if err != nil {
return fmt.Errorf("Redis 连接失败：%w", err)
}
logger.Info("Redis 连接成功")
return nil
}

// Close 关闭 Redis 连接
func Close() error {
if Client != nil {
return Client.Close()
}
return nil
}

// ========== 字符串操作 ==========

// Set 设置键值对，带过期时间
// 如果 expire 为 0，则永不过期
func Set(ctx context.Context, key string, value interface{}, expire time.Duration) error {
return Client.Set(ctx, key, value, expire).Err()
}

// Get 获取字符串值
func Get(ctx context.Context, key string) (string, error) {
return Client.Get(ctx, key).Result()
}

// Del 删除一个或多个键
func Del(ctx context.Context, keys ...string) error {
return Client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
result, err := Client.Exists(ctx, key).Result()
if err != nil {
return false, err
}
return result > 0, nil
}

// ========== 哈希操作 ==========

// HSet 设置哈希表字段值
func HSet(ctx context.Context, key string, field string, value interface{}) error {
return Client.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希表字段值
func HGet(ctx context.Context, key string, field string) (string, error) {
return Client.HGet(ctx, key, field).Result()
}

// HMSet 批量设置哈希表字段值
func HMSet(ctx context.Context, key string, values map[string]interface{}) error {
return Client.HMSet(ctx, key, values).Err()
}

// HMGet 批量获取哈希表字段值
func HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
return Client.HMGet(ctx, key, fields...).Result()
}

// HGetAll 获取整个哈希表
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
return Client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希表字段
func HDel(ctx context.Context, key string, fields ...string) error {
return Client.HDel(ctx, key, fields...).Err()
}

// ========== 列表操作 ==========

// LPush 左推入列表
func LPush(ctx context.Context, key string, values ...interface{}) error {
return Client.LPush(ctx, key, values...).Err()
}

// RPush 右推入列表
func RPush(ctx context.Context, key string, values ...interface{}) error {
return Client.RPush(ctx, key, values...).Err()
}

// LPop 左弹出列表
func LPop(ctx context.Context, key string) (string, error) {
return Client.LPop(ctx, key).Result()
}

// RPop 右弹出列表
func RPop(ctx context.Context, key string) (string, error) {
return Client.RPop(ctx, key).Result()
}

// LRange 获取列表指定范围的元素
func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
return Client.LRange(ctx, key, start, stop).Result()
}

// LLen 获取列表长度
func LLen(ctx context.Context, key string) (int64, error) {
return Client.LLen(ctx, key).Result()
}

// ========== 集合操作 ==========

// SAdd 添加元素到集合
func SAdd(ctx context.Context, key string, members ...interface{}) error {
return Client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有元素
func SMembers(ctx context.Context, key string) ([]string, error) {
return Client.SMembers(ctx, key).Result()
}

// SIsMember 判断元素是否在集合中
func SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
return Client.SIsMember(ctx, key, member).Result()
}

// SRem 从集合中删除元素
func SRem(ctx context.Context, key string, members ...interface{}) error {
return Client.SRem(ctx, key, members...).Err()
}

// ========== 有序集合操作 ==========

// ZAdd 添加元素到有序集合
func ZAdd(ctx context.Context, key string, score float64, member interface{}) error {
return Client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

// ZRange 获取有序集合指定范围的元素
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
return Client.ZRange(ctx, key, start, stop).Result()
}

// ZRevRange 倒序获取有序集合指定范围的元素
func ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
return Client.ZRevRange(ctx, key, start, stop).Result()
}

// ZRem 从有序集合中删除元素
func ZRem(ctx context.Context, key string, members ...interface{}) error {
return Client.ZRem(ctx, key, members...).Err()
}

// ZScore 获取元素的分数
func ZScore(ctx context.Context, key string, member string) (float64, error) {
return Client.ZScore(ctx, key, member).Result()
}

// ========== 分布式锁 ==========

// Lock 尝试获取分布式锁
// 参数：key 锁的键，value 锁的值（用于标识持有者），expire 过期时间
// 返回：是否成功获取锁
func Lock(ctx context.Context, key string, value string, expire time.Duration) (bool, error) {
// 使用 SET NX EX 命令实现原子操作
result, err := Client.SetNX(ctx, key, value, expire).Result()
if err != nil {
return false, err
}
return result, nil
}

// Unlock 释放分布式锁
// 参数：key 锁的键，value 锁的值（需要验证是同一个持有者）
// 返回：是否成功释放锁
func Unlock(ctx context.Context, key string, value string) (bool, error) {
// 使用 Lua 脚本保证原子性：先检查值是否匹配，再删除
script := `
if redis.call("get", KEYS[1]) == ARGV[1] then
return redis.call("del", KEYS[1])
else
return 0
end
`
result, err := Client.Eval(ctx, script, []string{key}, value).Int64()
if err != nil {
return false, err
}
return result == 1, nil
}

// TryLockWithRetry 带重试的分布式锁
// 参数：key 锁的键，value 锁的值，expire 过期时间，retryCount 重试次数，retryInterval 重试间隔
// 返回：是否成功获取锁
func TryLockWithRetry(ctx context.Context, key string, value string, expire time.Duration, retryCount int, retryInterval time.Duration) (bool, error) {
for i := 0; i < retryCount; i++ {
locked, err := Lock(ctx, key, value, expire)
if err != nil {
logger.Error("获取分布式锁失败", zap.String("key", key), zap.Error(err))
return false, err
}
if locked {
return true, nil
}
time.Sleep(retryInterval)
}
return false, nil
}

// ========== 限流器 ==========

// RateLimit 基于滑动窗口的限流器
// 参数：key 限流键，maxRequests 最大请求数，window 时间窗口
// 返回：是否允许请求
func RateLimit(ctx context.Context, key string, maxRequests int, window time.Duration) (bool, error) {
now := time.Now().UnixNano()
windowStart := now - window.Nanoseconds()

// 移除窗口外的数据
_ = Client.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", windowStart)).Err()

// 获取当前窗口内的请求数
count, err := Client.ZCard(ctx, key).Result()
if err != nil {
return false, err
}

if count >= int64(maxRequests) {
return false, nil
}

// 添加当前请求
_ = Client.ZAdd(ctx, key, float64(now), fmt.Sprintf("%d", now)).Err()
// 设置过期时间
_ = Client.Expire(ctx, key, window+time.Second).Err()

return true, nil
}

// ========== 缓存辅助函数 ==========

// CacheGetOrSet 缓存穿透保护：获取缓存，不存在则设置并返回
// 参数：key 缓存键，getter 获取数据的函数，expire 过期时间
// 返回：缓存数据和错误信息
func CacheGetOrSet[T any](ctx context.Context, key string, getter func() (T, error), expire time.Duration) (T, error) {
var zero T

// 尝试从缓存获取
data, err := Get(ctx, key)
if err == nil && data != "" {
// TODO: 这里需要根据具体类型进行反序列化，可以使用 JSON
// 简化处理，实际使用时建议传入序列化/反序列化函数
return zero, fmt.Errorf("需要实现反序列化逻辑")
}

// 缓存未命中，执行 getter
value, err := getter()
if err != nil {
return zero, err
}

// TODO: 将结果存入缓存（需要序列化）
// _ = Set(ctx, key, value, expire)

return value, nil
}

// Pipeline 批量操作示例
func BatchSet(ctx context.Context, items map[string]interface{}, expire time.Duration) error {
pipe := Client.Pipeline()
for key, value := range items {
_ = pipe.Set(ctx, key, value, expire).Err()
}
_, err := pipe.Exec(ctx)
return err
}
