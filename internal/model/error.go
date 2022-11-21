package model

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrBadRequest     = errors.New("bad request")
	ErrInternal       = errors.New("internal error")
	ErrUnknown        = errors.New("unknown error")
	ErrAuthorized     = errors.New("authorized error")
)
