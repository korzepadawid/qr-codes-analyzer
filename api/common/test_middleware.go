package common

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"go.uber.org/zap"
)

func SetUpErrorHandler(r *gin.Engine) {
	logger, _ := zap.NewProduction()
	r.Use(errors.HandleErrors(logger))
}
