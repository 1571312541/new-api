[根目录](../CLAUDE.md) > **constant**

# Constant 模块 - 常量定义

> 定义系统全局常量和枚举值。

## 模块职责

Constant 模块定义：

1. **渠道类型**：所有支持的 AI 渠道
2. **API 类型**：请求格式类型
3. **缓存键**：Redis 缓存键前缀
4. **上下文键**：Gin Context 键
5. **环境变量**：环境变量名称
6. **任务类型**：异步任务类型

## 关键文件

### 渠道定义 (`channel.go`)

```go
const (
    ChannelTypeOpenAI         = 1
    ChannelTypeAzure          = 3
    ChannelTypeOllama         = 4
    ChannelTypeAnthropic      = 14
    ChannelTypeBaidu          = 15
    ChannelTypeAli            = 17
    ChannelTypeGemini         = 24
    ChannelTypeAws            = 33
    ChannelTypeDeepSeek       = 43
    // ... 50+ 渠道
)

var ChannelBaseURLs = []string{...}    // 渠道默认 URL
var ChannelTypeNames = map[int]string{...} // 渠道名称
```

### API 类型 (`api_type.go`)

```go
const (
    APITypeOpenAI   = 1
    APITypeClaude   = 2
    APITypeGemini   = 3
    // ...
)
```

### 缓存键 (`cache_key.go`)

```go
const (
    CacheKeyToken   = "token:"
    CacheKeyChannel = "channel:"
    CacheKeyUser    = "user:"
    // ...
)
```

### 上下文键 (`context_key.go`)

```go
const (
    ContextKeyToken    = "token"
    ContextKeyUser     = "user"
    ContextKeyChannel  = "channel"
    ContextKeyModel    = "model"
    // ...
)
```

### 环境变量 (`env.go`)

环境变量名称常量定义。

### 任务类型 (`task.go`)

```go
const (
    TaskTypeMidjourney = "midjourney"
    TaskTypeSuno       = "suno"
    TaskTypeVideo      = "video"
    // ...
)
```

### Midjourney (`midjourney.go`)

Midjourney 相关常量。

### 结束原因 (`finish_reason.go`)

```go
const (
    FinishReasonStop       = "stop"
    FinishReasonLength     = "length"
    FinishReasonToolCalls  = "tool_calls"
    // ...
)
```

### Azure (`azure.go`)

Azure OpenAI 相关常量。

### 端点类型 (`endpoint_type.go`)

API 端点类型常量。

## 使用示例

```go
import "github.com/QuantumNous/new-api/constant"

// 检查渠道类型
if channel.Type == constant.ChannelTypeOpenAI {
    // OpenAI 渠道处理
}

// 获取渠道名称
name := constant.GetChannelTypeName(channel.Type)

// 获取默认 URL
baseURL := constant.ChannelBaseURLs[channel.Type]
```

## 常见问题 (FAQ)

**Q: 如何添加新渠道类型？**
A:
1. 在 `channel.go` 添加常量
2. 添加到 `ChannelBaseURLs` 和 `ChannelTypeNames`
3. 实现渠道适配器

**Q: 常量命名规范？**
A: 使用 PascalCase，如 `ChannelTypeOpenAI`

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28T20:58:40+0800 | 初始化 | 首次生成模块文档 |
