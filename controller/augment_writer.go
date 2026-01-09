package controller

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// SSE 格式前缀和结束标记
const (
	sseDataPrefix = "data: "
	sseDoneMarker = "[DONE]"
)

// NDJSONResponseWriter 将 SSE 格式转换为 NDJSON 格式的 ResponseWriter 包装器
// SSE 格式: "data: {...}\n" -> NDJSON 格式: "{...}\n"
type NDJSONResponseWriter struct {
	gin.ResponseWriter
	buffer bytes.Buffer // 缓存不完整的行
	mu     sync.Mutex
}

// NewNDJSONResponseWriter 创建新的 NDJSON 响应写入器
func NewNDJSONResponseWriter(w gin.ResponseWriter) *NDJSONResponseWriter {
	return &NDJSONResponseWriter{
		ResponseWriter: w,
	}
}

// Write 拦截写入数据，将 SSE 格式转换为 NDJSON 格式
func (w *NDJSONResponseWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 将新数据追加到缓冲区
	w.buffer.Write(data)

	// 处理缓冲区中的完整行
	return w.processBuffer()
}

// processBuffer 处理缓冲区中的数据，转换完整的行
func (w *NDJSONResponseWriter) processBuffer() (int, error) {
	content := w.buffer.String()

	// 查找最后一个换行符的位置
	lastNewline := strings.LastIndex(content, "\n")
	if lastNewline == -1 {
		// 没有完整的行，保留在缓冲区
		return len(content), nil
	}

	// 分离完整的行和不完整的部分
	completeLines := content[:lastNewline+1]
	remaining := content[lastNewline+1:]

	// 重置缓冲区，保留不完整的部分
	w.buffer.Reset()
	w.buffer.WriteString(remaining)

	// 处理每一行
	var output bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(completeLines))
	for scanner.Scan() {
		line := scanner.Text()
		converted := w.convertLine(line)
		if converted != "" {
			output.WriteString(converted)
			output.WriteString("\n")
		}
	}

	// 写入转换后的数据
	if output.Len() > 0 {
		_, err := w.ResponseWriter.Write(output.Bytes())
		if err != nil {
			return 0, err
		}
		// 刷新输出
		if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
			flusher.Flush()
		}
	}

	return len(content), nil
}

