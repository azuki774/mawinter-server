package api

import (
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func isValidYearMonth(yyyymm string) (err error) {
	err = validation.Validate(yyyymm,
		is.Digit,
		validation.Length(6, 6),
	)
	return err
}
