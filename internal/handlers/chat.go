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
	UserInput string `json:"user_input" binding:"required" example:"考勤规则"`
	Caller    string `json:"caller" binding:"required" example:"zltx"`
	SessionID string `json:"session_id" example:"default"`
}

// ChatResponse HTTP聊天响应
type ChatResponse struct {
	Status   string `json:"status" example:"success"` // 处理状态
	Message  string `json:"message" example:""`       // 响应消息
	Data     any    `json:"data,omitempty"`           // 响应数据
	TaskCode string `json:"task_code,omitempty"`      // 任务代码
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

func (h *ChatHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ChatResponse{
			Status:  "error",
			Message: "请求参数错误",
		})
		return
	}

	// 转换为服务层请求
	serviceReq := &types.ChatRequest{
		UserInput: req.UserInput,
		Caller:    req.Caller,
		SessionID: req.SessionID,
		ChatRequestMeta: types.ChatRequestMeta{
			Authorization: c.GetHeader("Authorization"),
		},
	}

	// 调用路由服务
	response, err := h.routerService.Route(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ChatResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// 转换响应格式
	httpResponse := &ChatResponse{
		Message:  response.Message,
		Status:   response.Status,
		Data:     response.Data,
		TaskCode: response.TaskCode,
	}

	c.JSON(http.StatusOK, httpResponse)
}
