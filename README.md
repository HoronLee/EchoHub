# EchoHub

> 基于 Echo、GORM、Viper、Wire、Cobra 的 HTTP 快速开发框架

EchoHub 是一个生产级的 Go HTTP 服务快速开发框架，提供了完整的分层架构、依赖注入、响应封装、中间件、日志系统和 Swagger 文档等企业级特性。

## 特性

- **分层架构**: 清晰的 MVC 分层设计 (Handler -> Service -> Data)
- **依赖注入**: 使用 Google Wire 自动生成依赖注入代码
- **响应封装**: Execute 模式统一处理响应和错误
- **中间件系统**: 内置日志、恢复、JWT 认证中间件
- **双日志系统**: Web 日志 + ORM 日志，支持 Debug/Release 模式
- **配置管理**: 基于 Viper，支持嵌入式配置和外部配置覆盖
- **CLI/TUI**: 基于 Cobra 的命令行工具和交互式界面
- **Swagger 文档**: 自动生成 API 文档
- **JWT 认证**: 内置用户认证和授权

## 项目结构

```
EchoHub/
├── cmd/                    # CLI 命令入口
│   ├── init.go            # 命令初始化
│   └── root.go            # 根命令定义
├── configs/               # 外部配置文件
│   ├── debug-config.yaml  # Debug 模式配置
│   └── production-config.yaml # 生产环境配置
├── internal/
│   ├── cli/              # CLI/TUI 逻辑
│   ├── config/           # 配置管理
│   ├── data/             # 数据访问层 (Repository)
│   ├── di/               # 依赖注入 (Wire)
│   ├── handler/          # HTTP 处理器 (Controller)
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   ├── response/         # 响应封装
│   ├── router/           # 路由配置
│   ├── server/           # HTTP 服务器
│   ├── service/          # 业务逻辑层
│   ├── swagger/          # Swagger 文档
│   ├── tui/              # TUI 界面
│   └── util/             # 工具函数
├── main.go               # 程序入口
├── Makefile              # 构建脚本
└── go.mod                # Go 模块定义
```

## 核心概念

### 1. Execute 响应封装模式

EchoHub 的核心设计模式是 `Execute` 包装器，它统一处理 HTTP 响应和错误：

```go
import res "github.com/HoronLee/EchoHub/internal/response"

func (h *UserHandler) Register() echo.HandlerFunc {
    return res.Execute(func(ctx echo.Context) res.Response {
        var req user.RegisterRequest
        if err := ctx.Bind(&req); err != nil {
            return res.Response{Msg: "Invalid request body", Err: err}
        }

        if err := h.svc.Register(ctx.Request().Context(), req); err != nil {
            return res.Response{Msg: "Registration failed", Err: err}
        }

        return res.Response{
            Data: map[string]any{"message": "User registered successfully"},
            Msg:  "success",
        }
    })
}
```

**Response 结构**:
- `Code`: 自定义业务状态码 (0 表示成功)
- `Data`: 响应数据
- `Msg`: 响应消息
- `Err`: 错误信息 (不序列化到 JSON)

### 2. 依赖注入 (Wire)

使用 Google Wire 自动管理依赖关系：

```go
// internal/di/wire.go
//go:build wireinject
// +build wireinject

func InitServer(cfg *config.AppConfig) (*server.HTTPServer, func(), error) {
    wire.Build(
        util.NewLogger,
        data.ProviderSet,
        service.ProviderSet,
        handler.ProviderSet,
        server.ProviderSet,
    )
    return nil, nil, nil
}
```

运行 `make wire` 生成 `wire_gen.go`。

### 3. 分层架构

```
Request -> Middleware -> Router -> Handler -> Service -> Data -> Database
                  |         |         |         |       |
                  v         v         v         v       v
                Logger   Execute   Business  Repository  GORM
                Recovery  Response   Logic
                JWT
```

## 快速开始

### 安装依赖

```bash
# 安装 Swagger 和 Wire 工具
make swagger-install

# 整理 Go 依赖
make deps
```

### 配置

编辑 `configs/debug-config.yaml`:

```yaml
server:
  port: "8080"
  host: "0.0.0.0"
  mode: "debug"

database:
  type: "mysql"  # 或 "sqlite"
  source: "user:password@tcp(localhost:3306)/echohub?charset=utf8mb4&parseTime=True"
  logmode: "debug"

auth:
  jwt:
    secret: "your-secret-key"
    expires: 86400  # 24小时
```

### 开发

```bash
# 启动开发服务器
make dev

# 生成 Swagger 文档
make swagger

# 生成 Wire 代码
make wire

# 构建二进制文件
make build
```

### 访问 Swagger UI

启动服务后访问: `http://localhost:8080/api/v1/swagger/`

## 开发流程

### 添加新功能

以添加"用户资料"功能为例：

#### 1. 定义数据模型 (Model)

```go
// internal/model/user/profile.go
package user

type Profile struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    UserID   uint   `json:"user_id" gorm:"uniqueIndex"`
    Nickname string `json:"nickname" gorm:"size:50"`
    Bio      string `json:"bio" gorm:"size:500"`
}

type ProfileResponse struct {
    Nickname string `json:"nickname" example:"张三"`
    Bio      string `json:"bio" example:"全栈开发者"`
}

type UpdateProfileRequest struct {
    Nickname string `json:"nickname" binding:"required,min=1,max=50" example:"张三"`
    Bio      string `json:"bio" binding:"max=500" example:"全栈开发者"`
}
```

#### 2. 实现数据访问层 (Data)

