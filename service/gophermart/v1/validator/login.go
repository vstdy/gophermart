package validator

import (
	"fmt"
)

// ValidateLogin validates login.
func ValidateLogin(login string) error {
	if login == "" {
		return fmt.Errorf("empty")
	}

	return nil
}
