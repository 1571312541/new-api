[根目录](../CLAUDE.md) > **dto**

# DTO 模块 - 数据传输对象

> 定义 API 请求和响应的数据结构。

## 模块职责

DTO (Data Transfer Object) 模块定义了：

1. **请求结构**：各 API 的请求体格式
2. **响应结构**：各 API 的响应体格式
3. **通用结构**：共享的数据类型
4. **错误结构**：错误响应格式

## 关键文件

### OpenAI 格式

| 文件 | 说明 |
|------|------|
| `openai_request.go` | OpenAI 请求结构 |
| `openai_response.go` | OpenAI 响应结构 |
| `openai_image.go` | 图像 API 结构 |
| `openai_video.go` | 视频 API 结构 |

### Claude 格式

| 文件 | 说明 |
|------|------|
| `claude.go` | Claude Messages API 结构 |

### Gemini 格式

| 文件 | 说明 |
|------|------|
| `gemini.go` | Google Gemini API 结构 |

### 其他格式

| 文件 | 说明 |
|------|------|
| `audio.go` | 音频 API 结构 |
| `embedding.go` | 嵌入 API 结构 |
| `rerank.go` | 重排序 API 结构 |
| `realtime.go` | 实时 API 结构 |
| `midjourney.go` | Midjourney 结构 |
| `suno.go` | Suno 音乐结构 |
| `video.go` | 视频生成结构 |
| `task.go` | 异步任务结构 |

### 通用结构

| 文件 | 说明 |
|------|------|
| `request_common.go` | 通用请求结构 |
| `error.go` | 错误响应结构 |
| `pricing.go` | 定价结构 |
| `notify.go` | 通知结构 |

## 核心结构

### GeneralOpenAIRequest

```go
type GeneralOpenAIRequest struct {
    Model            string          `json:"model"`
    Messages         []Message       `json:"messages"`
    Stream           bool            `json:"stream,omitempty"`
    Temperature      *float64        `json:"temperature,omitempty"`
    MaxTokens        int             `json:"max_tokens,omitempty"`
    TopP             *float64        `json:"top_p,omitempty"`
    FrequencyPenalty *float64        `json:"frequency_penalty,omitempty"`
    PresencePenalty  *float64        `json:"presence_penalty,omitempty"`
    Tools            []Tool          `json:"tools,omitempty"`
    // ...
}
```

### Message

```go
type Message struct {
    Role    string      `json:"role"`
    Content interface{} `json:"content"`
    Name    string      `json:"name,omitempty"`
}
```

### ClaudeRequest

```go
type ClaudeRequest struct {
    Model         string         `json:"model"`
    Messages      []ClaudeMessage `json:"messages"`
    MaxTokens     int            `json:"max_tokens"`
    System        interface{}    `json:"system,omitempty"`
    Stream        bool           `json:"stream,omitempty"`
    Temperature   *float64       `json:"temperature,omitempty"`
    // ...
}
```

### ImageRequest

```go
type ImageRequest struct {
    Model          string `json:"model"`
    Prompt         string `json:"prompt"`
    N              int    `json:"n,omitempty"`
    Size           string `json:"size,omitempty"`
    Quality        string `json:"quality,omitempty"`
    ResponseFormat string `json:"response_format,omitempty"`
}
```

## 常见问题 (FAQ)

**Q: 如何添加新的请求字段？**
A: 在对应的结构体中添加字段，注意 JSON 标签

**Q: 如何处理可选字段？**
A: 使用指针类型或 `omitempty` 标签

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28T20:58:40+0800 | 初始化 | 首次生成模块文档 |
