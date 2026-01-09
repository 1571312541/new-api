package controller

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

// AugmentChatStream 处理 Augment 插件的 chat-stream 请求
// 1. 解析 Augment 格式请求
// 2. 转换为 OpenAI Chat Completions 格式
// 3. 使用 NDJSON 响应包装器
// 4. 替换请求体后调用 Relay 处理
func AugmentChatStream(c *gin.Context) {
	// 解析 Augment 请求
	var augmentReq dto.AugmentChatStreamRequest
	if err := c.ShouldBindJSON(&augmentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": types.OpenAIError{
				Message: fmt.Sprintf("invalid request: %s", err.Error()),
				Type:    "invalid_request_error",
				Code:    "invalid_request",
			},
		})
		return
	}

	// 转换为 OpenAI 请求
	openAIReq := convertAugmentToOpenAI(&augmentReq)

	// 序列化新请求体并替换
	newBody, err := common.Marshal(openAIReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": types.OpenAIError{
				Message: fmt.Sprintf("failed to marshal request: %s", err.Error()),
				Type:    "internal_error",
				Code:    "marshal_failed",
			},
		})
		return
	}

	// 替换请求体
	c.Request.Body = io.NopCloser(bytes.NewBuffer(newBody))
	c.Request.ContentLength = int64(len(newBody))

	// 更新缓存的请求体，确保后续 GetRequestBody 返回新的请求体
	c.Set(common.KeyRequestBody, newBody)

	// 包装 ResponseWriter 以将 SSE 格式转换为 NDJSON 格式
	// Augment 插件期望 NDJSON（每行纯 JSON），而非 SSE（data: {...}）
	cleanup := WrapContextForNDJSON(c)
	defer cleanup()

	// 使用标准 OpenAI 格式调用 Relay，响应会被转换为 NDJSON
	Relay(c, types.RelayFormatOpenAI)
}

// convertAugmentToOpenAI 将 Augment 请求转换为 OpenAI Chat Completions 格式
func convertAugmentToOpenAI(req *dto.AugmentChatStreamRequest) *dto.GeneralOpenAIRequest {
	messages := make([]dto.Message, 0)

	// 1. 转换历史消息
	for _, h := range req.ChatHistory {
		messages = append(messages, dto.Message{
			Role:    h.Role,
			Content: h.Content,
		})
	}

	// 2. 构建当前用户消息（包含代码上下文）
	userContent := buildAugmentUserContent(req)
	messages = append(messages, dto.Message{
		Role:    "user",
		Content: userContent,
	})

	// 3. 确定模型名称
	model := resolveAugmentModel(req)

	// 4. 转换工具定义（如果有）
	var tools []dto.ToolCallRequest
	if len(req.ToolDefinitions) > 0 {
		tools = convertAugmentTools(req.ToolDefinitions)
	}

	return &dto.GeneralOpenAIRequest{
		Model:    model,
		Messages: messages,
		Stream:   true, // Augment 始终使用流式响应
		Tools:    tools,
	}
}

// buildAugmentUserContent 构建用户消息内容（包含代码上下文）
func buildAugmentUserContent(req *dto.AugmentChatStreamRequest) string {
	var sb strings.Builder

	// 添加用户消息
	sb.WriteString(req.Message)

	// 如果有选中的代码，添加代码上下文
	if req.SelectedCode != "" {
		sb.WriteString("\n\n--- Selected Code ---\n")
		sb.WriteString(req.SelectedCode)
	}

	// 如果有前缀上下文
	if req.Prefix != "" {
		sb.WriteString("\n\n--- Code Before ---\n")
		sb.WriteString(req.Prefix)
	}

	// 如果有后缀上下文
	if req.Suffix != "" {
		sb.WriteString("\n\n--- Code After ---\n")
		sb.WriteString(req.Suffix)
	}

	// 如果有文件路径
	if req.Path != "" {
		sb.WriteString("\n\n--- File Path ---\n")
		sb.WriteString(req.Path)
	}

	// 如果有语言信息
	if req.Lang != "" {
		sb.WriteString("\n\n--- Language ---\n")
		sb.WriteString(req.Lang)
	}

	// 如果有代码块（blobs）
	if len(req.Blobs) > 0 {
		sb.WriteString("\n\n--- Related Files ---\n")
		for _, blob := range req.Blobs {
			if blob.Path != "" {
				sb.WriteString(fmt.Sprintf("\n// %s\n", blob.Path))
			}
			sb.WriteString(blob.Content)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// resolveAugmentModel 解析目标模型
// 优先使用 third_party_override 中的模型，否则使用请求中的 model 字段
func resolveAugmentModel(req *dto.AugmentChatStreamRequest) string {
	// 尝试解析 third_party_override
	if req.ThirdPartyOverride != "" {
		override, err := decodeThirdPartyOverride(req.ThirdPartyOverride)
		if err == nil && override.ProviderModelName != "" {
			return override.ProviderModelName
		}
	}

	// 使用请求中的模型
	if req.Model != "" {
		return req.Model
	}

	// 默认模型
	return "gpt-4o"
}

// decodeThirdPartyOverride 解码 Base64 编码的第三方覆盖配置
func decodeThirdPartyOverride(encoded string) (*dto.AugmentThirdPartyOverride, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		// 尝试 URL-safe Base64
		decoded, err = base64.URLEncoding.DecodeString(encoded)
		if err != nil {
			return nil, err
		}
	}

	var override dto.AugmentThirdPartyOverride
	if err := common.Unmarshal(decoded, &override); err != nil {
		return nil, err
	}

	return &override, nil
}

// convertAugmentTools 转换 Augment 工具定义为 OpenAI 格式
func convertAugmentTools(tools []dto.AugmentToolDefinition) []dto.ToolCallRequest {
	result := make([]dto.ToolCallRequest, 0, len(tools))
	for _, tool := range tools {
		result = append(result, dto.ToolCallRequest{
			Type: "function",
			Function: dto.FunctionRequest{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
			},
		})
	}
	return result
}
