package errors

import "errors"

var (
	ErrUserAlreadyExists          = errors.New("user's already exists")
	ErrInvalidCredentials         = errors.New("invalid credentials")
	ErrMissingAuthorizationHeader = errors.New("missing \"Authorization\" header")
	ErrInvalidAuthorizationType   = errors.New("you need to use Bearer token strategy")
	ErrFailedCurrentUserRetrieval = errors.New("failed current user retrieval")
	ErrInvalidParamFormat         = errors.New("invalid param format")
	ErrGroupNotFound              = errors.New("group not found")
	ErrQRCodeGenerationFailed     = errors.New("failed to generate new qr code")
)
