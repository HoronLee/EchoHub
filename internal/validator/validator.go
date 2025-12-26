package validator

import (
	"reflect"
	"strings"

	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/labstack/echo/v4"
)

// 支持的语言
const (
	LocaleZhCN = "zh_CN"
	LocaleEnUS = "en_US"
)

// Validator 验证器结构体
type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
	locale   string
}

// NewValidator 创建验证器实例
func NewValidator(cfg *config.AppConfig) *Validator {
	v := &Validator{
		validate: validator.New(),
		locale:   cfg.Server.Locale,
	}

	if v.locale == "" {
		v.locale = LocaleZhCN
	}

	// 使用json tag作为字段名
	v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})

	// 初始化翻译器
	v.initTranslator()

	// 注册自定义验证规则
	v.registerCustomValidations()

	return v
}

// initTranslator 初始化翻译器
func (v *Validator) initTranslator() {
	zhLocale := zh.New()
	enLocale := en.New()
	uni := ut.New(enLocale, zhLocale, enLocale)

	switch v.locale {
	case LocaleEnUS:
		v.trans, _ = uni.GetTranslator("en")
		_ = enTranslations.RegisterDefaultTranslations(v.validate, v.trans)
	default:
		v.trans, _ = uni.GetTranslator("zh")
		_ = zhTranslations.RegisterDefaultTranslations(v.validate, v.trans)
	}
}

// Validate 实现 echo.Validator 接口
func (v *Validator) Validate(i any) error {
	return v.validate.Struct(i)
}

// ValidateStruct 验证结构体
func (v *Validator) ValidateStruct(i any) error {
	return v.validate.Struct(i)
}

// ValidateVar 验证单个变量
func (v *Validator) ValidateVar(field any, tag string) error {
	return v.validate.Var(field, tag)
}

// BindAndValidate 绑定并验证请求数据
func (v *Validator) BindAndValidate(c echo.Context, i any) (bool, string) {
	if err := c.Bind(i); err != nil {
		return false, v.getBindErrorMessage()
	}
	if err := v.validate.Struct(i); err != nil {
		return false, v.FirstErrorMessage(err)
	}
	return true, ""
}

// TranslateErrors 将验证错误转换为友好的错误信息
func (v *Validator) TranslateErrors(err error) ValidationErrors {
	if err == nil {
		return nil
	}

	var errors ValidationErrors

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			var msg string
			if v.trans != nil {
				msg = e.Translate(v.trans)
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
func (v *Validator) FirstErrorMessage(err error) string {
	errors := v.TranslateErrors(err)
	if len(errors) > 0 {
		return errors[0].Message
	}
	return ""
}

// GetLocale 获取当前语言设置
func (v *Validator) GetLocale() string {
	return v.locale
}

// getBindErrorMessage 获取绑定错误消息
func (v *Validator) getBindErrorMessage() string {
	if v.locale == LocaleEnUS {
		return "Invalid request format"
	}
	return "请求数据格式错误"
}

// SetupEcho 配置Echo实例的验证器
func (v *Validator) SetupEcho(e *echo.Echo) {
	e.Validator = v
}
