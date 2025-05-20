package apperrors

import "errors"

var (
	ErrEmptyString      = errors.New("string is empty")
	ErrNotMatch         = errors.New("this variables don't equael")
	ErrInvalidToken     = errors.New("invalid Token")
	ErrMethodNotSupport = errors.New("request has unsupported method")
	ErrSaveUser         = errors.New("error with saving user to db")
	ErrUserExist        = errors.New("user already exists")
	ErrEmptyField       = errors.New("request data has empty fields")
)
