package errors

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"go.uber.org/zap"
	"net/http"
)

func HandleErrors(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		err := ctx.Errors.Last()

		if err == nil {
			return
		}

		logger.Error(err.Error())

		var verr validator.ValidationErrors

		if errors.As(err, &verr) {
			ctx.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Errorf("validation error")))
		} else if errors.Is(err, ErrInvalidParamFormat) {
			ctx.JSON(http.StatusBadRequest, NewErrorResponse(ErrInvalidParamFormat))
		} else if errors.Is(err, token.ErrExpiredToken) || errors.Is(err, token.ErrInvalidToken) {
			ctx.JSON(http.StatusUnauthorized, NewErrorResponse(err))
		} else if errors.Is(err, ErrGroupNotFound) {
			ctx.JSON(http.StatusNotFound, NewErrorResponse(err))
		} else if errors.Is(err, ErrUserAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, NewErrorResponse(err))
		} else if errors.Is(err, ErrMissingAuthorizationHeader) || errors.Is(err, ErrInvalidAuthorizationType) {
			ctx.JSON(http.StatusUnauthorized, NewErrorResponse(err))
		} else if errors.Is(err, ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, NewErrorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Errorf("internal error")))
		}
	}
}
