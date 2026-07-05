package variables

import (
	"testing"
)

func TestValues_Replace(t *testing.T) {
	v := &Values{
		ProjectName: "my-api",
		ModuleName:  "github.com/user/my-api",
		Year:        "2026",
		Author:      "Alice",
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "replace all",
			input:    "Module: {{MODULE_NAME}}, Project: {{PROJECT_NAME}}, Year: {{YEAR}}, Author: {{AUTHOR}}",
			expected: "Module: github.com/user/my-api, Project: my-api, Year: 2026, Author: Alice",
		},
		{
			name:     "no replacements needed",
			input:    "Just plain text",
			expected: "Just plain text",
		},
		{
			name:     "multiple occurrences",
			input:    "{{AUTHOR}} wrote {{PROJECT_NAME}}. Yes, {{AUTHOR}} did.",
			expected: "Alice wrote my-api. Yes, Alice did.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := v.Replace(tt.input)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
