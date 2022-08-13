package errors

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		} else if errors.Is(err, ErrUserAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, NewErrorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Errorf("internal error")))
		}
	}
}
