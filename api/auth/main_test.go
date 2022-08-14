package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"go.uber.org/zap"
	"os"
	"testing"
)

func setUpHandler(store db.Store, maker token.Maker, hasher util.Hasher) *gin.Engine {
	r := gin.Default()
	handler := NewAuthHandler(store, maker, hasher)
	logger, _ := zap.NewProduction()
	r.Use(errors.HandleErrors(logger))
	handler.RegisterRoutes(r)
	return r
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
