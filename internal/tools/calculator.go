package tools

import (
	"ai_scheduler/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"math"
)

// CalculatorTool 计算器工具
type CalculatorTool struct{}

// NewCalculatorTool 创建计算器工具
func NewCalculatorTool() *CalculatorTool {
	return &CalculatorTool{}
}

// Name 返回工具名称
func (c *CalculatorTool) Name() string {
	return "calculate"
}

// Description 返回工具描述
func (c *CalculatorTool) Description() string {
	return "执行基本的数学运算，支持加减乘除和幂运算"
}

// Definition 返回工具定义
func (c *CalculatorTool) Definition() types.ToolDefinition {
	return types.ToolDefinition{
		Type: "function",
		Function: types.FunctionDef{
			Name:        c.Name(),
			Description: c.Description(),
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"operation": map[string]interface{}{
						"type":        "string",
						"description": "运算类型",
						"enum":        []string{"add", "subtract", "multiply", "divide", "power"},
					},
					"a": map[string]interface{}{
						"type":        "number",
						"description": "第一个数字",
					},
					"b": map[string]interface{}{
						"type":        "number",
						"description": "第二个数字",
					},
				},
				"required": []string{"operation", "a", "b"},
			},
		},
	}
}

// CalculateRequest 计算请求参数
type CalculateRequest struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

// CalculateResponse 计算响应
type CalculateResponse struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
	Result    float64 `json:"result"`
	Expression string `json:"expression"`
}

// Execute 执行计算
func (c *CalculatorTool) Execute(ctx context.Context, args json.RawMessage) (interface{}, error) {
	var req CalculateRequest
	if err := json.Unmarshal(args, &req); err != nil {
		return nil, fmt.Errorf("invalid calculate request: %w", err)
	}

	var result float64
	var expression string

	switch req.Operation {
	case "add":
		result = req.A + req.B
		expression = fmt.Sprintf("%.2f + %.2f = %.2f", req.A, req.B, result)
	case "subtract":
		result = req.A - req.B
		expression = fmt.Sprintf("%.2f - %.2f = %.2f", req.A, req.B, result)
	case "multiply":
		result = req.A * req.B
		expression = fmt.Sprintf("%.2f × %.2f = %.2f", req.A, req.B, result)
	case "divide":
		if req.B == 0 {
			return nil, fmt.Errorf("division by zero is not allowed")
		}
		result = req.A / req.B
		expression = fmt.Sprintf("%.2f ÷ %.2f = %.2f", req.A, req.B, result)
	case "power":
		result = math.Pow(req.A, req.B)
		expression = fmt.Sprintf("%.2f ^ %.2f = %.2f", req.A, req.B, result)
	default:
		return nil, fmt.Errorf("unsupported operation: %s", req.Operation)
	}

	// 检查结果是否有效
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return nil, fmt.Errorf("calculation resulted in invalid number")
	}

	return &CalculateResponse{
		Operation:  req.Operation,
		A:          req.A,
		B:          req.B,
		Result:     result,
		Expression: expression,
	}, nil
}