package model

import "fmt"

var ErrInvalidValue error = fmt.Errorf("invalid value error")
var ErrUnknownCategoryID error = fmt.Errorf("unknown category ID")
var ErrNotFound error = fmt.Errorf("record not found")
var ErrAlreadyRecorded error = fmt.Errorf("record already recorded")
