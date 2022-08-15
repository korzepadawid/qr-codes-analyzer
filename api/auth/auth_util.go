package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
)

func GetCurrentUserUsername(ctx *gin.Context) (string, error) {
	v, exists := ctx.Get(currentUserKey)

	if !exists {
		return "", errors.ErrFailedCurrentUserRetrieval
	}

	owner, ok := v.(string)

	if !ok {
		return "", errors.ErrFailedCurrentUserRetrieval
	}

	return owner, nil
}
