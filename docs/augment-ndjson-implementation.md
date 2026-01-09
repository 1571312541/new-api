# Augment 插件响应格式适配方案

> 本文档描述如何修改 new-api 以支持 Augment VSCode 插件所需的响应格式。

## 问题背景

### 问题 1：SSE vs NDJSON 格式

Augment 插件期望的流式响应格式为 **NDJSON**（Newline Delimited JSON），而 new-api 当前返回的是 **SSE**（Server-Sent Events）格式。

| 项目 | SSE 格式 (当前) | NDJSON 格式 (期望) |
|------|----------------|-------------------|
| 数据前缀 | `data: {...}\n` | `{...}\n` |
| 结束标记 | `data: [DONE]\n` | 无 |
| 插件解析 | 需要去除 `data: ` 前缀 | 直接 `JSON.parse()` |

**错误表现**：`JSON parse failed for data: {...}`

### 问题 2：响应结构不匹配

Augment 插件期望 **BackChatResult** 格式，而 new-api 返回的是 **OpenAI Chat Completions** 格式。

| 项目 | OpenAI 格式 (当前) | Augment 格式 (期望) |
|------|-------------------|-------------------|
| 内容字段 | `choices[0].delta.content` | `text` |
| 结束标记 | `finish_reason: "stop"` | `stop_reason: "end_turn"` |

**错误表现**：`Value of BackChatResult.text has unexpected type. Expected string, received undefined`

---

## 实现方案：ResponseWriter 包装器 + 响应结构转换

### 设计思路

使用 **响应流转换器** 方案，在 Controller 层拦截 SSE 输出并转换为 Augment 格式，**完全不修改上游核心代码**。

```
请求 → AugmentController → [包装 Writer] → 标准 Relay → SSE 输出
                                ↓
                          [转换器拦截]
                                ↓
                     1. 去除 SSE 前缀 (data: )
                     2. 转换 JSON 结构
                                ↓
                          Augment 输出 → 插件
```

### 转换流程

```
OpenAI 流式响应:
data: {"id":"...","choices":[{"delta":{"content":"Hello"}}]}

     ↓ 去除 SSE 前缀

{"id":"...","choices":[{"delta":{"content":"Hello"}}]}

     ↓ 转换 JSON 结构

{"text":"Hello"}
```

### 优势

| 对比项 | 修改渠道代码方案 | ResponseWriter 包装器方案 |
|--------|------------------|--------------------------|
| 修改文件数 | 6 个（含核心代码） | 2 个（仅 Augment 相关） |
| 上游合并风险 | 高 | 极低 |
| 代码侵入性 | 高 | 低 |
| 维护成本 | 高 | 低 |

---

## 实现代码

### 1. 新增文件：`controller/augment_writer.go`

创建 SSE 到 Augment 格式的流转换器：

```go
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

// OpenAI 流式响应结构（用于解析）
type openAIStreamChunk struct {
	ID      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	Model   string `json:"model,omitempty"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason,omitempty"`
	} `json:"choices,omitempty"`
}

// Augment BackChatResult 响应结构
type augmentBackChatResult struct {
	Text       string `json:"text"`
	StopReason string `json:"stop_reason,omitempty"`
}

// NDJSONResponseWriter 将 SSE + OpenAI 格式转换为 Augment BackChatResult 格式
type NDJSONResponseWriter struct {
	gin.ResponseWriter
	buffer bytes.Buffer
	mu     sync.Mutex
}

// convertToAugmentFormat 将 OpenAI 流式响应转换为 Augment BackChatResult 格式
func (w *NDJSONResponseWriter) convertToAugmentFormat(jsonData string) string {
	var chunk openAIStreamChunk
	if err := json.Unmarshal([]byte(jsonData), &chunk); err != nil {
		return jsonData // 解析失败，返回原始数据
	}

	result := augmentBackChatResult{}
	if len(chunk.Choices) > 0 {
		choice := chunk.Choices[0]
		result.Text = choice.Delta.Content
		if choice.FinishReason != nil && *choice.FinishReason != "" {
			result.StopReason = convertFinishReason(*choice.FinishReason)
		}
	}

	output, err := json.Marshal(result)
	if err != nil {
		return jsonData
	}
	return string(output)
}

