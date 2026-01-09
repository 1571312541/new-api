[根目录](../../../CLAUDE.md) > [relay](../../CLAUDE.md) > [channel](../CLAUDE.md) > gemini

# Gemini 渠道适配器

> Google Gemini API 的完整适配实现，支持聊天、嵌入、图像生成及原生 API 透传。

## 模块概述

本模块实现了 Google Gemini API 的适配器，主要功能包括：

- OpenAI 格式到 Gemini 格式的请求转换
- Gemini 响应到 OpenAI 格式的转换
- 流式 SSE 响应处理
- Gemini 特有功能支持（Thinking、Grounding、Code Execution）
- Imagen 图像生成支持
- 嵌入向量 API 支持

## 文件清单与职责

| 文件 | 行数 | 职责 |
|------|------|------|
| `adaptor.go` | 286 | **入口文件** - 实现 Adaptor 接口，处理 URL 构建、请求头设置、请求转换、响应路由 |
| `constant.go` | 39 | 常量定义：支持的模型列表、安全设置类别、渠道名称 |
| `relay-gemini.go` | 1365 | **核心逻辑** - OpenAI↔Gemini 格式转换、思考适配器、流式/非流式响应处理 |
| `relay-gemini-native.go` | 106 | 原生 Gemini 格式处理：直接透传 Gemini 请求，保留原始格式 |

## 支持的模型

```go
var ModelList = []string{
    // 稳定版
    "gemini-1.5-pro", "gemini-1.5-flash", "gemini-1.5-flash-8b",
    "gemini-2.0-flash",
    // 最新版
    "gemini-1.5-pro-latest", "gemini-1.5-flash-latest",
    // 预览版
    "gemini-2.0-flash-lite-preview", "gemini-3-pro-preview",
    // 实验版
    "gemini-exp-1206", "gemini-2.0-flash-exp", "gemini-2.0-pro-exp",
    // 思考模型
    "gemini-2.0-flash-thinking-exp", "gemini-2.5-pro-exp-03-25",
    // Imagen 图像生成
    "imagen-3.0-generate-002",
    // 嵌入模型
    "gemini-embedding-exp-03-07", "text-embedding-004", "embedding-001",
}
```

## 核心类型定义

### 请求结构 (dto/gemini.go)

```go
type GeminiChatRequest struct {
    Contents           []GeminiChatContent        `json:"contents"`
    SafetySettings     []GeminiChatSafetySettings `json:"safetySettings,omitempty"`
    GenerationConfig   GeminiChatGenerationConfig `json:"generationConfig,omitempty"`
    Tools              json.RawMessage            `json:"tools,omitempty"`
    SystemInstructions *GeminiChatContent         `json:"systemInstruction,omitempty"`
}

type GeminiThinkingConfig struct {
    IncludeThoughts bool   `json:"includeThoughts,omitempty"`
    ThinkingBudget  *int   `json:"thinkingBudget,omitempty"`
    ThinkingLevel   string `json:"thinkingLevel,omitempty"`
}

type GeminiChatTool struct {
    GoogleSearch          any `json:"googleSearch,omitempty"`
    CodeExecution         any `json:"codeExecution,omitempty"`
    FunctionDeclarations  any `json:"functionDeclarations,omitempty"`
    URLContext            any `json:"urlContext,omitempty"`
}
```

### 响应结构

```go
type GeminiChatResponse struct {
    Candidates     []GeminiChatCandidate     `json:"candidates"`
    UsageMetadata  GeminiUsageMetadata       `json:"usageMetadata"`
}

type GeminiUsageMetadata struct {
    PromptTokenCount     int `json:"promptTokenCount"`
    CandidatesTokenCount int `json:"candidatesTokenCount"`
    TotalTokenCount      int `json:"totalTokenCount"`
    ThoughtsTokenCount   int `json:"thoughtsTokenCount"`
}
```

## 关键函数签名

### adaptor.go

```go
func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error)
func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error
func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)
func (a *Adaptor) ConvertGeminiRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeminiChatRequest) (any, error)
func (a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error)
func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error)
func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)
```

### relay-gemini.go

```go
func CovertOpenAI2Gemini(c *gin.Context, textRequest dto.GeneralOpenAIRequest, info *relaycommon.RelayInfo) (*dto.GeminiChatRequest, error)
func ThinkingAdaptor(geminiRequest *dto.GeminiChatRequest, info *relaycommon.RelayInfo, oaiRequest ...dto.GeneralOpenAIRequest)
func GeminiChatStreamHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)
func GeminiChatHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)
func GeminiEmbeddingHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)
func GeminiImageHandler(c *gin.Context, info *relaycommon.RelayInfo, resp *http.Response) (*dto.Usage, *types.NewAPIError)
```

## 请求处理流程

