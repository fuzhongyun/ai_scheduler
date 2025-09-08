//go:build wireinject
// +build wireinject

package main

import (
	"ai_scheduler/internal/config"
	"ai_scheduler/internal/services"
	"ai_scheduler/internal/tools"
	"ai_scheduler/pkg/ollama"
	"ai_scheduler/pkg/types"

	"github.com/google/wire"
)

// InitializeApp 初始化应用程序
func InitializeApp(configPath string) (*App, error) {
	wire.Build(
		// 配置
		config.LoadConfig,

		// Ollama客户端
		provideOllamaClient,

		// 工具管理器
		provideToolsConfig,
		tools.NewManager,

		// 路由服务
		provideRouterService,

		// 应用程序
		NewApp,
	)
	return &App{}, nil
}

// provideOllamaClient 提供Ollama客户端
func provideOllamaClient(cfg *config.Config) types.AIClient {
	client, _ := ollama.NewClient(&cfg.Ollama)
	return client
}

// provideToolsConfig 提供工具配置
func provideToolsConfig(cfg *config.Config) *config.ToolsConfig {
	return &cfg.Tools
}

// provideRouterService 提供路由服务
func provideRouterService(aiClient types.AIClient, toolManager *tools.Manager) types.RouterService {
	return services.NewRouterService(aiClient, toolManager)
}

// App 应用程序结构
type App struct {
	Config        *config.Config
	RouterService types.RouterService
}

// NewApp 创建应用程序
func NewApp(cfg *config.Config, routerService types.RouterService) *App {
	return &App{
		Config:        cfg,
		RouterService: routerService,
	}
}
