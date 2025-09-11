package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Ollama  OllamaConfig  `mapstructure:"ollama"`
	Tools   ToolsConfig   `mapstructure:"tools"`
	Logging LoggingConfig `mapstructure:"logging"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// OllamaConfig Ollama配置
type OllamaConfig struct {
	BaseURL string        `mapstructure:"base_url"`
	Model   string        `mapstructure:"model"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// ToolsConfig 工具配置
type ToolsConfig struct {
	Weather         ToolConfig `mapstructure:"weather"`
	Calculator      ToolConfig `mapstructure:"calculator"`
	ZltxOrderDetail ToolConfig `mapstructure:"zltxOrderDetail"`
	ZltxOrderLog    ToolConfig `mapstructure:"zltxOrderLog"`
	Knowledge       ToolConfig `mapstructure:"knowledge"`
}

// ToolConfig 单个工具配置
type ToolConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	BaseURL   string `mapstructure:"base_url"`
	APIKey    string `mapstructure:"api_key"`
	BizSystem string `mapstructure:"biz_system"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("ollama.base_url", "http://localhost:11434")
	viper.SetDefault("ollama.model", "llama2")
	viper.SetDefault("ollama.timeout", "30s")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
