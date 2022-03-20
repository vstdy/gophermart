package validator

import (
	"github.com/ShiraazMoollatjie/goluhn"

	"github.com/vstdy0/go-diploma/pkg"
)

// ValidateOrderNumber validates order number.
func ValidateOrderNumber(orderID string) error {
	if len(orderID) < 5 || len(orderID) > 16 {
		return pkg.ErrInvalidInput
	}
	if err := goluhn.Validate(orderID); err != nil {
		return pkg.ErrInvalidInput
	}

	return nil
}
