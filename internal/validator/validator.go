package validator

import (
	goValidator "github.com/go-playground/validator/v10"
)

type Validate struct {
	*goValidator.Validate
}

var (
	instance *Validate
)

func Instance() *Validate {
	if instance == nil {
		v := goValidator.New()
		RegisterPasswordValidator(v)

		return &Validate{v}
	}

	return instance
}
