package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"github.com/korzepadawid/qr-codes-analyzer/util"
)

const (
	routerGroupPrefix = "/auth"
	signUpUrl         = "/signup"
	signInUrl         = "/signin"
)

type authHandler struct {
	tokenService    token.Maker
	passwordService util.Hasher
	store           db.Store
}

func NewAuthHandler(store db.Store) *authHandler {
	return &authHandler{
		store:           store,
		tokenService:    token.NewJWTMaker("asdfsafd", time.Hour),
		passwordService: util.NewBCryptHasher(),
	}
}

func (h *authHandler) RegisterRoutes(router *gin.Engine) {
	r := router.Group(routerGroupPrefix)
	r.POST(signUpUrl, h.signUp)
	r.POST(signInUrl, func(ctx *gin.Context) {})
}
