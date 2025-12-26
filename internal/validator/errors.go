package validator

import "strings"

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors 验证错误集合
type ValidationErrors []ValidationError

// Error 实现error接口
func (ve ValidationErrors) Error() string {
	var msgs []string
	for _, e := range ve {
		msgs = append(msgs, e.Message)
	}
	return strings.Join(msgs, "; ")
}

// First 返回第一个错误消息
func (ve ValidationErrors) First() string {
	if len(ve) > 0 {
		return ve[0].Message
	}
	return ""
}
