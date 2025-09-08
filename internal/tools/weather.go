package tools

import (
	"ai_scheduler/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// WeatherTool 天气查询工具
type WeatherTool struct {
	mockData bool
}

// NewWeatherTool 创建天气工具
func NewWeatherTool(mockData bool) *WeatherTool {
	return &WeatherTool{
		mockData: mockData,
	}
}

// Name 返回工具名称
func (w *WeatherTool) Name() string {
	return "get_weather"
}

// Description 返回工具描述
func (w *WeatherTool) Description() string {
	return "获取指定城市的天气信息"
}

// Definition 返回工具定义
func (w *WeatherTool) Definition() types.ToolDefinition {
	return types.ToolDefinition{
		Type: "function",
		Function: types.FunctionDef{
			Name:        w.Name(),
			Description: w.Description(),
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"city": map[string]interface{}{
						"type":        "string",
						"description": "城市名称，如：北京、上海、广州",
					},
					"unit": map[string]interface{}{
						"type":        "string",
						"description": "温度单位，celsius(摄氏度)或fahrenheit(华氏度)",
						"enum":        []string{"celsius", "fahrenheit"},
						"default":     "celsius",
					},
				},
				"required": []string{"city"},
			},
		},
	}
}

// WeatherRequest 天气请求参数
type WeatherRequest struct {
	City string `json:"city"`
	Unit string `json:"unit,omitempty"`
}

// WeatherResponse 天气响应
type WeatherResponse struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Unit        string  `json:"unit"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
	Timestamp   string  `json:"timestamp"`
}

// Execute 执行天气查询
func (w *WeatherTool) Execute(ctx context.Context, args json.RawMessage) (interface{}, error) {
	var req WeatherRequest
	if err := json.Unmarshal(args, &req); err != nil {
		return nil, fmt.Errorf("invalid weather request: %w", err)
	}

	if req.City == "" {
		return nil, fmt.Errorf("city is required")
	}

	if req.Unit == "" {
		req.Unit = "celsius"
	}

	if w.mockData {
		return w.getMockWeather(req.City, req.Unit), nil
	}

	// 这里可以集成真实的天气API
	return w.getMockWeather(req.City, req.Unit), nil
}

// getMockWeather 获取模拟天气数据
func (w *WeatherTool) getMockWeather(city, unit string) *WeatherResponse {
	rand.Seed(time.Now().UnixNano())

	// 模拟不同城市的基础温度
	baseTemp := map[string]float64{
		"北京": 15.0,
		"上海": 18.0,
		"广州": 25.0,
		"深圳": 26.0,
		"杭州": 17.0,
		"成都": 16.0,
	}

	temp := baseTemp[city]
	if temp == 0 {
		temp = 20.0 // 默认温度
	}

	// 添加随机变化
	temp += (rand.Float64() - 0.5) * 10

	// 转换温度单位
	if unit == "fahrenheit" {
		temp = temp*9/5 + 32
	}

	conditions := []string{"晴朗", "多云", "阴天", "小雨", "中雨"}
	condition := conditions[rand.Intn(len(conditions))]

	return &WeatherResponse{
		City:        city,
		Temperature: float64(int(temp*10)) / 10, // 保留一位小数
		Unit:        unit,
		Condition:   condition,
		Humidity:    rand.Intn(40) + 40, // 40-80%
		WindSpeed:   float64(rand.Intn(20)) + 1.0,
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
	}
}