package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

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

// TranslateErrors 将validator错误转换为友好的错误信息
func TranslateErrors(err error) ValidationErrors {
	if err == nil {
		return nil
	}

	var errors ValidationErrors
	t := GetTranslator()

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			var msg string
			if t != nil {
				msg = e.Translate(t)
			} else {
				msg = e.Error()
			}
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: msg,
			})
		}
	} else {
		errors = append(errors, ValidationError{
			Field:   "",
			Message: err.Error(),
		})
	}

	return errors
}

// FirstErrorMessage 获取第一个错误消息
func FirstErrorMessage(err error) string {
	errors := TranslateErrors(err)
	return errors.First()
}
