package model

// Result 定义统一的API响应格式
type Result[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    T      `json:"data"`
}

const (
	DEFAULT_SUCCESS_CODE = 1
	DEFAULT_FAILED_CODE  = 0
)

// OK 返回成功的结果
func OK[T any](data T, messages ...string) Result[T] {
	// 如果没有传入自定义消息，则使用默认消息
	message := SUCCESS_MESSAGE
	if len(messages) > 0 {
		message = messages[0] // 如果有自定义消息，就使用它
	}

	return Result[T]{
		Code:    DEFAULT_SUCCESS_CODE,
		Message: message,
		Data:    data,
	}
}

// Fail 返回失败的结果
func Fail[T any](message string) Result[T] {
	var zero T
	return Result[T]{
		Code:    DEFAULT_FAILED_CODE,
		Message: message,
		Data:    zero,
	}
}

// OKWithCode 返回成功的结果，并允许自定义状态码
func OKWithCode[T any](data T, code int, messages ...string) Result[T] {
	// 如果没有传入自定义消息，则使用默认消息
	message := SUCCESS_MESSAGE
	if len(messages) > 0 {
		message = messages[0] // 如果有自定义消息，就使用它
	}

	return Result[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
