package middleware

import (
	"net/http"
	"strings"

	"github.com/HoronLee/EchoHub/internal/config"
	commonModel "github.com/HoronLee/EchoHub/internal/model/common"
	"github.com/HoronLee/EchoHub/internal/model/user"
	jwtUtil "github.com/HoronLee/EchoHub/internal/util/jwt"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization Header 提取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				commonModel.Fail[string]("Token not found"))
			return
		}

		// 验证 Token 格式（Bearer <token>）
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				commonModel.Fail[string]("Token format invalid"))
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				commonModel.Fail[string]("Token not found"))
			return
		}

		// 解析和验证 Token
		jwtService := jwtUtil.NewJWT[user.Claims](&jwtUtil.Config{
			SecretKey: string(config.JWT_SECRET),
		})

		claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				commonModel.Fail[string]("Token invalid or expired"))
			return
		}

		// 将 UserID 存入 Context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// 也可以使用 jwt 包提供的上下文存储方式
		// ctx := jwtUtil.NewContext(c.Request.Context(), claims)
		// c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
