package model

// ServerError 定义服务器错误信息
type ServerError struct {
	Msg string
	Err error
}
