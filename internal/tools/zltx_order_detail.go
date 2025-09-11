package tools

import (
	"ai_scheduler/internal/config"
	"ai_scheduler/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ZltxOrderDetailTool 直连天下订单详情工具
type ZltxOrderDetailTool struct {
	config config.ToolConfig
}

// NewZltxOrderDetailTool 创建直连天下订单详情工具
func NewZltxOrderDetailTool(config config.ToolConfig) *ZltxOrderDetailTool {
	return &ZltxOrderDetailTool{config: config}
}

// Name 返回工具名称
func (w *ZltxOrderDetailTool) Name() string {
	return "zltxOrderDetail"
}

// Description 返回工具描述
func (w *ZltxOrderDetailTool) Description() string {
	return "获取直连天下订单详情"
}

// Definition 返回工具定义
func (w *ZltxOrderDetailTool) Definition() types.ToolDefinition {
	return types.ToolDefinition{
		Type: "function",
		Function: types.FunctionDef{
			Name:        w.Name(),
			Description: w.Description(),
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"number": map[string]interface{}{
						"type":        "string",
						"description": "订单编号/流水号",
					},
				},
				"required": []string{"number"},
			},
		},
	}
}

// ZltxOrderDetailRequest 直连天下订单详情请求参数
type ZltxOrderDetailRequest struct {
	Number string `json:"number"`
}

// ZltxOrderDetailResponse 直连天下订单详情响应
type ZltxOrderDetailResponse struct {
	Code  int                 `json:"code"`
	Error string              `json:"error"`
	Data  ZltxOrderDetailData `json:"data"`
}

// ZltxOrderDetailData 直连天下订单详情数据
type ZltxOrderDetailData struct {
	Direct map[string]any `json:"direct"`
	Order  map[string]any `json:"order"`
}

// Execute 执行直连天下订单详情查询
func (w *ZltxOrderDetailTool) Execute(ctx context.Context, args json.RawMessage) (interface{}, error) {
	var req ZltxOrderDetailRequest
	if err := json.Unmarshal(args, &req); err != nil {
		return nil, fmt.Errorf("invalid zltxOrderDetail request: %w", err)
	}

	if req.Number == "" {
		return nil, fmt.Errorf("number is required")
	}

	// 这里可以集成真实的直连天下订单详情API
	return w.getZltxOrderDetail(ctx, req.Number), nil
}

// getMockZltxOrderDetail 获取模拟直连天下订单详情数据
func (w *ZltxOrderDetailTool) getZltxOrderDetail(ctx context.Context, number string) *ZltxOrderDetailResponse {
	url := fmt.Sprintf("%s/admin/direct/ai/%s", w.config.BaseURL, number)
	authorization := fmt.Sprintf("Bearer %s", w.config.APIKey)

	// 发送http请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &ZltxOrderDetailResponse{}
	}
	req.Header.Set("Authorization", authorization)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &ZltxOrderDetailResponse{}
	}
	defer resp.Body.Close()

	return &ZltxOrderDetailResponse{}
}
