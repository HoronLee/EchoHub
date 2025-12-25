package middleware

import (
	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORS 跨域中间件
// 使用 Echo 官方的 CORS 中件，支持可配置的跨域策略
func CORS(cfg *config.AppConfig) echo.MiddlewareFunc {
	// 默认配置
	corsConfig := middleware.CORSConfig{
		// 允许所有来源 (开发环境默认)
		// 生产环境建议通过配置文件限制允许的来源
		AllowOrigins: []string{"*"},
		// 允许的 HTTP 方法
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.PATCH,
			echo.DELETE,
			echo.OPTIONS,
			echo.HEAD,
		},
		// 允许的请求头
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXRequestedWith,
			"X-CSRF-Token",
		},
		// 暴露的响应头 (浏览器可以访问的响应头)
		ExposeHeaders: []string{
			echo.HeaderContentLength,
			echo.HeaderContentEncoding,
			echo.HeaderContentType,
			echo.HeaderAuthorization,
		},
		// 允许发送凭证 (Cookie、Authorization 等)
		// 如果设置为 true，AllowOrigins 不能使用 "*"
		AllowCredentials: false,
		// 预检请求的缓存时间 (秒)
		MaxAge: 86400, // 24 小时
	}

	// 如果配置文件中设置了允许的来源，则使用配置的值
	if len(cfg.CORS.AllowOrigins) > 0 {
		corsConfig.AllowOrigins = cfg.CORS.AllowOrigins
	}

	// 如果配置文件中设置了允许的方法，则使用配置的值
	if len(cfg.CORS.AllowMethods) > 0 {
		corsConfig.AllowMethods = cfg.CORS.AllowMethods
	}

	// 如果配置文件中设置了允许的请求头，则使用配置的值
	if len(cfg.CORS.AllowHeaders) > 0 {
		corsConfig.AllowHeaders = cfg.CORS.AllowHeaders
	}

	// 如果配置文件中设置了暴露的响应头，则使用配置的值
	if len(cfg.CORS.ExposeHeaders) > 0 {
		corsConfig.ExposeHeaders = cfg.CORS.ExposeHeaders
	}

	// 如果配置文件中设置了是否允许凭证，则使用配置的值
	// 注意：当 AllowCredentials 为 true 时，AllowOrigins 不能为 "*"
	if cfg.CORS.AllowCredentials {
		corsConfig.AllowCredentials = true
		// 如果 AllowOrigins 为 "*" 且 AllowCredentials 为 true，需要将 AllowOrigins 设置为具体域名
		if len(corsConfig.AllowOrigins) == 1 && corsConfig.AllowOrigins[0] == "*" {
			// 警告：AllowCredentials 为 true 时，AllowOrigins 不能为 "*"
			// 这里保持默认，实际使用时应该配置具体的域名
			corsConfig.AllowOrigins = []string{"*"} // 实际生产环境应配置具体域名
		}
	}

	// 如果配置文件中设置了预检请求缓存时间，则使用配置的值
	if cfg.CORS.MaxAge > 0 {
		corsConfig.MaxAge = cfg.CORS.MaxAge
	}

	return middleware.CORSWithConfig(corsConfig)
}
