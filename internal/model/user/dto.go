package user

// RegisterRequest 注册请求
// swagger:model RegisterRequest
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"john_doe" description:"用户名，长度3-50字符"`
	Password string `json:"password" binding:"required,min=6" example:"password123" description:"密码，最少6个字符"`
}

// LoginRequest 登录请求
// swagger:model LoginRequest
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe" description:"用户名"`
	Password string `json:"password" binding:"required" example:"password123" description:"密码"`
}

// LoginResponse 登录响应
// swagger:model LoginResponse
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"JWT访问令牌"`
}
