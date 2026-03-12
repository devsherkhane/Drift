package middleware

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatValidationErrors_MalformedJSON(t *testing.T) {
	// Simulate a non-validator error (e.g., bad JSON syntax)
	err := &json.SyntaxError{Offset: 1}
	errors := FormatValidationErrors(err)

	assert.Len(t, errors, 1)
	assert.Equal(t, "body", errors[0]["field"])
	assert.Contains(t, errors[0]["message"], "Invalid request body")
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"BoardID", "board_i_d"},
		{"Title", "title"},
		{"NewPassword", "new_password"},
		{"Email", "email"},
	}

	for _, tt := range tests {
		result := toSnakeCase(tt.input)
		assert.Equal(t, tt.expected, result, "toSnakeCase(%q)", tt.input)
	}
}
