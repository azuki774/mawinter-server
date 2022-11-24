package model

import "fmt"

var ErrInvalidValue error = fmt.Errorf("invalid value error")
var ErrNotFound error = fmt.Errorf("record not found")
var ErrTableNotFound error = fmt.Errorf("table not found")
