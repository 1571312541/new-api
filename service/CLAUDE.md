[根目录](../CLAUDE.md) > **service**

# Service 模块 - 业务逻辑层

> 封装核心业务逻辑，为控制器层提供服务接口。

## 模块职责

Service 模块实现了 new-api 的核心业务逻辑：

1. **渠道服务**：渠道选择、负载均衡
2. **配额服务**：配额计算、预扣费、结算
3. **Token 计数**：各模型的 Token 统计
4. **文件处理**：图片、音频转换
5. **通知服务**：用户通知、Webhook

## 入口与启动

服务层由控制器和中间件调用，无独立启动入口。

## 对外接口

### 渠道服务 (`channel.go`, `channel_select.go`)

| 函数 | 说明 |
|------|------|
| `GetRandomSatisfiedChannel()` | 获取可用渠道 |
| `SelectChannel()` | 渠道选择（带负载均衡） |

### 配额服务 (`quota.go`, `pre_consume_quota.go`)

| 函数 | 说明 |
|------|------|
| `PreConsumeQuota()` | 预扣配额 |
| `PostConsumeQuota()` | 结算配额 |
| `ReturnPreConsumedQuota()` | 退还预扣配额 |

### Token 计数服务 (`tokenizer.go`, `token_counter.go`, `token_estimator.go`)

| 函数 | 说明 |
|------|------|
| `InitTokenEncoders()` | 初始化编码器 |
| `CountTokens()` | 计算 Token 数量 |
| `EstimateTokens()` | 估算 Token |

### 文件处理服务

| 函数 | 文件 | 说明 |
|------|------|------|
| `ProcessImage()` | `image.go` | 图片处理 |
| `ProcessAudio()` | `audio.go` | 音频处理 |
| `ConvertFile()` | `convert.go` | 文件格式转换 |
| `DecodeFile()` | `file_decoder.go` | 文件解码 |
| `DownloadFile()` | `download.go` | 文件下载 |

### HTTP 服务 (`http.go`, `http_client.go`)

| 函数 | 说明 |
|------|------|
| `InitHttpClient()` | 初始化 HTTP 客户端 |
| `DoRequest()` | 执行 HTTP 请求 |

### 通知服务 (`user_notify.go`, `webhook.go`)

| 函数 | 说明 |
|------|------|
| `NotifyUser()` | 发送用户通知 |
| `SendWebhook()` | 发送 Webhook |

## 关键依赖与配置

- **model/**：数据访问
- **common/**：工具函数
- **dto/**：数据结构
- **setting/**：配置读取

## 配额计算逻辑

### 预扣费流程

```
请求 -> 估算 Token -> 计算预扣金额 -> 扣除用户配额
    -> 执行 API 请求
    -> 获取实际用量 -> 计算实际费用 -> 结算差额
```

### 价格计算

```go
// 输入价格
inputPrice = tokens * modelInputRatio * groupRatio

// 输出价格
outputPrice = tokens * modelOutputRatio * groupRatio

// 总价
totalPrice = inputPrice + outputPrice
```

## Token 计数

支持多种编码器：

- **cl100k_base**：GPT-4, GPT-3.5-turbo
- **p50k_base**：GPT-3
- **Claude tokenizer**：Claude 系列
- **通用估算**：其他模型

## 常见问题 (FAQ)

**Q: 如何添加新的定价策略？**
A: 在 `setting/ratio_setting/` 中添加配置

**Q: 如何实现自定义渠道选择逻辑？**
A: 修改 `channel_select.go` 中的选择算法

**Q: Token 计数不准确怎么办？**
A: 检查模型对应的编码器配置

## 相关文件清单

| 文件 | 说明 |
|------|------|
| `channel.go` | 渠道服务 |
| `channel_select.go` | 渠道选择 |
| `quota.go` | 配额管理 |
| `pre_consume_quota.go` | 预扣费 |
| `tokenizer.go` | Token 编码 |
| `token_counter.go` | Token 计数 |
| `token_estimator.go` | Token 估算 |
| `image.go` | 图片处理 |
| `audio.go` | 音频处理 |
| `http_client.go` | HTTP 客户端 |
| `user_notify.go` | 用户通知 |
| `webhook.go` | Webhook |
| `error.go` | 错误处理 |
| `log_info_generate.go` | 日志生成 |
| `usage_helpr.go` | 用量计算 |

### Passkey 子模块 (`passkey/`)

| 文件 | 说明 |
|------|------|
| `service.go` | Passkey 服务 |
| `session.go` | 会话管理 |
| `user.go` | 用户关联 |

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28T20:58:40+0800 | 初始化 | 首次生成模块文档 |
