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

	ok := h.checkIfUserExists(ctx, request)

	if !ok {
		return
	}

	hashedPassword, err := h.passwordHasher.HashPassword(request.Password)

	if err != nil {
		ctx.Error(err)
		return
	}

	user, ok := h.createNewUser(ctx, request, hashedPassword)

	if !ok {
		return
	}

	ctx.JSON(http.StatusCreated, mapUserToResponse(user))
}

func (h *authHandler) createNewUser(ctx *gin.Context, request signUpRequest, hashedPassword string) (db.User, bool) {
	arg := db.CreateUserParams{
		Username: request.Username,
		Password: hashedPassword,
		Email:    request.Email,
		FullName: request.FullName,
	}

	user, err := h.store.CreateUser(ctx, arg)

	if err != nil {
		ctx.Error(err)
		return db.User{}, false
	}

	return user, true
}

func (h *authHandler) checkIfUserExists(ctx *gin.Context, request signUpRequest) bool {
	arg := db.GetUserByUsernameOrEmailParams{
		Username: strings.ToLower(request.Username),
		Email:    strings.ToLower(request.Email),
	}

	_, err := h.store.GetUserByUsernameOrEmail(ctx, arg)

	if err != nil {
		if err != sql.ErrNoRows {
			ctx.Error(err)
			return false
		}
	}

	if err == nil {
		ctx.Error(errors.ErrUserAlreadyExists)
		return false
	}
	return true
}

func mapUserToResponse(user db.User) userResponse {
	return userResponse{
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
	}
}
