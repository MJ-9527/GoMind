// Package logger 日志组件
// 职责：提供结构化、可配置的日志能力，封装zap避免上层直接依赖第三方库
package logger

import (
	"os"
	"time"

	"github.com/MJ-9527/GoMind/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 全局日志实例
var Logger *zap.Logger

// InitLogger 初始化日志组件
// 依赖：需先加载config.GlobalConfig
// 返回：错误信息
func InitLogger() error {
	// 1. 从配置读取日志级别
	level := zapcore.InfoLevel
	switch config.GlobalConfig.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel //默认级别
	}

	// 2. 配置日志编码器（结构化JSON格式，大厂日志采用标准）
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"                     //时间字段名
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   //时间格式标准话
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder //级别大写(INFO/ERROR)
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 3. 配置日志写入器（按大小切割，自动清理旧日志）
	writer := &lumberjack.Logger{
		Filename:   config.GlobalConfig.Log.Path,       //日志文件路径
		MaxSize:    config.GlobalConfig.Log.MaxSize,    //单个文件最大MB
		MaxBackups: config.GlobalConfig.Log.MaxBackups, //保留文件数
		Compress:   true,                               //压缩旧日志
	}

	// 4. 构建zap核心
	core := zapcore.NewCore(
		encoder, //编码器
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer), zapcore.AddSync(os.Stdout)), // 同时输出到文件+控制台  （写入器）
		level, //级别启动器
	)
	// 5. 创建logger实例（添加调用者信息，便于定位问题）
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(Logger) // 替换zap全局logger

	return nil
}

// ========== 基础日志方法封装 ==========

// Debug 封装常用日志方法（简化上层调用，统一风格）
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// ========== 日志字段构造函数封装 ==========

func String(key string, value string) zap.Field {
	return zap.String(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

func Uint(key string, value uint) zap.Field {
	return zap.Uint(key, value)
}

func Uint64(key string, value uint64) zap.Field {
	return zap.Uint64(key, value)
}

func Float32(key string, value float32) zap.Field {
	return zap.Float32(key, value)
}

func Float64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

func Duration(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func Time(key string, value time.Time) zap.Field {
	return zap.Time(key, value)
}

func ErrorField(key string, value error) zap.Field {
	return zap.NamedError(key, value)
}
