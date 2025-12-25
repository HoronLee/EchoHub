package middleware

import (
	"net/http"
	"strings"

	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/HoronLee/EchoHub/internal/model/user"
	jwtUtil "github.com/HoronLee/EchoHub/internal/util/jwt"
	"github.com/labstack/echo/v4"
)

func JwtAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// 如果匹配到 /api/v1/* 兜底通配符，说明没有具体路由匹配
			// 跳过认证，让框架返回 404
			// 注意：这里精确匹配，不影响真正的通配符路由（如 /api/v1/files/*）
			path := ctx.Path()
			if path == "" || path == "/api/v1/*" {
				return next(ctx)
			}

			authHeader := ctx.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token not found")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token format invalid")
			}

			tokenString := parts[1]
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token not found")
			}

			jwtService := jwtUtil.NewJWT[user.Claims](&jwtUtil.Config{
				SecretKey: string(config.JWT_SECRET),
			})

			claims, err := jwtService.ParseToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token invalid or expired")
			}

			ctx.Set("user_id", claims.UserID)
			ctx.Set("username", claims.Username)

			return next(ctx)
		}
	}
}
