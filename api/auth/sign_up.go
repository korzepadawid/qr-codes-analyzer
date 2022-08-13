package auth

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
)

type signUpRequest struct {
	Username string `json:"username" binding:"required,min=3,max=255,alphanum"`
	Email    string `json:"email" binding:"required,email,min=5,max=255"`
	FullName string `json:"full_name" binding:"required,min=3,max=255,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type userResponse struct {
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	FullName  string    `json:"full_name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (h *authHandler) signUp(ctx *gin.Context) {
	var request signUpRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		return
	}

	gArg := db.GetUserByUsernameOrEmailParams{
		Username: strings.ToLower(request.Username),
		Email:    strings.ToLower(request.Email),
	}

	_, err := h.store.GetUserByUsernameOrEmail(ctx, gArg)

	if err != nil {
		if err != sql.ErrNoRows {
			ctx.Error(err)
			return
		}
	}

	if err == nil {
		ctx.Error(errors.ErrUserAlreadyExists)
		return
	}

	hashedPassword, err := h.passwordService.HashPassword(request.Password)

	if err != nil {
		ctx.Error(err)
		return
	}

	cArg := db.CreateUserParams{
		Username: request.Username,
		Password: hashedPassword,
		Email:    request.Email,
		FullName: request.FullName,
	}

	user, err := h.store.CreateUser(ctx, cArg)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, mapUserToResponse(user))
}

func mapUserToResponse(user db.User) userResponse {
	return userResponse{
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
	}
}
