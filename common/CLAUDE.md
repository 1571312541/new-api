[根目录](../CLAUDE.md) > **common**

# Common 模块 - 公共工具库

> 提供全局共享的工具函数、常量和基础设施。

## 模块职责

Common 模块提供：

1. **环境配置**：环境变量读取和初始化
2. **数据库工具**：数据库连接和操作工具
3. **缓存管理**：Redis 连接和缓存操作
4. **加密工具**：密码哈希、Token 生成
5. **验证工具**：邮箱、参数验证
6. **日志工具**：系统日志记录
7. **限流器**：请求限流实现

## 关键文件

### 环境配置

| 文件 | 说明 |
|------|------|
| `env.go` | 环境变量初始化 |
| `constants.go` | 全局常量定义 |

### 数据库

| 文件 | 说明 |
|------|------|
| `database.go` | 数据库工具 |
| `redis.go` | Redis 连接和操作 |

### 安全

| 文件 | 说明 |
|------|------|
| `crypto.go` | 加密工具 |
| `hash.go` | 哈希函数 |
| `totp.go` | TOTP 验证 |

### 验证

| 文件 | 说明 |
|------|------|
| `validate.go` | 参数验证 |
| `verification.go` | 验证码生成 |
| `email.go` | 邮件发送 |

### 工具

| 文件 | 说明 |
|------|------|
| `utils.go` | 通用工具函数 |
| `str.go` | 字符串处理 |
| `ip.go` | IP 地址处理 |
| `json.go` | JSON 工具 |
| `audio.go` | 音频处理工具 |

### 限流

| 文件 | 说明 |
|------|------|
| `rate-limit.go` | 限流工具 |
| `limiter/limiter.go` | 限流器实现 |
| `limiter/lua/rate_limit.lua` | Redis Lua 脚本 |

### 日志

| 文件 | 说明 |
|------|------|
| `sys_log.go` | 系统日志 |

## 对外接口

### 环境初始化

```go
func InitEnv()                          // 初始化环境变量
func GetEnvOrDefault(key string, def int) int
```

### 日志函数

```go
func SysLog(s string)                   // 系统日志
func SysError(s string)                 // 错误日志
func FatalLog(v ...any)                 // 致命错误
```

### 密码工具

```go
func Password2Hash(password string) (string, error)
func ValidatePassword(password, hash string) bool
```

### Redis 操作

```go
func InitRedisClient() error
func RedisSet(key string, value interface{}, exp time.Duration) error
func RedisGet(key string) (string, error)
```

### 验证工具

```go
func IsValidEmail(email string) bool
func IsValidUsername(username string) bool
```

## 全局变量

### 系统配置

```go
var Port *int                           // 服务端口
var Version string                      // 版本号
var SessionSecret string                // 会话密钥
var DebugEnabled bool                   // 调试模式
```

### 数据库配置

```go
var SQLitePath string                   // SQLite 路径
var UsingSQLite bool
var UsingMySQL bool
var UsingPostgreSQL bool
```

### Redis 配置

```go
var RedisEnabled bool
var MemoryCacheEnabled bool
var SyncFrequency int                   // 同步频率
```

## 常见问题 (FAQ)

**Q: 如何添加新的环境变量？**
A: 在 `env.go` 的 `InitEnv()` 函数中添加

**Q: 如何自定义日志格式？**
A: 修改 `sys_log.go` 中的日志函数

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28T20:58:40+0800 | 初始化 | 首次生成模块文档 |
