package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"strings"
)

const (
	currentUserKey         = "auth_user"
	authorizationHeaderKey = "Authorization"
	authorizationType      = "Bearer"
)

func SecureRoute(tokenService token.Provider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		is := isEmptyHeader(authorizationHeader)

		if is {
			abortWithError(ctx, errors.ErrMissingAuthorizationHeader)
			return
		}

		split := strings.Split(strings.TrimSpace(authorizationHeader), " ")

		if len(split) != 2 {
			abortWithError(ctx, errors.ErrInvalidAuthorizationType)
			return
		}

		givenAuthorizationType := split[0]
		givenToken := split[1]

		notValidAuthorizationType := isNotValidAuthorizationType(givenAuthorizationType, givenToken)

		if notValidAuthorizationType {
			abortWithError(ctx, errors.ErrInvalidAuthorizationType)
			return
		}

		payload, err := tokenService.VerifyToken(givenToken)

		if err != nil {
			abortWithError(ctx, err)
			return
		}

		ctx.Set(currentUserKey, payload.Username)
		ctx.Next()
	}
}

func abortWithError(ctx *gin.Context, err error) {
	ctx.Error(err)
	ctx.Abort()
}

func isEmptyHeader(authorizationHeader string) bool {
	return len(authorizationHeader) == 0
}

func isNotValidAuthorizationType(givenAuthorizationType string, givenToken string) bool {
	return givenAuthorizationType != authorizationType || len(givenToken) < 1
}
