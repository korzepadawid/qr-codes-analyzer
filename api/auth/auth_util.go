package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
)

func GetCurrentUserUsername(ctx *gin.Context) (string, bool) {
	v, exists := ctx.Get(currentUserKey)

	if !exists {
		ctx.Error(errors.ErrFailedCurrentUserRetrieval)
		return "", false
	}

	owner, ok := v.(string)

	if !ok {
		ctx.Error(errors.ErrFailedCurrentUserRetrieval)
		return "", false
	}

	return owner, true
}
