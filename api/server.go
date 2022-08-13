package api

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/common"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"go.uber.org/zap"
)

type Server struct {
	config   config.Config
	store    db.Store
	router   *gin.Engine
	handlers []common.Handler
}

func NewServer(config config.Config, store db.Store) (*Server, error) {
	server := Server{
		config:   config,
		store:    store,
		router:   gin.Default(),
		handlers: make([]common.Handler, 0),
	}

	// setup uber's logger
	logger, _ := zap.NewProduction()

	// setup middlewares
	server.router.Use(errors.HandleErrors(logger))

	// route handlers
	server.handlers = append(server.handlers, auth.NewAuthHandler(store))

	for _, h := range server.handlers {
		h.RegisterRoutes(server.router)
	}

	return &server, nil
}

func (s *Server) Run() error {
	return s.router.Run(s.config.Addr)
}
