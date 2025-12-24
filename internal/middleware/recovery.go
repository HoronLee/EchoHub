package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	commonModel "github.com/HoronLee/EchoHub/internal/model/common"
	util "github.com/HoronLee/EchoHub/internal/util/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 自定义 panic 恢复中间件
func Recovery(logger *util.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// 检测是否为客户端断开连接
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
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
						zap.String("path", c.Request.URL.Path),
						zap.Any("error", err),
					)
					c.Abort()
					return
				}

				// Panic 错误：记录详细信息
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("ip", c.ClientIP()),
					zap.ByteString("request", httpRequest),
					zap.Stack("stacktrace"),
				)

				// 使用统一的响应格式
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					commonModel.Fail[string]("Internal server error"))
			}
		}()
		c.Next()
	}
}
