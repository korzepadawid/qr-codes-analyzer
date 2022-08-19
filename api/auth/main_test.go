package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"go.uber.org/zap"
	"net/http"
	"os"
	"testing"
)

const (
	securedTestRouteUrl = "/secured"
)

type testResponseForSecuredRoute struct {
	User string `json:"user,omitempty"`
}

// securedRoute is an endpoint for testing auth process
func securedRoute(ctx *gin.Context) {
	value, exists := ctx.Get(currentUserKey)

	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not found in the context"})
		return
	}

	user, ok := value.(string)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not found in the context"})
		return
	}

	ctx.JSON(http.StatusOK, testResponseForSecuredRoute{user})
}

func setUpAuthMiddleware(maker token.Provider) *gin.Engine {
	r := gin.Default()
	setUpErrorHandler(r)
	r.Use(SecureRoute(maker))
	r.GET(securedTestRouteUrl, securedRoute)
	return r
}

func setUpHandler(store db.Store, maker token.Provider, hasher util.PasswordService) *gin.Engine {
	r := gin.Default()
	handler := NewAuthHandler(store, maker, hasher)
	setUpErrorHandler(r)
	handler.RegisterRoutes(r)
	return r
}

func setUpErrorHandler(r *gin.Engine) {
	logger, _ := zap.NewProduction()
	r.Use(errors.HandleErrors(logger))
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
