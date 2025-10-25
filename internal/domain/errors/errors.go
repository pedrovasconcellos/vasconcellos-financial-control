package errors

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("resource conflict")
	ErrInternal     = errors.New("internal error")
)
