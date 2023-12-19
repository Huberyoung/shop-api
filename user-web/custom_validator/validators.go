package custom_validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	matched, err := regexp.MatchString(`^1(3\d{1}|4[57]|5[012356789]|6[6]|7[0135678]|8\d{1})\d{8}$`, mobile)
	if err != nil || !matched {
		return false
	}
	return true
}
