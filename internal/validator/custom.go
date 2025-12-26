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
		LocaleZhCN: "{0}必须是有效的手机号码",
		LocaleEnUS: "{0} must be a valid mobile number",
	},
	"username": {
		LocaleZhCN: "{0}必须以字母开头，只能包含字母、数字和下划线",
		LocaleEnUS: "{0} must start with a letter and contain only letters, numbers and underscores",
	},
	"strongpwd": {
		LocaleZhCN: "{0}必须包含大写字母、小写字母和数字",
		LocaleEnUS: "{0} must contain uppercase, lowercase letters and numbers",
	},
	"chinese_name": {
		LocaleZhCN: "{0}必须是有效的中文姓名",
		LocaleEnUS: "{0} must be a valid Chinese name",
	},
	"idcard": {
		LocaleZhCN: "{0}必须是有效的身份证号码",
		LocaleEnUS: "{0} must be a valid ID card number",
	},
}

// registerCustomValidations 注册自定义验证规则
func (v *Validator) registerCustomValidations() {
	// 手机号验证（中国大陆）
	_ = v.validate.RegisterValidation("mobile", validateMobile)

	// 用户名验证（字母开头，只能包含字母数字下划线）
	_ = v.validate.RegisterValidation("username", validateUsername)

	// 强密码验证（至少包含大小写字母和数字）
	_ = v.validate.RegisterValidation("strongpwd", validateStrongPassword)

	// 中文姓名验证
	_ = v.validate.RegisterValidation("chinese_name", validateChineseName)

	// 身份证号验证
	_ = v.validate.RegisterValidation("idcard", validateIDCard)

	// 注册自定义翻译
	v.registerCustomTranslations()
}

// registerCustomTranslations 注册自定义验证规则的翻译
func (v *Validator) registerCustomTranslations() {
	if v.trans == nil {
		return
	}

	for tag, translations := range customTranslations {
		tagCopy := tag
		msg := translations[v.locale]
		if msg == "" {
			msg = translations[LocaleZhCN]
		}
		msgCopy := msg

		_ = v.validate.RegisterTranslation(tagCopy, v.trans, func(trans ut.Translator) error {
			return trans.Add(tagCopy, msgCopy, true)
		}, func(trans ut.Translator, fe validator.FieldError) string {
			translated, _ := trans.T(fe.Tag(), fe.Field())
			return translated
		})
	}
}

// RegisterValidation 注册新的自定义验证规则
func (v *Validator) RegisterValidation(tag string, fn validator.Func, translations map[string]string) error {
	if err := v.validate.RegisterValidation(tag, fn); err != nil {
		return err
	}

	if v.trans == nil {
		return nil
	}

	msg := translations[v.locale]
	if msg == "" {
		msg = translations[LocaleZhCN]
	}

	return v.validate.RegisterTranslation(tag, v.trans, func(trans ut.Translator) error {
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
