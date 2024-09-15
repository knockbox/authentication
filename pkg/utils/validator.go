package utils

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/knockbox/authentication/pkg/responses"
	"net/http"
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

// DecodeAndValidateStruct decodes and validates the given strct. If we have written a response, this returns true.
func DecodeAndValidateStruct(w http.ResponseWriter, r *http.Request, strct interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(strct); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		msg := "malformed body, expected json"
		responses.NewGenericError(msg).Encode(w)
		return true
	}

	if errs := ValidateStruct(strct); errs != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		
		_ = json.NewEncoder(w).Encode(errs)
		return true
	}

	return false
}
