package api

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/util"
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

func (s *Server) signUp(ctx *gin.Context) {
	var request signUpRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	gArg := db.GetUserByUsernameOrEmailParams{
		Username: strings.ToLower(request.Username),
		Email:    strings.ToLower(request.Email),
	}

	_, err := s.store.GetUserByUsernameOrEmail(ctx, gArg)

	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errorResponse(ErrInternalError))
			return
		}
	}

	if err == nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrUserAlreadyExists))
		return
	}

	hashedPassword, err := util.HashPassword(request.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cArg := db.CreateUserParams{
		Username: request.Username,
		Password: hashedPassword,
		Email:    request.Email,
		FullName: request.FullName,
	}

	user, err := s.store.CreateUser(ctx, cArg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, mapUserToResponse(user))
}

func (s *Server) signIn(ctx *gin.Context) {
	ctx.JSON(http.StatusTeapot, gin.H{"ok": "ok"})
}
