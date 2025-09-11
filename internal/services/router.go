package services

import (
	"ai_scheduler/internal/constants"
	"ai_scheduler/internal/tools"
	"ai_scheduler/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
		// {
		// 	Role:    "system",
		// 	Content: r.buildSystemPrompt(),
		// },
		{
			Role:    "assistant",
			Content: r.buildIntentPrompt(req.UserInput),
		},
		{
			Role:    "user",
			Content: req.UserInput,
		},
	}

	// 第1次调用AI，获取用户意图
	intentResponse, err := r.aiClient.Chat(ctx, messages, nil)
	if err != nil {
		return nil, fmt.Errorf("AI响应失败: %w", err)
	}

	// 从AI响应中提取意图
	intent := r.extractIntent(intentResponse)
	if intent == "" {
		return nil, fmt.Errorf("未识别到用户意图")
	}

	switch intent {
	case "order_diagnosis":
		// 订单诊断意图
		return r.handleOrderDiagnosis(ctx, req, messages)
	case "knowledge_qa":
		// 知识问答意图
		return r.handleKnowledgeQA(ctx, req, messages)
	default:
		// 未知意图
		return nil, fmt.Errorf("意图识别失败，请明确您的需求呢，我可以为您")
	}

	// 获取工具定义
	toolDefinitions := r.toolManager.GetToolDefinitions(constants.Caller(req.Caller))

	// 第2次调用AI，获取是否需要使用工具
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

	log.Printf("Router processed request: %s, used %d tools", req.UserInput, len(toolResults))

	return finalResponse, nil
}

// buildSystemPrompt 构建系统提示词
func (r *RouterService) buildSystemPrompt() string {
	prompt := `你是一个智能路由系统，你的任务是根据用户输入判断用户的意图，并且执行对应的任务。`

	return prompt
}

// buildIntentPrompt 构建意图识别提示词
func (r *RouterService) buildIntentPrompt(userInput string) string {
	prompt := `请分析以下用户输入，判断用户的意图类型。

用户输入：{user_input}

意图类型说明：
1. order_diagnosis - 订单诊断：用户想要查询、诊断或了解订单相关信息  
2. knowledge_qa - 知识问答：用户想要进行一般性问答或获取知识信息

- 当用户意图不够清晰且不匹配 knowledge_qa 以外意图时，使用knowledge_qa
- 当用户意图非常不清晰时使用 unknown 

请只返回以下格式的JSON：
{
    "intent": "order_diagnosis" | "knowledge_qa" | "unknown",
    "confidence": 0.0-1.0,
    "reasoning": "判断理由"
}
`

	prompt = strings.ReplaceAll(prompt, "{user_input}", userInput)

	return prompt
}

// extractIntent 从AI响应中提取意图
func (r *RouterService) extractIntent(response *types.ChatResponse) string {
	if response == nil || response.Message == "" {
		return ""
	}

	// 尝试解析JSON
	var intent struct {
		Intent     string `json:"intent"`
		Confidence string `json:"confidence"`
		Reasoning  string `json:"reasoning"`
	}
	err := json.Unmarshal([]byte(response.Message), &intent)
	if err != nil {
		log.Printf("Failed to parse intent JSON: %v", err)
		return ""
	}

	return intent.Intent
}

// handleOrderDiagnosis 处理订单诊断意图
func (r *RouterService) handleOrderDiagnosis(ctx context.Context, req *types.ChatRequest, messages []types.Message) (*types.ChatResponse, error) {
	// 调用订单详情工具
	orderDetailTool, ok := r.toolManager.GetTool("zltxOrderDetail")
	if orderDetailTool == nil || !ok {
		return nil, fmt.Errorf("order detail tool not found")
	}
	orderDetailTool.Execute(ctx, json.RawMessage{})

	// 获取相关工具定义
	toolDefinitions := r.toolManager.GetToolDefinitions(constants.Caller(req.Caller))

	// 调用AI，获取是否需要使用工具
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

	return nil, nil
}

// handleKnowledgeQA 处理知识问答意图
func (r *RouterService) handleKnowledgeQA(ctx context.Context, req *types.ChatRequest, messages []types.Message) (*types.ChatResponse, error) {

	return nil, nil
}
