package router

import "github.com/HoronLee/EchoHub/internal/handler"

// setupV1UserRoutes 设置 v1 版本的用户路由
func setupV1UserRoutes(routerGroup *VersionedRouterGroup, h *handler.Handlers) {
	// Public routes - 公开路由，无需认证
	// 路径: POST /api/v1/register, POST /api/v1/login
	routerGroup.PublicRouterGroup.POST("/register", h.UserHandler.Register())
	routerGroup.PublicRouterGroup.POST("/login", h.UserHandler.Login())

	// Private routes - 私有路由，需要 JWT 认证
	// 路径: DELETE /api/v1/user
	routerGroup.PrivateRouterGroup.DELETE("/user", h.UserHandler.DeleteUser())
}
