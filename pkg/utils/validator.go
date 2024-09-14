package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/knockbox/authentication/pkg/responses"
)

var validate = validator.New()

// ValidateStruct validates the given struct and returns the errors, if any.
func ValidateStruct(strct interface{}) []*responses.ValidationError {
	var errors []*responses.ValidationError

	err := validate.Struct(strct)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el responses.ValidationError

			el.FailedField = err.StructNamespace()
			el.Tag = err.Tag()
			el.Value = err.Param()

			errors = append(errors, &el)
		}
	}

	return errors
}
