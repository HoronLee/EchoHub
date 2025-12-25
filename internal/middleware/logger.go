package middleware

import (
	"time"

	util "github.com/HoronLee/EchoHub/internal/util/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Logger Echo日志中间件
func Logger(logger *util.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()

			err := next(ctx)

			latency := time.Since(start)
			req := ctx.Request()

			logger.Info("HTTP Request",
				zap.Int("status", ctx.Response().Status),
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("query", req.URL.RawQuery),
				zap.String("ip", ctx.RealIP()),
				zap.Duration("latency", latency),
				zap.String("user-agent", req.UserAgent()),
				zap.Error(err),
			)

			return err
		}
	}
}
