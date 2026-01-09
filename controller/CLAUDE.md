[根目录](../CLAUDE.md) > **controller**

# Controller 模块 - API 控制器层

> 处理所有 HTTP 请求的入口点，负责参数验证、调用服务层、返回响应。

## 模块职责

Controller 模块实现了 new-api 的所有 API 端点：

1. **用户管理**：注册、登录、权限、2FA 认证
2. **渠道管理**：CRUD、测试、余额查询
3. **令牌管理**：API Key 的创建和管理
4. **计费管理**：充值、配额、账单
5. **中继转发**：AI API 请求代理
6. **任务管理**：Midjourney、视频生成等异步任务

## 入口与启动

控制器由路由器调用，主要入口在 `router/api-router.go` 和 `router/relay-router.go`。

## 对外接口

### 用户相关 (`user.go`)

| 路由 | 方法 | 函数 | 说明 |
|------|------|------|------|
| `/api/user/register` | POST | `Register` | 用户注册 |
| `/api/user/login` | POST | `Login` | 用户登录 |
| `/api/user/self` | GET | `GetSelf` | 获取当前用户信息 |

### 渠道相关 (`channel.go`)

| 路由 | 方法 | 函数 | 说明 |
|------|------|------|------|
| `/api/channel/` | GET | `GetAllChannels` | 获取所有渠道 |
| `/api/channel/:id` | GET | `GetChannel` | 获取单个渠道 |
| `/api/channel/` | POST | `AddChannel` | 添加渠道 |
| `/api/channel/test/:id` | GET | `TestChannel` | 测试渠道 |

### 令牌相关 (`token.go`)

| 路由 | 方法 | 函数 | 说明 |
|------|------|------|------|
| `/api/token/` | GET | `GetAllTokens` | 获取所有令牌 |
| `/api/token/` | POST | `AddToken` | 创建令牌 |

### 中继相关 (`relay.go`)

| 函数 | 说明 |
|------|------|
| `Relay()` | 统一的 API 中继入口 |
| `RelayMidjourney()` | Midjourney 请求处理 |
| `RelayTask()` | 异步任务请求处理 |

## 关键依赖与配置

- **model/**：数据模型
- **service/**：业务逻辑
- **middleware/**：认证中间件
- **dto/**：请求/响应结构

## 数据模型

控制器使用以下主要模型：

- `model.User` - 用户
- `model.Channel` - 渠道
- `model.Token` - 令牌
- `model.Log` - 日志
- `model.Redemption` - 兑换码

## 测试与质量

目前没有单元测试文件。建议添加：
- 请求参数验证测试
- 权限控制测试
- 边界条件测试

## 常见问题 (FAQ)

**Q: 如何添加新的 API 端点？**
A:
1. 在对应的控制器文件中添加处理函数
2. 在 `router/api-router.go` 中注册路由
3. 添加必要的中间件（认证、限流等）

**Q: 如何实现管理员专属功能？**
A: 使用 `middleware.AdminAuth()` 中间件

## 相关文件清单

| 文件 | 说明 |
|------|------|
| `user.go` | 用户管理 |
| `channel.go` | 渠道管理 |
| `token.go` | 令牌管理 |
| `relay.go` | API 中继 |
| `billing.go` | 计费管理 |
| `topup.go` | 充值管理 |
| `log.go` | 日志管理 |
| `midjourney.go` | Midjourney 任务 |
| `task.go` | 通用任务 |
| `pricing.go` | 定价管理 |
| `option.go` | 系统配置 |
| `model.go` | 模型管理 |
| `setup.go` | 初始化设置 |

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28T20:58:40+0800 | 初始化 | 首次生成模块文档 |
