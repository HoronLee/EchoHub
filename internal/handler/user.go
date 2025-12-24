package handler

import (
	"github.com/HoronLee/EchoHub/internal/model/user"
	res "github.com/HoronLee/EchoHub/internal/response"
	"github.com/HoronLee/EchoHub/internal/service"
	"github.com/gin-gonic/gin"
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
// @Success 200 {object} response.Response{data=map[string]string} "注册成功"
// @Failure 400 {object} response.Response "请求参数错误或注册失败"
// @Router /user/register [post]
func (h *UserHandler) Register() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		var req user.RegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			return res.Response{Msg: "Invalid request body", Err: err}
		}

		if err := h.svc.Register(ctx.Request.Context(), req); err != nil {
			return res.Response{Msg: "Registration failed", Err: err}
		}

		return res.Response{
			Data: gin.H{"message": "User registered successfully"},
			Msg:  "success",
		}
	})
}

// Login 用户登录处理器
// @Summary 用户登录
// @Description 用户身份验证并获取访问令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body user.LoginRequest true "登录请求参数"
// @Success 200 {object} response.Response{data=user.LoginResponse} "登录成功，返回JWT令牌"
// @Failure 400 {object} response.Response "请求参数错误或登录失败"
// @Router /user/login [post]
func (h *UserHandler) Login() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		var req user.LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			return res.Response{Msg: "Invalid request body", Err: err}
		}

		token, err := h.svc.Login(ctx.Request.Context(), req)
		if err != nil {
			return res.Response{Msg: "Login failed", Err: err}
		}

		return res.Response{
			Data: user.LoginResponse{Token: token},
			Msg:  "success",
		}
	})
}

// DeleteUser 删除用户处理器
// @Summary 删除用户
// @Description 删除当前登录的用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=map[string]string} "删除成功"
// @Failure 400 {object} response.Response "用户未认证或删除失败"
// @Failure 401 {object} response.Response "用户未认证"
// @Router /user/delete [delete]
func (h *UserHandler) DeleteUser() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 从JWT中间件获取用户ID（当前登录用户）
		userIDValue, exists := ctx.Get("user_id")
		if !exists {
			return res.Response{Msg: "User not authenticated", Err: nil}
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			return res.Response{Msg: "Invalid user ID format", Err: nil}
		}

		// 也可以从URL参数获取要删除的用户ID（如果需要管理员删除其他用户）
		// 这里简化为删除当前登录用户
		if err := h.svc.DeleteUser(ctx.Request.Context(), userID); err != nil {
			return res.Response{Msg: "Failed to delete user", Err: err}
		}

		return res.Response{
			Data: gin.H{"message": "User deleted successfully"},
			Msg:  "success",
		}
	})
}
