package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Initialize the validator instance
func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using validator tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// FormatValidationErrors converts validator errors to a readable format
func FormatValidationErrors(err error) map[string]string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		errorMap := make(map[string]string)
		for _, e := range validationErrs {
			// Convert the field name to JSON field name if needed
			fieldName := e.Field()
			
			// Create a user-friendly error message
			switch e.Tag() {
			case "required":
				errorMap[fieldName] = fieldName + " is required"
			case "min":
				errorMap[fieldName] = fieldName + " must be at least " + e.Param()
			case "gtefield":
				errorMap[fieldName] = fieldName + " must be after " + e.Param()
			default:
				errorMap[fieldName] = "Invalid value for " + fieldName
			}
		}
		return errorMap
	}
	return map[string]string{"error": err.Error()}
}
