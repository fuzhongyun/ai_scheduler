package tools

import (
	"ai_scheduler/internal/config"
	"ai_scheduler/internal/constants"
	"ai_scheduler/pkg/types"
	"context"
	"encoding/json"
	"fmt"
)

// Manager 工具管理器
type Manager struct {
	tools map[string]types.Tool
}

// NewManager 创建工具管理器
func NewManager(config *config.ToolsConfig) *Manager {
	m := &Manager{
		tools: make(map[string]types.Tool),
	}

	// 注册天气工具
	if config.Weather.Enabled {
		weatherTool := NewWeatherTool()
		m.tools[weatherTool.Name()] = weatherTool
	}

	// 注册计算器工具
	if config.Calculator.Enabled {
		calcTool := NewCalculatorTool()
		m.tools[calcTool.Name()] = calcTool
	}

	// 注册知识库工具
	// if config.Knowledge.Enabled {
	// 	knowledgeTool := NewKnowledgeTool()
	// 	m.tools[knowledgeTool.Name()] = knowledgeTool
	// }

	// 注册直连天下订单详情工具
	if config.ZltxOrderDetail.Enabled {
		zltxOrderDetailTool := NewZltxOrderDetailTool(config.ZltxOrderDetail)
		m.tools[zltxOrderDetailTool.Name()] = zltxOrderDetailTool
	}

	// 注册直连天下订单日志工具
	// if config.ZltxOrderLog.Enabled {
	// 	zltxOrderLogTool := NewZltxOrderLogTool(config.ZltxOrderLog)
	// 	m.tools[zltxOrderLogTool.Name()] = zltxOrderLogTool
	// }

	return m
}

// GetTool 获取工具
func (m *Manager) GetTool(name string) (types.Tool, bool) {
	tool, exists := m.tools[name]
	return tool, exists
}

// GetAllTools 获取所有工具
func (m *Manager) GetAllTools() []types.Tool {
	tools := make([]types.Tool, 0, len(m.tools))
	for _, tool := range m.tools {
		tools = append(tools, tool)
	}
	return tools
}

// GetToolDefinitions 获取所有工具定义
func (m *Manager) GetToolDefinitions(caller constants.Caller) []types.ToolDefinition {
	definitions := make([]types.ToolDefinition, 0, len(m.tools))
	for _, tool := range m.tools {
		definitions = append(definitions, tool.Definition())
	}

	return definitions
}

// ExecuteTool 执行工具
func (m *Manager) ExecuteTool(ctx context.Context, name string, args json.RawMessage) (interface{}, error) {
	tool, exists := m.GetTool(name)
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	return tool.Execute(ctx, args)
}

// ExecuteToolCalls 执行多个工具调用
func (m *Manager) ExecuteToolCalls(ctx context.Context, toolCalls []types.ToolCall) ([]types.ToolCall, error) {
	results := make([]types.ToolCall, len(toolCalls))

	for i, toolCall := range toolCalls {
		results[i] = toolCall

		// 执行工具
		result, err := m.ExecuteTool(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
		if err != nil {
			// 将错误信息作为结果返回
			errorResult := map[string]interface{}{
				"error": err.Error(),
			}
			resultBytes, _ := json.Marshal(errorResult)
			results[i].Result = resultBytes
		} else {
			// 将成功结果序列化
			resultBytes, err := json.Marshal(result)
			if err != nil {
				errorResult := map[string]interface{}{
					"error": fmt.Sprintf("failed to serialize result: %v", err),
				}
				resultBytes, _ = json.Marshal(errorResult)
			}
			results[i].Result = resultBytes
		}
	}

	return results, nil
}
