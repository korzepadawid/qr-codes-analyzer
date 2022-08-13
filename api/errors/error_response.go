package errors

type validationFieldError struct {
	Field  string `json:"field,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type ErrorResponse struct {
	Message string                 `json:"message,omitempty"`
	Fields  []validationFieldError `json:"fields,omitempty"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Message: err.Error(),
	}
}
