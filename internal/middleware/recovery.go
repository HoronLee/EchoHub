package middleware

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/HoronLee/EchoHub/internal/response"
	util "github.com/HoronLee/EchoHub/internal/util/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Recovery 自定义 panic 恢复中间件
func Recovery(logger *util.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					// 检测是否为客户端断开连接
					var brokenPipe bool
					if ne, ok := r.(*net.OpError); ok {
						if se, ok := ne.Err.(*os.SyscallError); ok {
							errStr := strings.ToLower(se.Error())
							if strings.Contains(errStr, "broken pipe") ||
								strings.Contains(errStr, "connection reset by peer") {
								brokenPipe = true
							}
						}
					}

					// 客户端断开连接：只记录简单日志
					if brokenPipe {
						logger.Error("Client disconnected",
							zap.String("path", ctx.Request().URL.Path),
							zap.Any("error", r),
						)
						err = http.ErrAbortHandler
						return
					}

					// Panic 错误：记录详细信息
					logger.Error("Panic recovered",
						zap.Any("error", r),
						zap.String("method", ctx.Request().Method),
						zap.String("path", ctx.Request().URL.Path),
						zap.String("ip", ctx.RealIP()),
						zap.Stack("stacktrace"),
					)

					err = ctx.JSON(http.StatusInternalServerError, response.InternalServerError("Internal server error"))
				}
			}()
			return next(ctx)
		}
	}
}
