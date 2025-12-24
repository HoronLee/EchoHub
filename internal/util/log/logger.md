# EchoHub 日志系统

基于 Zap 的高性能日志系统，支持依赖注入和模式切换。

## 特性

- ✅ 依赖注入：通过 Wire 自动注入
- ✅ 模式切换：根据 `server.mode` 自动配置
- ✅ Gin 中间件：自动记录 HTTP 请求
- ✅ Gorm 集成：数据库查询日志
- ✅ 高性能：基于 uber-go/zap

## 配置

### Debug 模式

```yaml
server:
  mode: "debug"
```

**特点：**
- 输出到控制台
- 带颜色显示
- 非 JSON 格式（易读）
- 日志级别：Debug

**示例输出：**
```
2025-12-09T18:41:37.651+0800	INFO	data/data.go:54	Database connected successfully	{"driver": "sqlite"}
2025-12-09T18:41:37.652+0800	DEBUG	log/gorm.go:66	Database Query	{"sql": "SELECT count(*) FROM...", "elapsed": 0.000030584}
```

### Release 模式

```yaml
server:
  mode: "release"
```

**特点：**
- 同时输出到控制台和文件
- JSON 格式（便于解析）
- 无颜色信息
- 日志级别：Info
- 文件路径：`logs/app.log`
- 自动轮转：100MB/文件，保留5个备份，30天

**示例输出：**
```json
{"level":"info","time":"2025-12-09T18:41:53.530+0800","caller":"data/data.go:54","msg":"Database connected successfully","driver":"sqlite"}
{"level":"info","time":"2025-12-09T18:41:56.544+0800","caller":"middleware/logger.go:21","msg":"HTTP Request","status":200,"method":"POST","path":"/api/helloworld"}
```

## 使用方法

### 1. 在服务中使用

通过依赖注入自动获取 logger：

```go
type MyService struct {
    logger *util.Logger
}

func NewMyService(logger *util.Logger) *MyService {
    return &MyService{logger: logger}
}

func (s *MyService) DoSomething() {
    s.logger.Info("Doing something", zap.String("key", "value"))
}
```

### 2. HTTP 请求日志

自动记录所有 HTTP 请求（已在 `server/http.go` 中配置）：

```go
engine.Use(middleware.Logger(logger))
```

记录内容：
- 状态码
- 请求方法
- 路径
- 查询参数
- 客户端 IP
- 响应时间
- User-Agent
- 错误信息

### 3. 数据库日志

Gorm 自动使用 Zap 记录数据库操作（已在 `data/data.go` 中配置）：

```go
gormLogger := util.NewGormLogger(logger)
db, err := gorm.Open(dialector, &gorm.Config{
    Logger: gormLogger,
})
```

记录内容：
- SQL 语句
- 执行时间
- 影响行数
- 错误信息
- 慢查询警告（>200ms）

## 架构设计

### 依赖注入流程

```
config.AppConfig
    ↓
server.NewLogger(mode) → *util.Logger
    ↓
├─→ data.NewDB(cfg, logger) → *gorm.DB
│       ↓
│   util.NewGormLogger(logger)
│
└─→ server.NewHTTPServer(cfg, handlers, db, logger)
        ↓
    middleware.Logger(logger)
```

### 文件结构

```
internal/
├── util/log/
│   ├── log.go       # 核心日志器
│   └── gorm.go      # Gorm 适配器
├── middleware/
│   └── logger.go    # Gin 中间件
├── server/
│   └── server.go    # Logger Provider
└── data/
    └── data.go      # 数据库集成
```

## 日志级别

- **Debug**: 详细的调试信息（仅 debug 模式）
- **Info**: 一般信息（默认级别）
- **Warn**: 警告信息（如慢查询）
- **Error**: 错误信息
- **Fatal**: 致命错误（会终止程序）

## 最佳实践

1. **使用结构化日志**
   ```go
   logger.Info("User login", 
       zap.String("username", username),
       zap.Int("user_id", userID),
   )
   ```

2. **避免敏感信息**
   ```go
   // ❌ 不要记录密码
   logger.Info("Login", zap.String("password", password))
   
   // ✅ 只记录必要信息
   logger.Info("Login", zap.String("username", username))
   ```

3. **合理使用日志级别**
   - Debug: 开发调试
   - Info: 正常业务流程
   - Warn: 需要关注但不影响运行
   - Error: 错误但可恢复
   - Fatal: 致命错误

4. **性能考虑**
   - Zap 是零分配日志库，性能极高
   - 避免在循环中频繁记录 Debug 日志
   - Release 模式会自动过滤 Debug 日志

## 参考

- [Zap 文档](https://github.com/uber-go/zap)
- [Lumberjack 日志轮转](https://github.com/natefinch/lumberjack)
