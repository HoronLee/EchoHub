package router

import (
	"github.com/HoronLee/EchoHub/internal/handler"
	_ "github.com/HoronLee/EchoHub/internal/swagger"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// setupResourceRoutes 设置资源路由
func setupResourceRoutes(routerGroup *VersionedRouterGroup, _ *handler.Handlers) {
	// Swagger UI - 使用公共路由组，无需认证
	// 使用 Any 方法处理所有 HTTP 方法
	routerGroup.PublicRouter.Any("/swagger/*", echoSwagger.WrapHandler)
}
