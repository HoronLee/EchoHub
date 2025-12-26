package validator

import (
	"reflect"
	"strings"
	"sync"

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

var (
	once     sync.Once
	validate *validator.Validate
	uni      *ut.UniversalTranslator
	trans    ut.Translator
	locale   string = LocaleZhCN // 默认中文
)

// CustomValidator Echo自定义验证器
type CustomValidator struct {
	validator *validator.Validate
}

// Validate 实现echo.Validator接口
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Init 初始化验证器（单例模式）
// loc: 语言设置，支持 "zh_CN" 和 "en_US"，默认 "zh_CN"
func Init(loc ...string) *validator.Validate {
	once.Do(func() {
		// 设置语言
		if len(loc) > 0 && loc[0] != "" {
			locale = loc[0]
		}

		validate = validator.New()

		// 使用json tag作为字段名
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return fld.Name
			}
			return name
		})

		// 初始化多语言翻译器
		initTranslator()

		// 注册自定义验证规则
		registerCustomValidations(validate)
	})
	return validate
}

// initTranslator 初始化翻译器
func initTranslator() {
	zhLocale := zh.New()
	enLocale := en.New()
	uni = ut.New(enLocale, zhLocale, enLocale)

	switch locale {
	case LocaleEnUS:
		trans, _ = uni.GetTranslator("en")
		_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	default: // 默认中文
		trans, _ = uni.GetTranslator("zh")
		_ = zhTranslations.RegisterDefaultTranslations(validate, trans)
	}
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

// GetLocale 获取当前语言设置
func GetLocale() string {
	return locale
}

// NewEchoValidator 创建Echo验证器
func NewEchoValidator() *CustomValidator {
	return &CustomValidator{validator: GetValidator()}
}

// SetupEcho 配置Echo实例的验证器
// loc: 语言设置，支持 "zh_CN" 和 "en_US"
func SetupEcho(e *echo.Echo, loc ...string) {
	if len(loc) > 0 {
		Init(loc[0])
	} else {
		Init()
	}
	e.Validator = NewEchoValidator()
}
