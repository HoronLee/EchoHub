package handler

import (
	"github.com/HoronLee/EchoHub/internal/model/user"
	res "github.com/HoronLee/EchoHub/internal/response"
	"github.com/HoronLee/EchoHub/internal/service"
	"github.com/labstack/echo/v4"
)

// UserHandler 用户处理器
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler 创建UserHandler实例
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

// Register 用户注册处理器
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body user.RegisterRequest true "注册请求参数"
// @Success 200 {object} map[string]string "注册成功"
// @Failure 400 {object} res.Response "请求参数错误或注册失败"
// @Router /v1/register [post]
func (h *UserHandler) Register() echo.HandlerFunc {
	return res.Execute(func(ctx echo.Context) res.Response {
		var req user.RegisterRequest
		if err := ctx.Bind(&req); err != nil {
			return res.BadRequest("Invalid request body", err)
		}

		if err := h.svc.Register(ctx.Request().Context(), req); err != nil {
			return res.InternalServerError("Registration failed", err)
		}

		return res.Success(map[string]any{"message": "User registered successfully"}, "success")
	})
}

// Login 用户登录处理器
// @Summary 用户登录
// @Description 用户身份验证并获取访问令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body user.LoginRequest true "登录请求参数"
// @Success 200 {object} user.LoginResponse "登录成功，返回JWT令牌"
// @Failure 400 {object} res.Response "请求参数错误或登录失败"
// @Router /v1/login [post]
func (h *UserHandler) Login() echo.HandlerFunc {
	return res.Execute(func(ctx echo.Context) res.Response {
		var req user.LoginRequest
		if err := ctx.Bind(&req); err != nil {
			return res.BadRequest("Invalid request body", err)
		}

		token, err := h.svc.Login(ctx.Request().Context(), req)
		if err != nil {
			return res.Unauthorized("Login failed", err)
		}

		return res.Success(user.LoginResponse{Token: token}, "success")
	})
}

// DeleteUser 删除用户处理器
// @Summary 删除用户
// @Description 删除当前登录的用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "删除成功"
// @Failure 400 {object} res.Response "用户未认证或删除失败"
// @Failure 401 {object} res.Response "用户未认证"
// @Router /v1/user [delete]
func (h *UserHandler) DeleteUser() echo.HandlerFunc {
	return res.Execute(func(ctx echo.Context) res.Response {
		userIDValue := ctx.Get("user_id")
		if userIDValue == nil {
			return res.Unauthorized("User not authenticated")
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			return res.BadRequest("Invalid user ID format")
		}

		if err := h.svc.DeleteUser(ctx.Request().Context(), userID); err != nil {
			return res.InternalServerError("Failed to delete user", err)
		}

		return res.Success(map[string]any{"message": "User deleted successfully"}, "success")
	})
}
