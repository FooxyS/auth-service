package apperrors

import "errors"

var (
	ErrUserExists       = errors.New("user with this email already exists")
	ErrPasswordMismatch = errors.New("password mismatch")
)
