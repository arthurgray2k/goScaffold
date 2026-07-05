package variables

import (
	"testing"
)

func TestValidateModuleName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid standard", "my-api", false},
		{"valid github", "github.com/user/my-api", false},
		{"valid complex", "gitlab.com/org-name/repo_name.v2", false},
		{"invalid spaces", "my api", true},
		{"invalid symbols", "github.com/user/my@api", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModuleName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateModuleName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
