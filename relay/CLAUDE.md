[根目录](../CLAUDE.md) > **relay**

# Relay 模块 - API 中继核心

> 本模块是 new-api 的核心，负责将用户请求代理转发到各种 AI 服务提供商。

## 模块职责

Relay 模块实现了多渠道 AI API 的统一代理功能：

1. **请求转换**：将统一的 OpenAI 格式请求转换为各渠道特定格式
2. **响应处理**：将各渠道响应转换回统一格式
3. **流式处理**：支持 SSE 流式响应
4. **错误处理**：统一的错误处理和重试机制

## 入口与启动

- **主要入口**：由 `router/relay-router.go` 调用 `controller.Relay()` 触发
- **音频处理**：`audio_handler.go`

## 对外接口

### Adaptor 接口 (relay/channel/adapter.go)

所有渠道适配器必须实现此接口：

```go
type Adaptor interface {
    Init(info *relaycommon.RelayInfo)
    GetRequestURL(info *relaycommon.RelayInfo) (string, error)
    SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error
    ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error)
    DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error)
    DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError)
    GetModelList() []string
    GetChannelName() string
    // ... 更多方法
}
```

### TaskAdaptor 接口

用于异步任务类 API（如 Midjourney、视频生成）：

```go
type TaskAdaptor interface {
    ValidateRequestAndSetAction(c *gin.Context, info *relaycommon.RelayInfo) *dto.TaskError
    BuildRequestURL(info *relaycommon.RelayInfo) (string, error)
    DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (taskID string, taskData []byte, err *dto.TaskError)
    FetchTask(baseUrl, key string, body map[string]any, proxy string) (*http.Response, error)
    ParseTaskResult(respBody []byte) (*relaycommon.TaskInfo, error)
}
```

## 关键依赖与配置

- **dto/**：请求/响应数据结构
- **types/**：类型定义
- **constant/**：渠道类型常量
- **relay/common/**：中继公共工具

## 渠道适配器目录

| 渠道 | 路径 | 说明 |
|------|------|------|
| OpenAI | `channel/openai/` | 支持 GPT 系列、DALL-E、TTS 等 |
| Claude | `channel/claude/` | Anthropic Claude API |
| Gemini | `channel/gemini/` | Google Gemini API |
| Ali/Qwen | `channel/ali/` | 阿里云通义千问 |
| AWS Bedrock | `channel/aws/` | AWS Bedrock 服务 |
| Baidu | `channel/baidu/` | 百度文心一言 |
| DeepSeek | `channel/deepseek/` | DeepSeek API |
| Ollama | `channel/ollama/` | 本地 Ollama 服务 |
| Cohere | `channel/cohere/` | Cohere API (含 Rerank) |
| MiniMax | `channel/minimax/` | MiniMax API |
| Mistral | `channel/mistral/` | Mistral AI |
| SiliconFlow | `channel/siliconflow/` | 硅基流动 |
| Coze | `channel/coze/` | 扣子 API |
| Dify | `channel/dify/` | Dify ChatFlow |
| Jina | `channel/jina/` | Jina Embeddings/Rerank |

## 添加新渠道指南

1. **在 `constant/channel.go` 添加渠道类型**：
   ```go
   const ChannelTypeNewChannel = 57
   ```

2. **创建渠道目录**：`relay/channel/newchannel/`

3. **实现必要文件**：
   - `adaptor.go` - 实现 Adaptor 接口
   - `constants.go` - 渠道常量和模型列表
   - `dto.go` - 请求/响应结构（如有特殊格式）

4. **在适配器工厂注册**（通常在 `relay/channel/` 下的初始化代码中）

## 常见问题 (FAQ)

**Q: 如何调试渠道请求？**
A: 设置环境变量 `GIN_MODE=debug` 并启用 `ERROR_LOG_ENABLED=true`

**Q: 流式响应超时怎么办？**
A: 增加 `STREAMING_TIMEOUT` 环境变量值（默认 300 秒）

**Q: 如何处理渠道特有的认证方式？**
A: 在 `SetupRequestHeader()` 方法中实现特定的认证逻辑

## 相关文件清单

| 文件 | 说明 |
|------|------|
| `audio_handler.go` | 音频请求处理 |
| `channel/adapter.go` | 适配器接口定义 |
| `channel/api_request.go` | API 请求工具 |
| `common/` | 中继公共工具 |

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28T20:58:40+0800 | 初始化 | 首次生成模块文档 |
