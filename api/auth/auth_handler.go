package auth

import (
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
	store          db.Store
	tokenMaker     token.Maker
	passwordHasher util.Hasher
}

func NewAuthHandler(store db.Store, tokenService token.Maker, hasher util.Hasher) *authHandler {
	return &authHandler{
		store:          store,
		tokenMaker:     tokenService,
		passwordHasher: hasher,
	}
}

func (h *authHandler) RegisterRoutes(router *gin.Engine) {
	r := router.Group(routerGroupPrefix)
	r.POST(signUpUrl, h.signUp)
	r.POST(signInUrl, h.signIn)
}
