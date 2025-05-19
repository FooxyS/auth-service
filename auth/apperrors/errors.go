package apperrors

import "errors"

var (
	ErrEmptyString  = errors.New("string is empty")
	ErrNotMatch     = errors.New("this variables don't equael")
	ErrInvalidToken = errors.New("invalid Token")
)
