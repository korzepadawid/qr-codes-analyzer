package valid

import (
	"github.com/go-playground/validator/v10"
	"net/url"
)

// IsUrlValid validates given url,
//returns true when url is valid
func isUrlValid(s string) bool {
	_, err := url.ParseRequestURI(s)

	if err != nil {
		return false
	}

	return true
}

var URL validator.Func = func(fieldLevel validator.FieldLevel) bool {
	v, ok := fieldLevel.Field().Interface().(string)

	if !ok {
		return false
	}

	return isUrlValid(v)
}
