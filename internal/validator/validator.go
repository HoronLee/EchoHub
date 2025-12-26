package validator

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/labstack/echo/v4"
)

var (
	once     sync.Once
	validate *validator.Validate
	trans    ut.Translator
)

// CustomValidator Echo自定义验证器
type CustomValidator struct {
	validator *validator.Validate
}

// Validate 实现echo.Validator接口
func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

// Init 初始化验证器（单例模式）
func Init() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		// 使用json tag作为字段名
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return fld.Name
			}
			return name
		})

		// 注册中文翻译器
		zhLocale := zh.New()
		uni := ut.New(zhLocale, zhLocale)
		trans, _ = uni.GetTranslator("zh")
		_ = zhTranslations.RegisterDefaultTranslations(validate, trans)

		// 注册自定义验证规则
		registerCustomValidations(validate)
	})
	return validate
}

// GetValidator 获取验证器实例
func GetValidator() *validator.Validate {
	if validate == nil {
		Init()
	}
	return validate
}

// GetTranslator 获取翻译器
func GetTranslator() ut.Translator {
	if trans == nil {
		Init()
	}
	return trans
}

// NewEchoValidator 创建Echo验证器
func NewEchoValidator() *CustomValidator {
	return &CustomValidator{validator: Init()}
}

// SetupEcho 配置Echo实例的验证器
func SetupEcho(e *echo.Echo) {
	e.Validator = NewEchoValidator()
}
