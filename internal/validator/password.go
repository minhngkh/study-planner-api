package validator

import (
	"fmt"

	goValidator "github.com/go-playground/validator/v10"
)

var (
	MinLength = 6
	MaxLength = 30
)

func PasswordValidator(fl goValidator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	for _, c := range password {
		// ! -> 33, ~ -> 126
		if c < 33 || c > 126 {
			return false
		}
	}

	return true
}

func RegisterPasswordValidator(v *goValidator.Validate) {
	v.RegisterValidation("password", PasswordValidator)
	v.RegisterAlias("default-password", fmt.Sprintf("password,min=%d,max=%d", MinLength, MaxLength))
}