// OpenAI 流式响应结构（用于解析）
type openAIStreamChunk struct {
	ID      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	Model   string `json:"model,omitempty"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role             string `json:"role,omitempty"`
			Content          string `json:"content,omitempty"`
			ReasoningContent string `json:"reasoning_content,omitempty"`
			Reasoning        string `json:"reasoning,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason,omitempty"`
	} `json:"choices,omitempty"`
}

// Augment BackChatResult 响应结构
// 注意：所有字段都需要存在，因为 Augment 插件期望完整的结构
type augmentBackChatResult struct {
	Text                        string        `json:"text"`
	UnknownBlobNames            []string      `json:"unknown_blob_names"`
	CheckpointNotFound          bool          `json:"checkpoint_not_found"`
	WorkspaceFileChunks         []any         `json:"workspace_file_chunks"`
	IncorporatedExternalSources []any         `json:"incorporated_external_sources"`
	Nodes                       []augmentNode `json:"nodes"`
	StopReason                  *int          `json:"stop_reason"` // 使用数字类型：1=end_turn, 2=max_tokens, 3=tool_use
}

// augmentNode Augment 节点结构
type augmentNode struct {
	ID              int                  `json:"id"`
	Type            int                  `json:"type"` // 8=thinking, 2=content_end, 3=response_end
	Content         string               `json:"content"`
	ToolUse         any                  `json:"tool_use"`
	Thinking        *augmentThinking     `json:"thinking"`
	BillingMetadata any                  `json:"billing_metadata"`
	Metadata        *augmentNodeMetadata `json:"metadata"`
	TokenUsage      any                  `json:"token_usage"`
}

// augmentThinking Augment thinking 结构
type augmentThinking struct {
	Summary                   string `json:"summary"`
	Content                   any    `json:"content"`
	OpenAIResponsesAPIItemID  any    `json:"openai_responses_api_item_id"`
}

// augmentNodeMetadata Augment 节点元数据
type augmentNodeMetadata struct {
	OpenAIID string `json:"openai_id"`
	GoogleTS any    `json:"google_ts"`
	Provider any    `json:"provider"`
}

// convertLine 转换单行数据
// 1. 去除 "data: " 前缀
// 2. 过滤 [DONE] 标记和空行
// 3. 将 OpenAI 格式转换为 Augment BackChatResult 格式
func (w *NDJSONResponseWriter) convertLine(line string) string {
	// 去除首尾空白
	line = strings.TrimSpace(line)

	// 过滤空行
	if line == "" {
		return ""
	}

	// 非 data 行（如 event: 或 id: 等 SSE 字段），过滤掉
	if strings.HasPrefix(line, "event:") || strings.HasPrefix(line, "id:") || strings.HasPrefix(line, "retry:") {
		return ""
	}

	// 处理 SSE data 行
	jsonData := line
	if strings.HasPrefix(line, sseDataPrefix) {
		// 去除 "data: " 前缀
		jsonData = strings.TrimPrefix(line, sseDataPrefix)
	}

	// 过滤 [DONE] 标记
	if jsonData == sseDoneMarker {
		return ""
	}

	// 将 OpenAI 格式转换为 Augment BackChatResult 格式
	return w.convertToAugmentFormat(jsonData)
}

// convertToAugmentFormat 将 OpenAI 流式响应转换为 Augment BackChatResult 格式
// OpenAI: {"id":"...","choices":[{"delta":{"content":"Hello","reasoning_content":"..."}}]}
// Augment: {"text":"Hello","unknown_blob_names":[],"checkpoint_not_found":false,"workspace_file_chunks":[],"incorporated_external_sources":[],"nodes":[{"thinking":{"summary":"..."}}],"stop_reason":null}
func (w *NDJSONResponseWriter) convertToAugmentFormat(jsonData string) string {
	var chunk openAIStreamChunk
	if err := json.Unmarshal([]byte(jsonData), &chunk); err != nil {
		// 解析失败，返回原始数据
		return jsonData
	}

	// 构建 Augment 响应，初始化所有必需字段
	result := augmentBackChatResult{
		Text:                        "",
		UnknownBlobNames:            []string{},
		CheckpointNotFound:          false,
		WorkspaceFileChunks:         []any{},
		IncorporatedExternalSources: []any{},
		Nodes:                       []augmentNode{},
		StopReason:                  nil,
	}

	// 提取内容
	if len(chunk.Choices) > 0 {
		choice := chunk.Choices[0]
		result.Text = choice.Delta.Content

		// 处理 reasoning_content（来自 Claude thinking）
		reasoningContent := choice.Delta.ReasoningContent
		if reasoningContent == "" {
			reasoningContent = choice.Delta.Reasoning
		}

		if reasoningContent != "" {
			// 添加 thinking 节点
			result.Nodes = append(result.Nodes, augmentNode{
				ID:      0,
				Type:    8, // 8 = thinking 类型
				Content: "",
				ToolUse: nil,
				Thinking: &augmentThinking{
					Summary:                  reasoningContent,
					Content:                  nil,
					OpenAIResponsesAPIItemID: nil,
				},
				BillingMetadata: nil,
				Metadata: &augmentNodeMetadata{
					OpenAIID: "",
					GoogleTS: nil,
					Provider: nil,
				},
				TokenUsage: nil,
			})
		}

		// 处理结束原因
		if choice.FinishReason != nil && *choice.FinishReason != "" {
			stopReason := convertFinishReason(*choice.FinishReason)
			result.StopReason = &stopReason
		}
	}

	// 序列化为 JSON
	output, err := json.Marshal(result)
	if err != nil {
		return jsonData
	}

	return string(output)
}

// convertFinishReason 将 OpenAI finish_reason 转换为 Augment stop_reason (数字类型)
// Augment stop_reason: 1=end_turn, 2=max_tokens, 3=tool_use
func convertFinishReason(reason string) int {
	switch reason {
	case "stop", "end_turn":
		return 1 // end_turn
	case "length", "max_tokens":
		return 2 // max_tokens
	case "tool_calls", "function_call", "tool_use":
		return 3 // tool_use
	default:
		return 1 // 默认返回 end_turn
	}
}

// Flush 刷新剩余缓冲区数据
func (w *NDJSONResponseWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 处理缓冲区中剩余的数据
	if w.buffer.Len() > 0 {
		line := w.convertLine(w.buffer.String())
		if line != "" {
			w.ResponseWriter.Write([]byte(line + "\n"))
		}
		w.buffer.Reset()
	}

	// 调用底层的 Flush
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack 实现 http.Hijacker 接口
func (w *NDJSONResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// WrapContextForNDJSON 包装 gin.Context 使用 NDJSON 响应写入器
// 返回清理函数，应在请求结束时调用以刷新剩余数据
func WrapContextForNDJSON(c *gin.Context) func() {
	originalWriter := c.Writer
	ndjsonWriter := NewNDJSONResponseWriter(originalWriter)
	c.Writer = ndjsonWriter

	// 返回清理函数
	return func() {
		ndjsonWriter.Flush()
		c.Writer = originalWriter
	}
}
