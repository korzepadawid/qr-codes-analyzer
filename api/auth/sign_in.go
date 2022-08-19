package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
	"regexp"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

type getUserFunc func(context.Context, string) (db.User, error)
type getUserExec func(*gin.Context, string) (db.User, error)

type signInRequest struct {
	Username string `json:"username" binding:"required,min=3,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type signInResponse struct {
	BearerToken string `json:"bearer_token,omitempty"`
}

func (h *authHandler) signIn(ctx *gin.Context) {
	var request signInRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		return
	}

	getUser := prepareGetUserQuery(h.store.GetUserByUsername)
	emailProvided := emailProvidedInsteadOfUsername(request.Username)

	if emailProvided {
		getUser = prepareGetUserQuery(h.store.GetUserByEmail)
	}

	user, err := getUser(ctx, request.Username)

	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.pass.VerifyPassword(user.Password, request.Password)

	if err != nil {
		ctx.Error(errors.ErrInvalidCredentials)
		return
	}

	token, err := h.token.CreateToken(user.Username)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, newSingInResponse(token))
}

func newSingInResponse(token string) signInResponse {
	return signInResponse{
		token,
	}
}

func prepareGetUserQuery(fn getUserFunc) getUserExec {
	return func(ctx *gin.Context, arg string) (db.User, error) {
		user, err := fn(ctx, arg)

		if err != nil {
			return db.User{}, errors.ErrInvalidCredentials
		}

		return user, nil
	}
}

func emailProvidedInsteadOfUsername(s string) bool {
	return emailRegex.MatchString(s)
}
