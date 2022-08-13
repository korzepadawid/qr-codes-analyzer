package errors

import "errors"

var (
	ErrUserAlreadyExists          = errors.New("user's already exists")
	ErrInvalidCredentials         = errors.New("invalid credentials")
	ErrMissingAuthorizationHeader = errors.New("missing \"Authorization\" header")
	ErrInvalidAuthorizationType   = errors.New("you need to use Bearer token strategy")
)
