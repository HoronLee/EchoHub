package response

import (
	"net/http"

	commonModel "github.com/HoronLee/EchoHub/internal/model/common"
	errorUtil "github.com/HoronLee/EchoHub/internal/util/err"
	"github.com/labstack/echo/v4"
)

// Response 代表 handler 层的执行结果封装
// swagger:model Response
type Response struct {
	// Code 状态码，非0时表示自定义HTTP业务状态码
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
		if res.Err != nil {
			return ctx.JSON(http.StatusBadRequest, commonModel.Fail[string](
				errorUtil.HandleError(&commonModel.ServerError{
					Msg: res.Msg,
					Err: res.Err,
				}),
			))
		}

		// 支持自定义 code
		if res.Code != 0 {
			return ctx.JSON(http.StatusOK, commonModel.OKWithCode(res.Data, res.Code, res.Msg))
		} else {
			return ctx.JSON(http.StatusOK, commonModel.OK(res.Data, res.Msg))
		}
	}
}
