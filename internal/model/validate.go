package model

import (
	"fmt"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// ValidRecordRequest
func ValidRecordRequest(req RecordRequest) (err error) {
	// req : YYYYMMDD or ""
	if req.Datetime != "" || validation.Validate(req.Datetime, validation.Length(6, 6), is.Digit) != nil {
		return fmt.Errorf("%s: %w", err.Error(), ErrInvalidValue)
	}

	return nil
}

func ValidYYYY(yyyy string) (yyyyint int, err error) {
	if err := validation.Validate(yyyy, validation.Length(4, 4), is.Digit); err != nil {
		return 0, fmt.Errorf("invalid YYYY: %w", ErrTableNotFound)
	}

	yyyyint, err = strconv.Atoi(yyyy)
	if err != nil {
		return 0, err
	}

	return yyyyint, err
}
