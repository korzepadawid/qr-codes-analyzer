package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrUserAlreadyExists = errors.New("user's already exists")
)

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
