package helloworld

// CreateRequest 创建 HelloWorld 的请求
// swagger:model CreateRequest
type CreateRequest struct {
	Message string `json:"message" binding:"required,min=1,max=20" example:"Hello, World!" description:"问候消息，长度1-20字符"`
}

// CreateResponse 创建 HelloWorld 的响应
// swagger:model CreateResponse
type CreateResponse struct {
	Message  string `json:"message" example:"Hello, World!" description:"返回的问候消息"`
	Version  string `json:"version" example:"1.0.0" description:"应用程序版本"`
	Database string `json:"database" example:"SQLite" description:"数据库类型"`
}
