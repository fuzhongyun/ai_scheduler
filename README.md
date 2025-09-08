# AI Scheduler - 智能路由调度系统

基于Go语言开发的智能AI助手，支持Function Calling和工具调用，可以智能路由用户请求到合适的工具进行处理。

## 功能特性

- 🤖 **智能对话**: 基于Ollama的AI对话能力
- 🔧 **工具调用**: 支持天气查询、计算器等工具
- 🎯 **智能路由**: 自动判断是否需要调用工具
- 📚 **API文档**: 集成Swagger文档
- ⚡ **高性能**: 基于Gin框架的HTTP服务
- 🏗️ **依赖注入**: 使用Wire进行依赖管理

## 项目结构

```
ai_scheduler/
├── cmd/                    # 应用程序入口
│   └── main.go
├── internal/               # 内部包
│   ├── config/            # 配置管理
│   ├── handlers/          # HTTP处理器
│   ├── models/            # 数据模型
│   ├── services/          # 业务服务
│   ├── tools/             # 工具实现
│   └── wire/              # 依赖注入
├── pkg/                   # 公共包
│   ├── ollama/           # Ollama客户端
│   └── types/            # 类型定义
├── docs/                  # API文档
├── config.yaml           # 配置文件
└── go.mod                # Go模块
```

## 快速开始

### 1. 环境要求

- Go 1.23+
- Ollama服务运行中

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置文件

编辑 `config.yaml` 文件，确保Ollama服务地址正确：

```yaml
server:
  port: "8080"
  host: "localhost"

ollama:
  base_url: "http://localhost:11434"
  model: "llama2"
  timeout: 30s

tools:
  weather:
    enabled: true
    mock_data: true
  calculator:
    enabled: true
    mock_data: false

logging:
  level: "info"
  format: "json"
```

### 4. 启动服务

```bash
go run cmd/main.go
```

服务启动后，可以访问：

- API服务: http://localhost:8080/api/v1/chat
- Swagger文档: http://localhost:8080/swagger/index.html
- 健康检查: http://localhost:8080/health

## API使用示例

### 聊天接口

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "北京今天天气怎么样？",
    "model": "llama2"
  }'
```

### 计算器示例

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "计算 15 + 25 * 3",
    "model": "llama2"
  }'
```

## 支持的工具

### 1. 天气查询工具
- 功能：查询指定城市的天气信息
- 示例："北京今天天气怎么样？"

### 2. 计算器工具
- 功能：执行数学计算
- 支持：加减乘除、幂运算
- 示例："计算 2 + 3 * 4"

## 开发说明

### 添加新工具

1. 在 `internal/tools/` 目录下创建新工具文件
2. 实现 `types.Tool` 接口
3. 在 `tools.Manager` 中注册新工具

### 配置管理

配置文件使用Viper加载，支持环境变量覆盖。

### 依赖注入

使用Google Wire进行依赖注入，修改依赖关系后需要重新生成代码。

## 许可证

MIT License