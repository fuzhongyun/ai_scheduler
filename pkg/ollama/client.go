package ollama

import (
	"ai_scheduler/internal/config"
	"ai_scheduler/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ollama/ollama/api"
)

// Client Ollama客户端适配器
type Client struct {
	client *api.Client
	config *config.OllamaConfig
}

// NewClient 创建新的Ollama客户端
func NewClient(config *config.OllamaConfig) (*Client, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama client: %w", err)
	}

	// 设置自定义HTTP客户端
	// 创建HTTP客户端（如果需要的话）
	_ = &http.Client{
		Timeout: config.Timeout,
	}

	return &Client{
		client: client,
		config: config,
	}, nil
}

// Chat 实现聊天功能
func (c *Client) Chat(ctx context.Context, messages []types.Message, tools []types.ToolDefinition) (*types.ChatResponse, error) {
	// 构建聊天请求
	req := &api.ChatRequest{
		Model: c.config.Model,
		Messages: make([]api.Message, len(messages)),
		Stream:   new(bool), // 设置为false，不使用流式响应
	}

	// 转换消息格式
	for i, msg := range messages {
		req.Messages[i] = api.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 添加工具定义
	if len(tools) > 0 {
		req.Tools = make([]api.Tool, len(tools))
		for i, tool := range tools {
			toolData, _ := json.Marshal(tool)
			var apiTool api.Tool
			json.Unmarshal(toolData, &apiTool)
			req.Tools[i] = apiTool
		}
	}

	// 发送请求
	responseChan := make(chan api.ChatResponse)
	errorChan := make(chan error)

	go func() {
		err := c.client.Chat(ctx, req, func(resp api.ChatResponse) error {
			responseChan <- resp
			return nil
		})
		if err != nil {
			errorChan <- err
		}
		close(responseChan)
		close(errorChan)
	}()

	// 等待响应
	select {
	case resp := <-responseChan:
		return c.convertResponse(&resp), nil
	case err := <-errorChan:
		return nil, fmt.Errorf("chat request failed: %w", err)
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(c.config.Timeout):
		return nil, fmt.Errorf("chat request timeout")
	}
}

// convertResponse 转换响应格式
func (c *Client) convertResponse(resp *api.ChatResponse) *types.ChatResponse {
	result := &types.ChatResponse{
		Message:  resp.Message.Content,
		Finished: resp.Done,
	}

	// 转换工具调用
	if len(resp.Message.ToolCalls) > 0 {
		result.ToolCalls = make([]types.ToolCall, len(resp.Message.ToolCalls))
		for i, toolCall := range resp.Message.ToolCalls {
			result.ToolCalls[i] = types.ToolCall{
				ID:   fmt.Sprintf("call_%d", i),
				Type: "function",
				Function: types.FunctionCall{
					Name:      toolCall.Function.Name,
					Arguments: json.RawMessage(fmt.Sprintf("%v", toolCall.Function.Arguments)),
				},
			}
		}
	}

	return result
}