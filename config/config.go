// Package config 配置管理模块
// 职责：加载/解析配置文件，提供全局可访问的配置对象
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// AppConfig 全局配置结构体
type AppConfig struct {
	Server struct {
		Host string `yaml:"host"` //服务监听地址
		Port int    `yaml:"port"` //服务监听端口
		Mode string `yaml:"mode"` //运行环境(dev/test/prod)
	} `yaml:"server"`

	Log struct {
		Level      string `yaml:"level"`       //日志级别
		Path       string `yaml:"path"`        //日志文件路径
		MaxSize    int    `yaml:"max_size"`    //单个日志文件最大大小(MB)
		MaxBackups int    `yaml:"max_backups"` //单个文件保留数量
	} `yaml:"log"`

	AI struct {
		Model      string `yaml:"model"`       // 大模型名称（gpt-3.5-turbo/qwen-turbo）
		Timeout    int    `yaml:"timeout"`     //请求超时时间（秒）
		MaxRetries int    `yaml:"max_retries"` //最大重试次数
	} `yaml:"ai"`

	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`

	RateLimit struct {
		MaxRequests int `yaml:"max_requests"`
	} `yaml:"rate_limit"`
}

// GlobalConfig 全局配置实例
// 注意：全局变量仅用于配置这类核心公共数据，禁止滥用
var GlobalConfig AppConfig

// LoadConfig 加载配置文件
// 参数：configPath 配置文件路径
// 返回：错误信息（所有对外函数必须返回error，便于上层处理）
func LoadConfig(configPath string) error {
	//1.校验文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s,err: %w", configPath, err)
	}

	//2.读取文件内容
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	//3.解析YAML到结构体
	if err := yaml.Unmarshal(data, &GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	//4.初始化日志目录（提前创建，避免日志写入失败）
	logDir := filepath.Dir(GlobalConfig.Log.Path)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	//5. 优先加载.env文件，生产环境可通过系统环境变量覆盖
	if err := godotenv.Load(); err != nil {
		// 非致命错误，仅打印警告（生产环境可能不使用.env文件）
		fmt.Printf("加载.env文件警告: %v\n", err)
	}

	return nil
}
