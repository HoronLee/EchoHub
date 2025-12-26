package validator

import (
	"github.com/labstack/echo/v4"
)

// BindAndValidate 绑定并验证请求数据
// 返回值: 是否成功, 错误信息
func BindAndValidate(c echo.Context, i any) (bool, string) {
	// 绑定请求数据
	if err := c.Bind(i); err != nil {
		return false, "请求数据格式错误"
	}

	// 验证数据
	if err := Validate(i); err != nil {
		return false, FirstErrorMessage(err)
	}

	return true, ""
}

// Validate 验证结构体
func Validate(i any) error {
	v := GetValidator()
	return v.Struct(i)
}

// ValidateVar 验证单个变量
func ValidateVar(field any, tag string) error {
	v := GetValidator()
	return v.Var(field, tag)
}

// ValidateMap 验证map数据
func ValidateMap(data map[string]any, rules map[string]any) map[string]any {
	v := GetValidator()
	return v.ValidateMap(data, rules)
}
