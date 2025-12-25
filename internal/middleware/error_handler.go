package middleware

import (
	"net/http"

	"github.com/HoronLee/EchoHub/internal/response"
	"github.com/HoronLee/EchoHub/internal/util/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CustomHTTPErrorHandler 自定义HTTP错误处理器
// 将所有错误（包括404、405、500等）统一转换为项目响应格式
func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if he.Message != nil {
			message = he.Message.(string)
		}
	} else {
		message = err.Error()
	}

	util.GetLogger().Error("HTTP error",
		zap.String("path", c.Path()),
		zap.String("method", c.Request().Method),
		zap.Int("http_status", code),
		zap.String("message", message),
		zap.Error(err),
	)

	c.JSON(code, response.Error(code, code, message))
}
