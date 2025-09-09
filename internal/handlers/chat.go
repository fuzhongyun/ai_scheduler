package handlers

import (
	"ai_scheduler/pkg/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ChatHandler 聊天处理器
type ChatHandler struct {
	routerService types.RouterService
}

// NewChatHandler 创建聊天处理器
func NewChatHandler(routerService types.RouterService) *ChatHandler {
	return &ChatHandler{
		routerService: routerService,
	}
}

// ChatRequest HTTP聊天请求
type ChatRequest struct {
	Message string `json:"message" binding:"required" example:"北京今天天气怎么样？"`
}

// ChatResponse HTTP聊天响应
type ChatResponse struct {
	Message   string             `json:"message" example:"北京今天天气晴朗，温度15.3°C"`
	ToolCalls []ToolCallResponse `json:"tool_calls,omitempty"`
	Finished  bool               `json:"finished" example:"true"`
}

// ToolCallResponse 工具调用响应
type ToolCallResponse struct {
	ID       string               `json:"id" example:"call_1"`
	Type     string               `json:"type" example:"function"`
	Function FunctionCallResponse `json:"function"`
	Result   interface{}          `json:"result,omitempty"`
}

// FunctionCallResponse 函数调用响应
type FunctionCallResponse struct {
	Name      string      `json:"name" example:"get_weather"`
	Arguments interface{} `json:"arguments"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"请求参数错误"`
}

// Chat 处理聊天请求
// @Summary 智能聊天
// @Description 发送消息给AI助手，支持工具调用
// @Tags chat
// @Accept json
// @Produce json
// @Param request body ChatRequest true "聊天请求"
// @Success 200 {object} ChatResponse "聊天响应"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 500 {object} ErrorResponse "服务器错误"
// @Router /api/v1/chat [post]
func (h *ChatHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// 转换为服务层请求
	serviceReq := &types.ChatRequest{
		Message: req.Message,
	}

	// 调用路由服务
	response, err := h.routerService.Route(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Service error",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	// 转换响应格式
	httpResponse := &ChatResponse{
		Message:  response.Message,
		Finished: response.Finished,
	}

	// 转换工具调用
	if len(response.ToolCalls) > 0 {
		httpResponse.ToolCalls = make([]ToolCallResponse, len(response.ToolCalls))
		for i, toolCall := range response.ToolCalls {
			httpResponse.ToolCalls[i] = ToolCallResponse{
				ID:   toolCall.ID,
				Type: toolCall.Type,
				Function: FunctionCallResponse{
					Name:      toolCall.Function.Name,
					Arguments: toolCall.Function.Arguments,
				},
				Result: toolCall.Result,
			}
		}
	}

	c.JSON(http.StatusOK, httpResponse)
}

// Health 健康检查
// @Summary 健康检查
// @Description 检查服务是否正常运行
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *ChatHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "ai-scheduler",
	})
}
