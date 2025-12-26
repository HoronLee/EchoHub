package validator

import (
	"regexp"
	"unicode"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// 自定义验证规则的多语言翻译
var customTranslations = map[string]map[string]string{
	"mobile": {
		"zh_CN": "{0}必须是有效的手机号码",
		"en_US": "{0} must be a valid mobile number",
	},
	"username": {
		"zh_CN": "{0}必须以字母开头，只能包含字母、数字和下划线",
		"en_US": "{0} must start with a letter and contain only letters, numbers and underscores",
	},
	"strongpwd": {
		"zh_CN": "{0}必须包含大写字母、小写字母和数字",
		"en_US": "{0} must contain uppercase, lowercase letters and numbers",
	},
	"chinese_name": {
		"zh_CN": "{0}必须是有效的中文姓名",
		"en_US": "{0} must be a valid Chinese name",
	},
	"idcard": {
		"zh_CN": "{0}必须是有效的身份证号码",
		"en_US": "{0} must be a valid ID card number",
	},
}

// registerCustomValidations 注册自定义验证规则
func registerCustomValidations(v *validator.Validate) {
	// 手机号验证（中国大陆）
	_ = v.RegisterValidation("mobile", validateMobile)

	// 用户名验证（字母开头，只能包含字母数字下划线）
	_ = v.RegisterValidation("username", validateUsername)

	// 强密码验证（至少包含大小写字母和数字）
	_ = v.RegisterValidation("strongpwd", validateStrongPassword)

	// 中文姓名验证
	_ = v.RegisterValidation("chinese_name", validateChineseName)

	// 身份证号验证
	_ = v.RegisterValidation("idcard", validateIDCard)

	// 注册自定义翻译
	registerCustomTranslations(v)
}

// registerCustomTranslations 注册自定义验证规则的翻译
func registerCustomTranslations(v *validator.Validate) {
	t := GetTranslator()
	if t == nil {
		return
	}

	currentLocale := GetLocale()

	for tag, translations := range customTranslations {
		tagCopy := tag
		msg := translations[currentLocale]
		if msg == "" {
			msg = translations[LocaleZhCN] // 回退到中文
		}
		msgCopy := msg

		_ = v.RegisterTranslation(tagCopy, t, func(trans ut.Translator) error {
			return trans.Add(tagCopy, msgCopy, true)
		}, func(trans ut.Translator, fe validator.FieldError) string {
			translated, _ := trans.T(fe.Tag(), fe.Field())
			return translated
		})
	}
}

// RegisterCustomValidation 注册新的自定义验证规则
// tag: 规则名称
// fn: 验证函数
// translations: 多语言翻译 map[locale]message
func RegisterCustomValidation(tag string, fn validator.Func, translations map[string]string) error {
	v := GetValidator()
	if err := v.RegisterValidation(tag, fn); err != nil {
		return err
	}

	// 注册翻译
	t := GetTranslator()
	if t == nil {
		return nil
	}

	currentLocale := GetLocale()
	msg := translations[currentLocale]
	if msg == "" {
		msg = translations[LocaleZhCN]
	}

	return v.RegisterTranslation(tag, t, func(trans ut.Translator) error {
		return trans.Add(tag, msg, true)
	}, func(trans ut.Translator, fe validator.FieldError) string {
		translated, _ := trans.T(fe.Tag(), fe.Field())
		return translated
	})
}

// validateMobile 验证手机号
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, mobile)
	return matched
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	pattern := `^[a-zA-Z][a-zA-Z0-9_]*$`
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

// validateStrongPassword 验证强密码
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var hasUpper, hasLower, hasDigit bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

// validateChineseName 验证中文姓名
func validateChineseName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	for _, r := range name {
		if !unicode.Is(unicode.Han, r) {
			return false
		}
	}
	return len(name) >= 2
}

// validateIDCard 验证身份证号（简单验证）
func validateIDCard(fl validator.FieldLevel) bool {
	idcard := fl.Field().String()
	pattern := `^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`
	matched, _ := regexp.MatchString(pattern, idcard)
	return matched
}