// convertFinishReason 将 OpenAI finish_reason 转换为 Augment stop_reason
func convertFinishReason(reason string) string {
	switch reason {
	case "stop":
		return "end_turn"
	case "length":
		return "max_tokens"
	case "tool_calls", "function_call":
		return "tool_use"
	default:
		return reason
	}
}
```

---

### 2. 修改文件：`controller/augment.go`

在 `AugmentChatStream` 函数中使用转换器：

```go
// AugmentChatStream 处理 Augment 插件的 chat-stream 请求
// 1. 解析 Augment 格式请求
// 2. 转换为 OpenAI Chat Completions 格式
// 3. 使用 NDJSON 响应包装器
// 4. 替换请求体后调用 Relay 处理
func AugmentChatStream(c *gin.Context) {
	// ... 请求解析和转换逻辑 ...

	// 替换请求体
	c.Request.Body = io.NopCloser(bytes.NewBuffer(newBody))
	c.Request.ContentLength = int64(len(newBody))

	// 更新缓存的请求体
	c.Set(common.KeyRequestBody, newBody)

	// 包装 ResponseWriter 以将 SSE 格式转换为 NDJSON 格式
	// Augment 插件期望 NDJSON（每行纯 JSON），而非 SSE（data: {...}）
	cleanup := WrapContextForNDJSON(c)
	defer cleanup()

	// 使用标准 OpenAI 格式调用 Relay，响应会被转换为 NDJSON
	Relay(c, types.RelayFormatOpenAI)
}
```

---

## 修改文件清单

| 文件路径 | 修改类型 | 说明 |
|----------|----------|------|
| `controller/augment_writer.go` | **新增** | SSE 到 NDJSON 流转换器 |
| `controller/augment.go` | 修改 | 使用 NDJSON 响应包装器 |

**注意**：此方案完全不修改任何上游核心代码（relay、types 等），最大程度降低合并冲突风险。

---

## 转换逻辑说明

### SSE 到 NDJSON

| 输入 (SSE) | 输出 (NDJSON) | 说明 |
|------------|---------------|------|
| `data: {"id":"..."}` | `{"text":"..."}` | 去除前缀并转换结构 |
| `data: [DONE]` | *(过滤)* | 不输出结束标记 |
| `event: message` | *(过滤)* | 过滤 SSE 事件字段 |
| *(空行)* | *(过滤)* | 过滤空行 |

### OpenAI 到 Augment 结构

| OpenAI 字段 | Augment 字段 | 说明 |
|-------------|--------------|------|
| `choices[0].delta.content` | `text` | 文本内容 |
| `finish_reason: "stop"` | `stop_reason: "end_turn"` | 正常结束 |
| `finish_reason: "length"` | `stop_reason: "max_tokens"` | 达到最大长度 |
| `finish_reason: "tool_calls"` | `stop_reason: "tool_use"` | 工具调用 |

---

## 测试验证

### 测试步骤

1. 构建并启动 new-api 服务
2. 使用 Augment 插件发送聊天请求
3. 检查响应格式

### 预期结果

- 响应每行为 Augment BackChatResult 格式 JSON
- 无 `data: ` SSE 前缀
- 流式结束时不发送 `[DONE]` 标记
- 插件能正常解析响应并显示内容

### 调试方法

```bash
curl -N -X POST http://localhost:3000/augment/v1/chat-stream \
  -H "Authorization: Bearer sk-xxx" \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello", "model": "gpt-4o"}'
```

观察输出是否为：
```
{"text":"Hello"}
{"text":" there"}
{"text":"!","stop_reason":"end_turn"}
```

而非：
```
data: {"id":"...","choices":[{"delta":{"content":"Hello"}}]}
data: {"id":"...","choices":[{"delta":{"content":" there"}}]}
data: [DONE]
```

---

## 注意事项

1. **向后兼容**：此修改仅影响 `/augment/v1/chat-stream` 端点，不影响其他 API
2. **上游合并**：所有修改都在 `controller/augment*.go` 文件中，与上游代码无冲突
3. **性能影响**：转换器使用缓冲和流式处理，性能开销极小
4. **线程安全**：使用 mutex 保护缓冲区，支持并发写入
