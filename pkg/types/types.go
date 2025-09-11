package types

import (
	"context"
	"encoding/json"
)

// ChatRequest 聊天请求
type ChatRequest struct {
	UserInput       string          `json:"user_input" binding:"required"`
	Caller          string          `json:"caller" binding:"required"`
	SessionID       string          `json:"session_id"`
	ChatRequestMeta ChatRequestMeta `json:"meta,omitempty"`
}

// ChatRequestMeta 聊天请求元数据
type ChatRequestMeta struct {
	Authorization string `json:"authorization"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Data     any    `json:"data,omitempty"`
	TaskCode string `json:"task_code,omitempty"`
}

// ToolCall 工具调用
type ToolCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Function FunctionCall    `json:"function"`
	Result   json.RawMessage `json:"result,omitempty"`
}

// FunctionCall 函数调用
type FunctionCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ToolDefinition 工具定义
type ToolDefinition struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

// FunctionDef 函数定义
type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Tool 工具接口
type Tool interface {
	Name() string
	Description() string
	Definition() ToolDefinition
	Execute(ctx context.Context, args json.RawMessage) (interface{}, error)
}

// AIClient AI客户端接口
type AIClient interface {
	Chat(ctx context.Context, messages []Message, tools []ToolDefinition) (*ChatResponse, error)
}

// Message 消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// RouterService 路由服务接口
type RouterService interface {
	Route(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
}
