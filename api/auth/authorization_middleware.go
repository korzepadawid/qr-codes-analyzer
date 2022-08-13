package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"strings"
)

const (
	CurrentUserKey         = "auth_user"
	authorizationHeaderKey = "Authorization"
	authorizationType      = "Bearer"
)

func SecureRoute(tokenService token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if isEmptyHeader(authorizationHeader) {
			ctx.Error(errors.ErrMissingAuthorizationHeader)
			ctx.Abort()
			return
		}

		split := strings.Split(strings.TrimSpace(authorizationHeader), " ")
		givenAuthorizationType := split[0]
		givenToken := split[1]

		if isNotValidAuthorizationType(givenAuthorizationType, givenToken) {
			ctx.Error(errors.ErrInvalidAuthorizationType)
			ctx.Abort()
			return
		}

		payload, err := tokenService.VerifyToken(givenToken)

		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set(CurrentUserKey, payload.Username)
		ctx.Next()
	}
}

func isEmptyHeader(authorizationHeader string) bool {
	return len(authorizationHeader) == 0
}

func isNotValidAuthorizationType(givenAuthorizationType string, givenToken string) bool {
	return givenAuthorizationType != authorizationType || len(givenToken) < 1
}
