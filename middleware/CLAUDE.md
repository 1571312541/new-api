[根目录](../CLAUDE.md) > **middleware**

# Middleware 模块 - HTTP 中间件

> 提供认证、限流、日志、CORS 等 HTTP 请求处理中间件。

## 模块职责

Middleware 模块实现了所有 HTTP 请求的前置/后置处理：

1. **认证授权**：Token 认证、用户认证、管理员认证
2. **限流控制**：全局限流、模型限流、邮件验证限流
3. **请求分发**：渠道选择和负载均衡
4. **日志记录**：请求日志和错误日志
5. **安全控制**：CORS、Turnstile 验证

## 入口与启动

中间件在路由器中注册，由 Gin 框架自动调用。

## 对外接口

### 认证中间件 (`auth.go`)

| 函数 | 说明 | 使用场景 |
|------|------|----------|
| `TokenAuth()` | API Token 认证 | AI API 请求 |
| `UserAuth()` | 用户会话认证 | 管理后台 |
| `TryUserAuth()` | 尝试用户认证（可选） | 公开 + 用户页面 |
| `AdminAuth()` | 管理员权限认证 | 管理功能 |
| `RootAuth()` | 超级管理员认证 | 系统设置 |

### 限流中间件

| 函数 | 文件 | 说明 |
|------|------|------|
| `GlobalAPIRateLimit()` | `rate-limit.go` | 全局 API 限流 |
| `CriticalRateLimit()` | `rate-limit.go` | 关键操作限流 |
| `ModelRequestRateLimit()` | `model-rate-limit.go` | 模型请求限流 |
| `EmailVerificationRateLimit()` | `email-verification-rate-limit.go` | 邮件验证限流 |

### 分发中间件 (`distributor.go`)

| 函数 | 说明 |
|------|------|
| `Distribute()` | 渠道选择和请求分发 |

### 其他中间件

| 函数 | 文件 | 说明 |
|------|------|------|
| `CORS()` | `cors.go` | 跨域资源共享 |
| `RequestId()` | `request-id.go` | 请求 ID 生成 |
| `SetUpLogger()` | `logger.go` | 日志中间件 |
| `TurnstileCheck()` | `turnstile-check.go` | Cloudflare Turnstile 验证 |
| `DecompressRequestMiddleware()` | `gzip.go` | 请求解压 |
| `StatsMiddleware()` | `stats.go` | 统计中间件 |

## 关键依赖与配置

- **model/**：用户和令牌验证
- **service/**：渠道选择服务
- **common/limiter/**：限流器实现

## 认证流程

### Token 认证流程

```
请求 -> TokenAuth() -> 提取 Authorization Header
                    -> 验证 Token 格式
                    -> 查询数据库/缓存
                    -> 检查配额和状态
                    -> 设置上下文信息
                    -> 继续请求
```

### 用户认证流程

```
请求 -> UserAuth() -> 提取 Session
                   -> 验证用户 ID
                   -> 查询用户信息
                   -> 检查用户状态
                   -> 设置上下文信息
                   -> 继续请求
```

## 渠道分发逻辑 (`distributor.go`)

1. **解析请求模型**
2. **获取用户可用渠道列表**
3. **根据权重随机选择渠道**
4. **检查渠道状态和配额**
5. **设置渠道信息到上下文**

## 常见问题 (FAQ)

**Q: 如何自定义认证逻辑？**
A: 在 `auth.go` 中修改对应的认证函数

**Q: 如何调整限流阈值？**
A: 通过系统设置或环境变量配置

**Q: 如何添加新的中间件？**
A:
1. 创建新的中间件文件
2. 实现 `gin.HandlerFunc` 函数
3. 在路由器中注册

## 相关文件清单

| 文件 | 说明 |
|------|------|
| `auth.go` | 认证中间件 |
| `rate-limit.go` | 通用限流 |
| `model-rate-limit.go` | 模型限流 |
| `distributor.go` | 渠道分发 |
| `cors.go` | CORS 处理 |
| `logger.go` | 日志记录 |
| `request-id.go` | 请求 ID |
| `turnstile-check.go` | Turnstile 验证 |
| `gzip.go` | 解压中间件 |
| `stats.go` | 统计中间件 |
| `cache.go` | 缓存控制 |
| `recover.go` | 错误恢复 |
| `secure_verification.go` | 安全验证 |

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28T20:58:40+0800 | 初始化 | 首次生成模块文档 |
