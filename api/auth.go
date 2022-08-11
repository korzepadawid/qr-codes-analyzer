package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) signUp(ctx *gin.Context) {
	ctx.JSON(http.StatusTeapot, gin.H{"ok": "ok"})
}

func (s *Server) signIn(ctx *gin.Context) {
	ctx.JSON(http.StatusTeapot, gin.H{"ok": "ok"})
}
