package validator

import (
	"fmt"
)

// ValidatePassword validates password.
func ValidatePassword(login string) error {
	if login == "" {
		return fmt.Errorf("empty")
	}

	return nil
}
