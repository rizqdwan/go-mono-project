package validator

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	v *validator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{v: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.v.Struct(i); err != nil {
		return err
	}
	return nil
}
