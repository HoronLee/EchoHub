package middleware

import (
	"net/http"
	"strings"

	"github.com/HoronLee/EchoHub/internal/config"
	commonModel "github.com/HoronLee/EchoHub/internal/model/common"
	"github.com/HoronLee/EchoHub/internal/model/user"
	jwtUtil "github.com/HoronLee/EchoHub/internal/util/jwt"
	"github.com/labstack/echo/v4"
)

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// 从 Authorization Header 提取 Token
			authHeader := ctx.Request().Header.Get("Authorization")
			if authHeader == "" {
				return ctx.JSON(http.StatusUnauthorized,
					commonModel.Fail[string]("Token not found"))
			}

			// 验证 Token 格式（Bearer <token>）
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return ctx.JSON(http.StatusUnauthorized,
					commonModel.Fail[string]("Token format invalid"))
			}

			tokenString := parts[1]
			if tokenString == "" {
				return ctx.JSON(http.StatusUnauthorized,
					commonModel.Fail[string]("Token not found"))
			}

			// 解析和验证 Token
			jwtService := jwtUtil.NewJWT[user.Claims](&jwtUtil.Config{
				SecretKey: string(config.JWT_SECRET),
			})

			claims, err := jwtService.ParseToken(tokenString)
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized,
					commonModel.Fail[string]("Token invalid or expired"))
			}

			// 将 UserID 存入 Context
			ctx.Set("user_id", claims.UserID)
			ctx.Set("username", claims.Username)

			// 也可以使用 jwt 包提供的上下文存储方式
			// ctxReq := jwtUtil.NewContext(ctx.Request().Context(), claims)
			// ctx.SetRequest(ctx.Request().WithContext(ctxReq))

			return next(ctx)
		}
	}
}