```go
// internal/data/user.go
func (d *Data) GetProfileByUserID(ctx context.Context, userID uint) (*user.Profile, error) {
    var profile user.Profile
    err := d.db.Where("user_id = ?", userID).First(&profile).Error
    if err != nil {
        return nil, err
    }
    return &profile, nil
}

func (d *Data) UpdateProfile(ctx context.Context, userID uint, nickname, bio string) error {
    return d.db.Where("user_id = ?", userID).
        Updates(&user.Profile{Nickname: nickname, Bio: bio}).Error
}
```

#### 3. 实现业务逻辑层 (Service)

```go
// internal/service/user.go
func (s *UserService) GetProfile(ctx context.Context, userID uint) (*user.ProfileResponse, error) {
    profile, err := s.data.GetProfileByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }
    return &user.ProfileResponse{
        Nickname: profile.Nickname,
        Bio:      profile.Bio,
    }, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req user.UpdateProfileRequest) error {
    return s.data.UpdateProfile(ctx, userID, req.Nickname, req.Bio)
}
```

#### 4. 实现 HTTP 处理器 (Handler)

```go
// internal/handler/user.go
// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的资料信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user.ProfileResponse "获取成功"
// @Failure 400 {object} res.Response "请求失败"
// @Router /user/profile [get]
func (h *UserHandler) GetProfile() echo.HandlerFunc {
    return res.Execute(func(ctx echo.Context) res.Response {
        userID := ctx.Get("user_id").(uint)
        profile, err := h.svc.GetProfile(ctx.Request().Context(), userID)
        if err != nil {
            return res.Response{Msg: "Failed to get profile", Err: err}
        }
        return res.Response{Data: profile, Msg: "success"}
    })
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前登录用户的资料信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.UpdateProfileRequest true "更新资料请求"
// @Success 200 {object} map[string]string "更新成功"
// @Failure 400 {object} res.Response "请求失败"
// @Router /user/profile [put]
func (h *UserHandler) UpdateProfile() echo.HandlerFunc {
    return res.Execute(func(ctx echo.Context) res.Response {
        userID := ctx.Get("user_id").(uint)
        var req user.UpdateProfileRequest
        if err := ctx.Bind(&req); err != nil {
            return res.Response{Msg: "Invalid request", Err: err}
        }
        if err := h.svc.UpdateProfile(ctx.Request().Context(), userID, req); err != nil {
            return res.Response{Msg: "Update failed", Err: err}
        }
        return res.Response{Msg: "Profile updated successfully"}
    })
}
```

#### 5. 注册路由 (Router)

```go
// internal/router/user.go
func setupV1UserRoutes(routerGroup *VersionedRouterGroup, h *handler.Handlers) {
    // 公开路由
    routerGroup.PublicRouter.POST("/register", h.UserHandler.Register())
    routerGroup.PublicRouter.POST("/login", h.UserHandler.Login())
    
    // 需要认证的路由
    routerGroup.PrivateRouter.DELETE("/user", h.UserHandler.DeleteUser())
    routerGroup.PrivateRouter.GET("/user/profile", h.UserHandler.GetProfile())
    routerGroup.PrivateRouter.PUT("/user/profile", h.UserHandler.UpdateProfile())
}
```

#### 6. 更新依赖注入 (DI)

```go
// internal/handler/handler.go
var ProviderSet = wire.NewSet(
    NewUserHandler,
    NewHelloWorldHandler,
)
```

#### 7. 生成并运行

```bash
# 生成 Swagger 文档
make swagger

# 生成 Wire 代码
make wire

# 启动开发服务器
make dev
```

### 中间件开发

```go
// internal/middleware/custom.go
func CustomMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(ctx echo.Context) error {
            // 前置处理
            log.Info("Request received")
            
            // 调用下一个处理器
            err := next(ctx)
            
            // 后置处理
            log.Info("Request completed")
            
            return err
        }
    }
}

// 在 server/http.go 中使用
e.Use(middleware.CustomMiddleware())
```

## 命令行工具

```bash
# 启动 HTTP 服务
./bin/echohub serve -c configs/debug-config.yaml

# 启动 TUI 界面
./bin/echohub tui

# 查看版本
./bin/echohub version

# 查看信息
./bin/echohub info

# 显示 Logo
./bin/echohub hello
```

## API 规范

### 统一响应格式

成功响应:
```json
{
  "code": 0,
  "data": { ... },
  "msg": "success"
}
```

错误响应:
```json
{
  "code": 0,
  "data": null,
  "msg": "error message"
}
```

### 认证方式

使用 Bearer Token 认证:

```
Authorization: Bearer <your-jwt-token>
```

## 环境变量

- `JWT_SECRET`: JWT 签名密钥 (优先级高于配置文件)

## 生产部署

```bash
# 构建生产版本
make build-prod

# 运行
./bin/echohub-prod serve -c configs/production-config.yaml
```

## 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| Web 框架 | [Echo](https://echo.labstack.com/) | 高性能 Go Web 框架 |
| ORM | [GORM](https://gorm.io/) | Go ORM 库 |
| 配置管理 | [Viper](https://github.com/spf13/viper) | 配置文件管理 |
| 依赖注入 | [Wire](https://github.com/google/wire) | 编译时依赖注入 |
| CLI 框架 | [Cobra](https://github.com/spf13/cobra) | 命令行工具 |
| 日志 | [Zap](https://github.com/uber-go/zap) | 结构化日志 |
| API 文档 | [Swaggo](https://github.com/swaggo/swag) | Swagger 自动生成 |

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
