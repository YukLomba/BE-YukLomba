package util

import (
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/common"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateStruct(s any) error {
	validate = validator.New(validator.WithRequiredStructEnabled())
	RegisterDatetimeValidators(validate)
	return validate.Struct(s)
}

func GenerateValidationErrorMessage(err error) []string {
	var errors []string

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()

		// Create user-friendly messages
		switch tag {
		case "required":
			errors = append(errors, field+" is required")
		case "email":
			errors = append(errors, field+" must be a valid email address")
		case "min":
			errors = append(errors, field+" must be at least "+err.Param()+" characters")
		case "max":
			errors = append(errors, field+" cannot exceed "+err.Param()+" characters")
		default:
			errors = append(errors, field+" is invalid")
		}
	}

	return errors
}

func RegisterDatetimeValidators(v *validator.Validate) {
	v.RegisterValidation("future", func(fl validator.FieldLevel) bool {
		dt, ok := fl.Field().Interface().(common.Datetime)
		if !ok {
			return false
		}
		return dt.ToTime().After(time.Now())
	})

	v.RegisterValidation("past", func(fl validator.FieldLevel) bool {
		dt, ok := fl.Field().Interface().(common.Datetime)
		if !ok {
			return false
		}
		return dt.ToTime().Before(time.Now())
	})
}
