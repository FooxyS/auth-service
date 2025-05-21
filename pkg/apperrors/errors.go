package apperrors

import "errors"

var (
	ErrUserExists       = errors.New("user with this email already exists")
	ErrPasswordMismatch = errors.New("password mismatch")
	ErrIPMismatch       = errors.New("IP mismatch")
	ErrAgentMismatch    = errors.New("agent mismatch, session was deleted")
)
