package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
)

type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance with configured field name function
func NewValidator() *CustomValidator {
	validate := validator.New()

	// Register function to get field name from json tag
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})

	return &CustomValidator{
		validator: validate,
	}
}

// Validate performs validation and returns formatted errors
func (cv *CustomValidator) Validate(s interface{}) error {
	if err := cv.validator.Struct(s); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		if len(validationErrors) > 0 {
			errorMap := make(map[string]string)

			// Get first error for each field
			for _, e := range validationErrors {
				fieldName := e.Field()
				if _, exists := errorMap[fieldName]; !exists {
					errorMap[fieldName] = getErrorMsg(e)
				}
			}

			// Construct error message
			var errMsgs []string
			for field, msg := range errorMap {
				// Fix: Use constant format string
				errorMsg := fmt.Sprintf("%s: %s", field, msg)
				errMsgs = append(errMsgs, errorMsg)
			}

			return errors.New(strings.Join(errMsgs, "; "))
		}
	}
	return nil
}

// Helper function to generate friendly error messages
func getErrorMsg(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("should be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("should not be longer than %s characters", e.Param())
	case "gt":
		return fmt.Sprintf("should be greater than %s", e.Param())
	case "gte":
		return fmt.Sprintf("should be greater than or equal to %s", e.Param())
	default:
		return fmt.Sprintf("failed validation on %s", e.Tag())
	}
}
