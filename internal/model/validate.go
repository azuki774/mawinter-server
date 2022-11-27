package model

import (
	"fmt"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)
func ValidYYYY(yyyy string) (yyyyint int, err error) {
	if err := validation.Validate(yyyy, validation.Length(4, 4), is.Digit); err != nil {
		return 0, fmt.Errorf("invalid YYYY: %w", ErrInvalidValue)
	}

	yyyyint, err = strconv.Atoi(yyyy)
	if err != nil {
		return 0, err
	}

	return yyyyint, err
}
