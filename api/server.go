package api

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/common"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"go.uber.org/zap"
)

type Server struct {
	Config         config.Config
	Store          db.Store
	Router         *gin.Engine
	TokenMaker     token.Maker
	PasswordHasher util.Hasher
	Handlers       []common.Handler
}

func NewServer(config config.Config, store db.Store, maker token.Maker, hasher util.Hasher) (*Server, error) {
	server := Server{
		Config:         config,
		Store:          store,
		Router:         gin.Default(),
		TokenMaker:     maker,
		PasswordHasher: hasher,
		Handlers:       make([]common.Handler, 0),
	}

	// setup gin
	gin.SetMode(gin.DebugMode)

	// setup logger
	logger, _ := zap.NewProduction()

	// setup middlewares
	server.Router.Use(errors.HandleErrors(logger))

	// route Handlers
	authHandler := auth.NewAuthHandler(server.Store, server.TokenMaker, server.PasswordHasher)
	server.Handlers = append(server.Handlers, authHandler)

	for _, h := range server.Handlers {
		h.RegisterRoutes(server.Router)
	}

	server.Router.Use(auth.SecureRoute(server.TokenMaker))

	return &server, nil
}

func (s *Server) Run() error {
	return s.Router.Run(s.Config.Addr)
}
