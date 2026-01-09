[根目录](../CLAUDE.md) > **router**

# Router 模块 - 路由定义

> 定义所有 HTTP 路由，将 URL 路径映射到对应的控制器处理函数。

## 模块职责

Router 模块负责：

1. **API 路由**：管理后台 API
2. **中继路由**：AI API 代理
3. **视频路由**：视频生成 API
4. **静态资源**：前端文件服务
5. **Dashboard**：数据看板 API

## 入口与启动

- **主入口**：`main.go` -> `router.SetRouter(server, buildFS, indexPage)`
- **路由设置**：`main.go`

## 对外接口

### 路由组

| 文件 | 函数 | 路径前缀 | 说明 |
|------|------|----------|------|
| `api-router.go` | `SetApiRouter()` | `/api` | 管理后台 API |
| `relay-router.go` | `SetRelayRouter()` | `/v1` | OpenAI 兼容 API |
| `video-router.go` | `SetVideoRouter()` | `/v1/video` | 视频生成 API |
| `web-router.go` | `SetWebRouter()` | `/` | 前端静态文件 |
| `dashboard.go` | `SetDashboardRouter()` | `/dashboard` | 数据看板 |

## 关键依赖与配置

- **controller/**：控制器函数
- **middleware/**：中间件

## 路由结构

### API 路由 (`/api`)

```
/api
├── /setup              # 初始化设置
├── /status             # 系统状态
├── /user               # 用户管理
│   ├── POST /register  # 注册
│   ├── POST /login     # 登录
│   └── /self           # 个人信息
├── /channel            # 渠道管理 (Admin)
├── /token              # 令牌管理
├── /log                # 日志查询
├── /redemption         # 兑换码 (Admin)
├── /option             # 系统配置 (Root)
├── /pricing            # 定价查询
├── /models             # 模型管理 (Admin)
├── /deployments        # 部署管理 (Admin)
└── /mj                 # Midjourney 任务
```

### 中继路由 (`/v1`)

```
/v1
├── /models             # 模型列表
├── /chat/completions   # 聊天补全 (OpenAI)
├── /completions        # 文本补全
├── /messages           # Claude Messages
├── /embeddings         # 嵌入向量
├── /images/generations # 图像生成
├── /audio/*            # 音频处理
├── /rerank             # 重排序
├── /realtime           # 实时对话 (WebSocket)
└── /models/*           # Gemini 格式
```

### Midjourney 路由 (`/mj`)

```
/mj
├── /submit/imagine     # 生成图像
├── /submit/change      # 变换图像
├── /submit/describe    # 描述图像
├── /task/:id/fetch     # 获取任务
└── /insight-face/swap  # 换脸
```

## 中间件应用

### 认证中间件

| 中间件 | 应用场景 |
|--------|----------|
| `TokenAuth()` | AI API 请求 |
| `UserAuth()` | 登录用户功能 |
| `AdminAuth()` | 管理员功能 |
| `RootAuth()` | 超级管理员功能 |

### 限流中间件

| 中间件 | 应用场景 |
|--------|----------|
| `GlobalAPIRateLimit()` | 所有 API |
| `CriticalRateLimit()` | 关键操作 |
| `ModelRequestRateLimit()` | 模型请求 |

## 常见问题 (FAQ)

**Q: 如何添加新的 API 端点？**
A:
1. 在控制器中添加处理函数
2. 在对应路由文件中添加路由定义
3. 配置必要的中间件

**Q: 如何修改路由前缀？**
A: 修改对应路由文件中的 `router.Group()` 参数

**Q: 如何添加新的路由组？**
A:
1. 创建新的路由文件
2. 实现 `SetXxxRouter()` 函数
3. 在 `main.go` 中调用

## 相关文件清单

| 文件 | 说明 |
|------|------|
| `main.go` | 路由入口和合并 |
| `api-router.go` | 管理 API 路由 |
| `relay-router.go` | 中继 API 路由 |
| `video-router.go` | 视频 API 路由 |
| `web-router.go` | 前端静态资源 |
| `dashboard.go` | 数据看板路由 |

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28T20:58:40+0800 | 初始化 | 首次生成模块文档 |
