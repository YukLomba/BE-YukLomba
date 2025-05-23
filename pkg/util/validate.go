package util

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func ValidateStruct(s any) error {
	validate = validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(s)
}
