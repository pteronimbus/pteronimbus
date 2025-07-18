package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringArray_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected StringArray
		hasError bool
	}{
		{
			name:     "empty array",
			input:    "{}",
			expected: StringArray{},
			hasError: false,
		},
		{
			name:     "single item",
			input:    "{owner}",
			expected: StringArray{"owner"},
			hasError: false,
		},
		{
			name:     "multiple items",
			input:    "{owner,admin,moderator}",
			expected: StringArray{"owner", "admin", "moderator"},
			hasError: false,
		},
		{
			name:     "quoted items",
			input:    `{"owner","admin","moderator"}`,
			expected: StringArray{"owner", "admin", "moderator"},
			hasError: false,
		},
		{
			name:     "mixed quoted and unquoted",
			input:    `{owner,"admin",moderator}`,
			expected: StringArray{"owner", "admin", "moderator"},
			hasError: false,
		},
		{
			name:     "items with spaces",
			input:    `{"server admin","game moderator"}`,
			expected: StringArray{"server admin", "game moderator"},
			hasError: false,
		},
		{
			name:     "nil input",
			input:    nil,
			expected: StringArray{},
			hasError: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: StringArray{},
			hasError: false,
		},
		{
			name:     "byte slice input",
			input:    []byte("{owner,admin}"),
			expected: StringArray{"owner", "admin"},
			hasError: false,
		},
		{
			name:     "invalid format",
			input:    "invalid",
			expected: StringArray{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sa StringArray
			err := sa.Scan(tt.input)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, sa)
			}
		})
	}
}

func TestStringArray_Value(t *testing.T) {
	tests := []struct {
		name     string
		input    StringArray
		expected string
	}{
		{
			name:     "empty array",
			input:    StringArray{},
			expected: "{}",
		},
		{
			name:     "single item",
			input:    StringArray{"owner"},
			expected: `{"owner"}`,
		},
		{
			name:     "multiple items",
			input:    StringArray{"owner", "admin", "moderator"},
			expected: `{"owner","admin","moderator"}`,
		},
		{
			name:     "items with special characters",
			input:    StringArray{"server admin", "game moderator"},
			expected: `{"server admin","game moderator"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.input.Value()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, value)
		})
	}
}