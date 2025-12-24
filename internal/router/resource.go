package router

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/HoronLee/EchoHub/internal/handler"
	_ "github.com/HoronLee/EchoHub/internal/swagger"
)

// setupResourceRoutes 设置资源路由
func setupResourceRoutes(routerGroup *VersionedRouterGroup, _ *handler.Handlers) {
	// Swagger UI - 使用公共路由组，无需认证
	routerGroup.PublicRouterGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
