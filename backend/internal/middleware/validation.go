package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// FormatValidationErrors converts Gin binding errors into user-friendly messages.
// This replaces the hardcoded error strings scattered across handlers with a
// single, consistent format.
//
// Usage in any handler:
//
//	if err := c.ShouldBindJSON(&input); err != nil {
//	    c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
//	    return
//	}
func FormatValidationErrors(err error) []gin.H {
	var errors []gin.H

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := toSnakeCase(e.Field())
			errors = append(errors, gin.H{
				"field":   field,
				"message": buildMessage(e),
			})
		}
	} else {
		// Fallback for non-validation errors (e.g., malformed JSON)
		errors = append(errors, gin.H{
			"field":   "body",
			"message": "Invalid request body — check your JSON syntax",
		})
	}

	return errors
}

// buildMessage converts a validator.FieldError into a human-readable string
func buildMessage(e validator.FieldError) string {
	field := toSnakeCase(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	default:
		return fmt.Sprintf("%s failed validation: %s", field, e.Tag())
	}
}

// toSnakeCase converts PascalCase field names (e.g., "BoardID") to snake_case ("board_id").
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(r + 32) // lowercase
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
