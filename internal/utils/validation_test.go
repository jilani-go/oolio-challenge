package utils

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

// Test struct to validate
type TestStruct struct {
	Name     string `validate:"required"`
	Age      int    `validate:"min=18"`
	Email    string `validate:"required,email"`
	Score    int    `validate:"min=0,max=100"`
	IsActive bool
}

// TestValidateStruct tests the ValidateStruct function
func TestValidateStruct(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		input         TestStruct
		shouldBeValid bool
	}{
		{
			name: "Valid struct - All fields valid",
			input: TestStruct{
				Name:     "John Doe",
				Age:      25,
				Email:    "john@example.com",
				Score:    85,
				IsActive: true,
			},
			shouldBeValid: true,
		},
		{
			name: "Invalid struct - Missing required name",
			input: TestStruct{
				Name:     "",
				Age:      25,
				Email:    "john@example.com",
				Score:    85,
				IsActive: true,
			},
			shouldBeValid: false,
		},
		{
			name: "Invalid struct - Age below minimum",
			input: TestStruct{
				Name:     "John Doe",
				Age:      16,
				Email:    "john@example.com",
				Score:    85,
				IsActive: true,
			},
			shouldBeValid: false,
		},
		{
			name: "Invalid struct - Invalid email format",
			input: TestStruct{
				Name:     "John Doe",
				Age:      25,
				Email:    "invalid-email",
				Score:    85,
				IsActive: true,
			},
			shouldBeValid: false,
		},
		{
			name: "Invalid struct - Score above maximum",
			input: TestStruct{
				Name:     "John Doe",
				Age:      25,
				Email:    "john@example.com",
				Score:    110,
				IsActive: true,
			},
			shouldBeValid: false,
		},
		{
			name: "Invalid struct - Multiple validation errors",
			input: TestStruct{
				Name:     "",
				Age:      16,
				Email:    "invalid-email",
				Score:    110,
				IsActive: true,
			},
			shouldBeValid: false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate the struct
			err := ValidateStruct(tt.input)

			// Assert result
			if tt.shouldBeValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.IsType(t, validator.ValidationErrors{}, err)
			}
		})
	}
}

// TestFormatValidationErrors tests the FormatValidationErrors function
func TestFormatValidationErrors(t *testing.T) {
	// Create a struct with validation errors
	invalidStruct := TestStruct{
		Name:     "",
		Age:      16,
		Email:    "invalid-email",
		Score:    110,
		IsActive: true,
	}

	// Validate and get errors
	err := ValidateStruct(invalidStruct)
	assert.Error(t, err)

	// Format the errors
	formattedErrors := FormatValidationErrors(err)

	// Assert the formatted errors
	assert.Contains(t, formattedErrors, "Name")
	assert.Contains(t, formattedErrors, "Age")
	assert.Contains(t, formattedErrors, "Email")
	assert.Contains(t, formattedErrors, "Score")

	assert.Equal(t, "Name is required", formattedErrors["Name"])
	assert.Equal(t, "Age must be at least 18", formattedErrors["Age"])
	assert.Contains(t, formattedErrors["Email"], "Invalid value for Email")
	assert.Contains(t, formattedErrors["Score"], "Invalid value for Score")

	// Test with non-validation error
	nonValidationErr := errors.New("random error")
	formattedNonValidationErr := FormatValidationErrors(nonValidationErr)
	assert.Equal(t, "random error", formattedNonValidationErr["error"])
}
