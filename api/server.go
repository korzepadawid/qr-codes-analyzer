package api

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
)

type Server struct {
	config config.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(config config.Config, store db.Store) (*Server, error) {
	server := Server{
		config: config,
		store:  store,
		router: gin.Default(),
	}

	server.router.POST("/auth/signup", server.signUp)
	server.router.POST("/auth/signin", server.signIn)

	return &server, nil
}

func (s *Server) Run() error {
	return s.router.Run(s.config.Addr)
}