```mermaid
flowchart TD
    A[收到请求] --> B{请求格式?}

    B -->|OpenAI 格式| C[ConvertOpenAIRequest]
    B -->|Gemini 原生| D[ConvertGeminiRequest]
    B -->|Claude 格式| E[ConvertClaudeRequest]
    B -->|Embedding| F[ConvertEmbeddingRequest]
    B -->|Image| G[ConvertImageRequest]

    C --> H[CovertOpenAI2Gemini]
    H --> I[ThinkingAdaptor]
    I --> J[构建 SafetySettings]
    J --> K[处理 Tools]
    K --> L[转换 Messages]

    E --> C

    subgraph URL构建
        M[GetRequestURL]
        M --> N{模型类型?}
        N -->|imagen| O[/models/{model}:predict]
        N -->|embedding| P[/models/{model}:embedContent]
        N -->|chat 非流式| Q[/models/{model}:generateContent]
        N -->|chat 流式| R[/models/{model}:streamGenerateContent?alt=sse]
    end

    L --> M

    subgraph 响应处理
        S[DoResponse]
        S --> T{RelayMode?}
        T -->|Gemini 原生| U[GeminiTextGenerationHandler]
        T -->|OpenAI 兼容| V{流式?}
        V -->|是| W[GeminiChatStreamHandler]
        V -->|否| X[GeminiChatHandler]
    end
```

## Gemini 特有功能

### 1. Thinking 思考模式

支持通过模型名称后缀或 extra_body 配置思考模式：

```go
// 模型名称后缀方式
"gemini-2.5-pro-thinking"      // 自动计算 budget
"gemini-2.5-pro-thinking-8192" // 指定 budget
"gemini-2.5-pro-nothinking"    // 禁用思考

// extra_body 方式
{
    "extra_body": {
        "google": {
            "thinking_config": {
                "thinking_budget": 5324,
                "include_thoughts": true
            }
        }
    }
}
```

**Budget 限制：**
- gemini-2.5-pro 新版：128 - 32768
- gemini-2.5-flash-lite：512 - 24576
- 其他模型：0 - 24576

### 2. Grounding 搜索增强

通过 Tools 传递 Google Search：

```json
{
    "tools": [{"function": {"name": "googleSearch"}}]
}
```

### 3. Code Execution 代码执行

```json
{
    "tools": [{"function": {"name": "codeExecution"}}]
}
```

### 4. URL Context

```json
{
    "tools": [{"function": {"name": "urlContext"}}]
}
```

### 5. Safety Settings 安全设置

```go
var SafetySettingList = []string{
    "HARM_CATEGORY_HARASSMENT",
    "HARM_CATEGORY_HATE_SPEECH",
    "HARM_CATEGORY_SEXUALLY_EXPLICIT",
    "HARM_CATEGORY_DANGEROUS_CONTENT",
}
```

### 6. 图像生成 (Imagen)

支持 OpenAI 图像 API 格式转换：

- 尺寸映射：`1024x1024` → `1:1`, `1792x1024` → `16:9`
- 质量映射：`hd`/`high` → `2K`, `standard` → `1K`

## 配置项 (setting/model_setting/gemini.go)

```go
type GeminiSettings struct {
    SafetySettings                        map[string]string // 安全设置阈值
    VersionSettings                       map[string]string // API 版本
    SupportedImagineModels                []string          // 支持图像输出的模型
    ThinkingAdapterEnabled                bool              // 思考适配器开关
    ThinkingAdapterBudgetTokensPercentage float64           // 默认思考预算比例
    FunctionCallThoughtSignatureEnabled   bool              // 函数调用签名
}
```

## 依赖关系

```
gemini/
├── dto/gemini.go              # Gemini 数据结构
├── relay/common/              # 中继公共工具
├── relay/channel/openai/      # OpenAI 适配器（复用部分逻辑）
├── relay/helper/              # 流式处理辅助
├── service/                   # 文件处理、响应工具
├── setting/model_setting/     # Gemini 配置
└── setting/reasoning/         # 推理后缀处理
```

## 错误处理

1. **请求转换错误**：返回 `types.NewOpenAIError`
2. **响应解析错误**：返回 `types.ErrorCodeBadResponseBody`
3. **内容过滤**：映射 `FinishReason` 为 `content_filter`
4. **图片数量限制**：通过 `GeminiVisionMaxImageNum` 控制

## 流式传输实现

```go
func geminiStreamHandler(c *gin.Context, info *relaycommon.RelayInfo,
    resp *http.Response,
    callback func(data string, geminiResponse *dto.GeminiChatResponse) bool) (*dto.Usage, *types.NewAPIError)
```

- 使用 `helper.StreamScannerHandler` 处理 SSE
- 实时统计 token 使用量
- 支持图片计数和 token 估算
- 最终发送 `[DONE]` 标记

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28 | 创建 | 首次生成模块文档 |
