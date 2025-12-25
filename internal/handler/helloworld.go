package handler

import (
	commonModel "github.com/HoronLee/EchoHub/internal/model/common"
	"github.com/HoronLee/EchoHub/internal/model/helloworld"
	res "github.com/HoronLee/EchoHub/internal/response"
	"github.com/HoronLee/EchoHub/internal/service"
	"github.com/labstack/echo/v4"
)

type HelloWorldHandler struct {
	svc *service.HelloWorldService
}

func NewHelloWorldHandler(svc *service.HelloWorldService) *HelloWorldHandler {
	return &HelloWorldHandler{svc: svc}
}

// PostHelloWorld 处理POST /helloworld请求
// @Summary 创建HelloWorld消息
// @Description 创建一个新的HelloWorld消息并返回系统信息
// @Tags HelloWorld
// @Accept json
// @Produce json
// @Param request body helloworld.CreateRequest true "HelloWorld创建请求参数"
// @Success 200 {object} helloworld.CreateResponse "创建成功，返回消息和系统信息"
// @Failure 400 {object} res.Response "请求参数错误或创建失败"
// @Router /helloworld [post]
func (h *HelloWorldHandler) PostHelloWorld() echo.HandlerFunc {
	return res.Execute(func(ctx echo.Context) res.Response {
		var req helloworld.CreateRequest
		if err := ctx.Bind(&req); err != nil {
			return res.Response{Msg: "", Err: err} // Msg 为空，自动使用 err 详情
		}

		if err := h.svc.PostHelloWorld(ctx.Request().Context(), req.Message); err != nil {
			return res.Response{Msg: "Failed to create hello world", Err: err}
		}

		dbInfo, err := h.svc.GetDatabaseInfo(ctx.Request().Context())
		if err != nil {
			return res.Response{Msg: "Failed to get database info", Err: err}
		}

		return res.Response{
			Data: helloworld.CreateResponse{
				Message:  req.Message,
				Version:  commonModel.Version,
				Database: dbInfo,
			},
			Msg: "success",
		}
	})
}
