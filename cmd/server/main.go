package main

import (
	"ai_scheduler/internal/handlers"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "ai_scheduler/docs"

	"github.com/gin-gonic/gin"
)

// @title AI Scheduler API
// @version 1.0
// @description 智能路由调度系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	port := flag.String("port", "8080", "服务端口")
	flag.Parse()

	// 初始化应用程序
	app, err := InitializeApp(*configPath)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// 设置Gin模式为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 设置路由
	router := handlers.SetupRoutes(app.RouterService)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         ":" + *port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("Starting server on port %s", *port)
		log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", *port)
		log.Printf("Health check: http://localhost:%s/health", *port)
		log.Printf("Chat API: http://localhost:%s/api/v1/chat", *port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server exited gracefully")
	}
}
