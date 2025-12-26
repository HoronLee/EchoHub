package response

import (
	"net/http"

	log "github.com/HoronLee/EchoHub/internal/util/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Response 代表 handler 层的执行结果封装
// swagger:model Response
type Response struct {
	// HTTPStatus HTTP状态码，默认为200
	HTTPStatus int `json:"-"`

	// Code 业务状态码，0表示成功，非0表示业务错误
	Code int `json:"code" example:"0" description:"状态码，0表示成功，非0表示自定义业务状态码"`

	// Data 响应数据，具体内容因接口而异
	Data any `json:"data,omitempty" description:"响应数据，具体内容因接口而异"`

	// Msg 返回信息，通常是状态描述
	Msg string `json:"msg" example:"success" description:"返回信息，通常是状态描述"`

	// Err 错误信息，序列化时忽略（仅供内部日志使用）
	Err error `json:"-"`
}

// Execute 包装器，自动根据 Response 返回统一格式的 HTTP 响应 (仅处理返回类型为JSON的handler)
func Execute(fn func(ctx echo.Context) Response) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		res := fn(ctx)

		if res.HTTPStatus == 0 {
			res.HTTPStatus = http.StatusOK
		}

		if res.Err != nil {
			log.GetLogger().Error("Business error",
				zap.String("path", ctx.Path()),
				zap.String("method", ctx.Request().Method),
				zap.Int("http_status", res.HTTPStatus),
				zap.Int("business_code", res.Code),
				zap.String("message", res.Msg),
				zap.Error(res.Err),
			)
		}

		return ctx.JSON(res.HTTPStatus, res)
	}
}

func Success(data any, msg ...string) Response {
	message := "success"
	if len(msg) > 0 {
		message = msg[0]
	}
	return Response{
		HTTPStatus: http.StatusOK,
		Code:       0,
		Data:       data,
		Msg:        message,
	}
}

func Error(httpStatus int, code int, msg string, err ...error) Response {
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	return Response{
		HTTPStatus: httpStatus,
		Code:       code,
		Msg:        msg,
		Err:        e,
	}
}

func BadRequest(msg string, err ...error) Response {
	return Error(http.StatusBadRequest, 400, msg, err...)
}

// ValidationError 验证错误响应
func ValidationError(msg string, err ...error) Response {
	return Error(http.StatusUnprocessableEntity, 422, msg, err...)
}

func Unauthorized(msg string, err ...error) Response {
	return Error(http.StatusUnauthorized, 401, msg, err...)
}

func Forbidden(msg string, err ...error) Response {
	return Error(http.StatusForbidden, 403, msg, err...)
}

func NotFound(msg string, err ...error) Response {
	return Error(http.StatusNotFound, 404, msg, err...)
}

func InternalServerError(msg string, err ...error) Response {
	return Error(http.StatusInternalServerError, 500, msg, err...)
}
