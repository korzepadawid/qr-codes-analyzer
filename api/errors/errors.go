package errors

import "errors"

var (
	ErrInternalError     = errors.New("internal server error")
	ErrUserAlreadyExists = errors.New("user's already exists")
)
