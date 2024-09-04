package response

import (
	"errors"
	"thelist-app/gokit/translation"

	"github.com/go-playground/validator/v10"
)

type Validations map[string][]string

type ValidationError struct {
	Property string `json:"property"`
	Message  string `json:"message"`
}

type validation struct {
	Translation translation.Translation
}

func newValidationTranslator(
	t translation.Translation,
) *validation {
	return &validation{
		Translation: t,
	}
}

func (v validation) translate(err error) Validations {
	var validationError []ValidationError

	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		return nil
	}

	for _, err := range err.(validator.ValidationErrors) { //nolint: all
		validationError = append(validationError, ValidationError{
			Property: err.Field(),
			Message: func() string {
				return v.Translation.Trans(
					"validation."+err.Tag(),
					map[string]interface{}{
						"attribute": v.Translation.Trans("attributes."+err.Field(), nil),
						err.Tag():   err.Param(),
					},
				)
			}(),
		})
	}

	validations := Validations{}
	for _, err := range validationError {
		if _, ok := validations[err.Property]; ok { //nolint: gosimple
			validations[err.Property] = append(validations[err.Property], err.Message)
		} else {
			validations[err.Property] = []string{err.Message}
		}
	}

	return validations
}
