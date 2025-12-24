package router

import "github.com/HoronLee/EchoHub/internal/handler"

// setupV1HelloWorldRoutes 设置 v1 版本的 HelloWorld 路由
func setupV1HelloWorldRoutes(routerGroup *VersionedRouterGroup, h *handler.Handlers) {
	// Public routes - 公开路由，无需认证
	// 路径: POST /api/v1/helloworld
	routerGroup.PublicRouterGroup.POST("/helloworld", h.HelloWorldHandler.PostHelloWorld())

	// Private routes - 私有路由，需要 JWT 认证
	// 可以在这里添加需要认证的 HelloWorld 相关路由
}
