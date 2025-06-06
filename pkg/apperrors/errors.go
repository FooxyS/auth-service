package apperrors

import "errors"

var (
	ErrUserExists       = errors.New("user with this email already exists")
	ErrPasswordMismatch = errors.New("password mismatch")
	ErrIPMismatch       = errors.New("IP mismatch")
	ErrAgentMismatch    = errors.New("agent mismatch, session was deleted")
	ErrNotValid         = errors.New("this argument is not valid")
	ErrSessionNotFound  = errors.New("session not found")

	ErrFindSession = errors.New("session is not exists")

	ErrBearer = errors.New("bearer header is incorrect")
)
