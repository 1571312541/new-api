[根目录](../CLAUDE.md) > **setting**

# Setting 模块 - 配置管理

> 管理系统配置，包括定价、限流、渠道设置等。

## 模块职责

Setting 模块负责：

1. **价格配置**：模型定价、计费比例
2. **限流配置**：请求限流规则
3. **系统设置**：认证、支付、监控
4. **渠道设置**：渠道特有配置

## 目录结构

```
setting/
├── ratio_setting/          # 价格比例设置
│   ├── model_ratio.go      # 模型价格比例
│   ├── group_ratio.go      # 用户组价格比例
│   ├── cache_ratio.go      # 缓存价格
│   ├── expose_ratio.go     # 暴露价格配置
│   └── exposed_cache.go    # 缓存配置
├── operation_setting/      # 运营设置
│   ├── operation_setting.go
│   ├── general_setting.go
│   ├── payment_setting.go
│   ├── quota_setting.go
│   └── monitor_setting.go
├── system_setting/         # 系统设置
│   ├── discord.go
│   ├── oidc.go
│   ├── passkey.go
│   └── legal.go
├── model_setting/          # 模型设置
│   ├── global.go
│   ├── claude.go
│   └── gemini.go
├── console_setting/        # 控制台设置
│   ├── config.go
│   └── validation.go
├── reasoning/              # 推理设置
│   └── suffix.go
├── config/                 # 通用配置
│   └── config.go
├── auto_group.go           # 自动分组
├── chat.go                 # 聊天设置
├── midjourney.go           # Midjourney 设置
├── payment_creem.go        # Creem 支付
├── payment_stripe.go       # Stripe 支付
├── rate_limit.go           # 限流设置
├── sensitive.go            # 敏感词设置
└── user_usable_group.go    # 用户可用组
```

## 关键配置

### 模型价格比例 (`ratio_setting/model_ratio.go`)

```go
// 模型输入/输出价格比例
var ModelRatio map[string]float64
var ModelCompletionRatio map[string]float64
```

### 用户组价格比例 (`ratio_setting/group_ratio.go`)

```go
// 用户组价格倍率
var GroupRatio map[string]float64
```

### 限流配置 (`rate_limit.go`)

```go
// 请求限流配置
var RateLimitConfig RateLimitSetting
```

### 支付配置

| 文件 | 说明 |
|------|------|
| `payment_stripe.go` | Stripe 支付集成 |
| `payment_creem.go` | Creem 支付集成 |

## 配置加载

配置通过以下方式加载：
1. **环境变量**：启动时读取
2. **数据库**：Option 表存储
3. **内存缓存**：运行时缓存

### 初始化流程

```go
// main.go
ratio_setting.InitRatioSettings()
model.InitOptionMap()
```

## 常见问题 (FAQ)

**Q: 如何修改模型定价？**
A: 在管理后台或修改 `ratio_setting/model_ratio.go`

**Q: 如何添加新的支付方式？**
A: 参考 `payment_stripe.go` 实现新的支付模块

**Q: 配置优先级是什么？**
A: 数据库配置 > 环境变量 > 代码默认值

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28T20:58:40+0800 | 初始化 | 首次生成模块文档 |
