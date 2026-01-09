[根目录](../CLAUDE.md) > **model**

# Model 模块 - 数据模型层

> 定义所有数据库表结构和 ORM 操作，使用 GORM 作为 ORM 框架。

## 模块职责

Model 模块负责：

1. **数据库连接**：支持 MySQL、PostgreSQL、SQLite
2. **表结构定义**：所有业务实体的 ORM 映射
3. **数据库迁移**：自动建表和字段更新
4. **缓存管理**：渠道、令牌、用户的内存缓存
5. **CRUD 操作**：所有数据库读写操作

## 入口与启动

- **初始化入口**：`main.go` -> `model.InitDB()`
- **日志数据库**：可配置独立的日志数据库 (`LOG_SQL_DSN`)

## 对外接口

### 数据库初始化 (`main.go`)

```go
func InitDB() error        // 初始化主数据库
func InitLogDB() error     // 初始化日志数据库
func CloseDB() error       // 关闭数据库连接
```

### 核心模型

| 模型 | 文件 | 说明 |
|------|------|------|
| `User` | `user.go` | 用户账户 |
| `Channel` | `channel.go` | AI 渠道配置 |
| `Token` | `token.go` | API 令牌 |
| `Log` | `log.go` | 请求日志 |
| `Ability` | `ability.go` | 渠道能力（模型支持） |
| `Redemption` | `redemption.go` | 兑换码 |
| `Midjourney` | `midjourney.go` | MJ 任务 |
| `Task` | `task.go` | 通用异步任务 |
| `Model` | `model_meta.go` | 模型元数据 |
| `Vendor` | `vendor_meta.go` | 供应商元数据 |
| `Option` | `option.go` | 系统配置 |
| `Pricing` | `pricing.go` | 定价配置 |

## 关键依赖与配置

### 环境变量

| 变量 | 说明 | 示例 |
|------|------|------|
| `SQL_DSN` | 主数据库连接 | `root:pass@tcp(localhost:3306)/newapi` |
| `LOG_SQL_DSN` | 日志数据库（可选） | 同上 |
| `SQL_MAX_IDLE_CONNS` | 最大空闲连接数 | `100` |
| `SQL_MAX_OPEN_CONNS` | 最大连接数 | `1000` |

### 数据库支持

```go
// PostgreSQL
postgres://user:pass@host:5432/db

// MySQL
user:pass@tcp(host:3306)/db

// SQLite (默认)
local 或不设置 SQL_DSN
```

## 数据模型

### User 用户模型

```go
type User struct {
    Id          int    `gorm:"primaryKey"`
    Username    string `gorm:"unique;index"`
    Password    string
    Role        int    // 1-普通用户, 10-管理员, 100-超级管理员
    Status      int
    Quota       int64
    AccessToken *string
    // ...
}
```

### Channel 渠道模型

```go
type Channel struct {
    Id          int    `gorm:"primaryKey"`
    Name        string
    Type        int    // 渠道类型 (见 constant/channel.go)
    Key         string // API 密钥
    BaseURL     *string
    Models      string // 支持的模型列表
    Status      int
    Priority    *int64 // 优先级（用于负载均衡）
    // ...
}
```

### Token 令牌模型

```go
type Token struct {
    Id            int    `gorm:"primaryKey"`
    UserId        int
    Key           string `gorm:"unique;index"`
    Name          string
    RemainQuota   int64
    Models        string // 限制可用模型
    ExpiredTime   int64
    // ...
}
```

## 缓存机制

### 渠道缓存 (`channel_cache.go`)

- **InitChannelCache()** - 初始化缓存
- **SyncChannelCache()** - 定时同步
- **GetRandomSatisfiedChannel()** - 获取可用渠道

### 令牌缓存 (`token_cache.go`)

- **ValidateAndGetToken()** - 验证并获取令牌

### 用户缓存 (`user_cache.go`)

- **GetUserById()** - 带缓存的用户查询

## 测试与质量

目前没有单元测试。建议添加：
- 模型 CRUD 测试
- 缓存一致性测试
- 数据库迁移测试

## 常见问题 (FAQ)

**Q: 如何添加新的数据模型？**
A:
1. 创建模型文件定义结构体
2. 在 `migrateDB()` 中添加 AutoMigrate
3. 实现必要的 CRUD 方法

**Q: 如何执行数据库迁移？**
A: GORM AutoMigrate 会自动处理，启动时自动执行

**Q: 如何支持新数据库类型？**
A: 在 `chooseDB()` 函数中添加新的数据库驱动支持

## 相关文件清单

| 文件 | 说明 |
|------|------|
| `main.go` | 数据库初始化和迁移 |
| `user.go` | 用户模型 |
| `channel.go` | 渠道模型 |
| `token.go` | 令牌模型 |
| `log.go` | 日志模型 |
| `ability.go` | 渠道能力 |
| `channel_cache.go` | 渠道缓存 |
| `token_cache.go` | 令牌缓存 |
| `option.go` | 系统配置 |
| `pricing.go` | 定价管理 |

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28T20:58:40+0800 | 初始化 | 首次生成模块文档 |
