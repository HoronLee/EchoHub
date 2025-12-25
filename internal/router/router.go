package router

import (
	"github.com/HoronLee/EchoHub/internal/handler"
	"github.com/HoronLee/EchoHub/internal/middleware"
	"github.com/labstack/echo/v4"
)

// VersionedRouterGroup 版本化路由组
type VersionedRouterGroup struct {
	PublicRouter  *echo.Group
	PrivateRouter *echo.Group
}

// SetupRouter 配置路由
func SetupRouter(e *echo.Echo, h *handler.Handlers) {
	// 设置 v1 版本路由
	v1RouterGroup := setupV1RouterGroup(e)
	setupV1Routes(v1RouterGroup, h)

	// 设置资源路由（包括 Swagger UI）
	setupResourceRoutes(v1RouterGroup, h)

	// 未来可以添加 v2 版本路由
	// v2RouterGroup := setupV2RouterGroup(e)
	// setupV2Routes(v2RouterGroup, h)
}

// setupV1RouterGroup 初始化 v1 版本路由组
func setupV1RouterGroup(e *echo.Echo) *VersionedRouterGroup {
	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1")

	public := v1Group.Group("")
	private := v1Group.Group("")
	private.Use(middleware.JwtAuth()) // JWT认证中间件

	return &VersionedRouterGroup{
		PublicRouter:  public,
		PrivateRouter: private,
	}
}

// setupV1Routes 设置 v1 版本的所有路由
func setupV1Routes(routerGroup *VersionedRouterGroup, h *handler.Handlers) {
	setupV1HelloWorldRoutes(routerGroup, h)
	setupV1UserRoutes(routerGroup, h)
}
