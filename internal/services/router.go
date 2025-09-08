package services

import (
	"ai_scheduler/internal/tools"
	"ai_scheduler/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// RouterService 智能路由服务
type RouterService struct {
	aiClient    types.AIClient
	toolManager *tools.Manager
}

// NewRouterService 创建路由服务
func NewRouterService(aiClient types.AIClient, toolManager *tools.Manager) *RouterService {
	return &RouterService{
		aiClient:    aiClient,
		toolManager: toolManager,
	}
}

// Route 执行智能路由
func (r *RouterService) Route(ctx context.Context, req *types.ChatRequest) (*types.ChatResponse, error) {
	// 构建消息
	messages := []types.Message{
		{
			Role:    "system",
			Content: r.buildSystemPrompt(),
		},
		{
			Role:    "user",
			Content: req.Message,
		},
	}

	// 获取工具定义
	toolDefinitions := r.toolManager.GetToolDefinitions()

	// 第一次调用AI，获取是否需要使用工具
	response, err := r.aiClient.Chat(ctx, messages, toolDefinitions)
	if err != nil {
		return nil, fmt.Errorf("failed to chat with AI: %w", err)
	}

	// 如果没有工具调用，直接返回
	if len(response.ToolCalls) == 0 {
		return response, nil
	}

	// 执行工具调用
	toolResults, err := r.toolManager.ExecuteToolCalls(ctx, response.ToolCalls)
	if err != nil {
		return nil, fmt.Errorf("failed to execute tools: %w", err)
	}

	// 构建包含工具结果的消息
	messages = append(messages, types.Message{
		Role:    "assistant",
		Content: response.Message,
	})

	// 添加工具调用结果
	for _, toolResult := range toolResults {
		toolResultStr, _ := json.Marshal(toolResult.Result)
		messages = append(messages, types.Message{
			Role:    "tool",
			Content: fmt.Sprintf("Tool %s result: %s", toolResult.Function.Name, string(toolResultStr)),
		})
	}

	// 第二次调用AI，生成最终回复
	finalResponse, err := r.aiClient.Chat(ctx, messages, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate final response: %w", err)
	}

	// 合并工具调用信息到最终响应
	finalResponse.ToolCalls = toolResults

	log.Printf("Router processed request: %s, used %d tools", req.Message, len(toolResults))

	return finalResponse, nil
}

// buildSystemPrompt 构建系统提示词
func (r *RouterService) buildSystemPrompt() string {
	prompt := `你是一个智能助手，可以帮助用户解决各种问题。你有以下工具可以使用：

`

	// 添加工具描述
	tools := r.toolManager.GetAllTools()
	for _, tool := range tools {
		prompt += fmt.Sprintf("- %s: %s\n", tool.Name(), tool.Description())
	}

	prompt += `
请根据用户的问题，判断是否需要使用工具。如果需要，请调用相应的工具获取信息，然后基于工具返回的结果给出完整的回答。

注意事项：
1. 只有在确实需要获取实时信息或进行计算时才使用工具
2. 如果用户只是普通聊天，不需要使用工具
3. 使用工具后，请基于工具返回的结果给出自然、友好的回复
4. 如果工具执行出错，请告知用户并提供替代建议`

	return prompt
}