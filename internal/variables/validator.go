package variables

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidModuleName = errors.New("invalid Go module name")
	moduleNameRegex      = regexp.MustCompile(`^[a-zA-Z0-9_\-\./]+$`)
)

// ValidateModuleName checks if the provided module name is roughly valid.
func ValidateModuleName(name string) error {
	if name == "" {
		return errors.New("module name cannot be empty")
	}
	if !moduleNameRegex.MatchString(name) {
		return ErrInvalidModuleName
	}
	return nil
}
