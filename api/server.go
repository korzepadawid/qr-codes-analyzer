package api

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/common"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"go.uber.org/zap"
	"time"
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
	// deps
	tokenService := token.NewJWTMaker("asdfsafd", time.Hour)

	// setup gin
	gin.SetMode(gin.DebugMode)

	// setup uber's zap logger
	logger, _ := zap.NewProduction()

	// setup middlewares
	server.router.Use(errors.HandleErrors(logger))

	// route handlers
	server.handlers = append(server.handlers, auth.NewAuthHandler(store, tokenService))

	for _, h := range server.handlers {
		h.RegisterRoutes(server.router)
	}

	server.router.Use(auth.SecureRoute(tokenService))

	server.router.GET("/test", func(context *gin.Context) {
		value, exists := context.Get(auth.CurrentUserKey)

		if !exists {
			context.JSON(418, gin.H{"error": "not exists"})
			return
		}

		s, ok := value.(string)
		if !ok {
			context.JSON(418, gin.H{"error": "not exists"})
			return
		}

		context.JSON(418, gin.H{"user": s})
	})

	return &server, nil
}

func (s *Server) Run() error {
	return s.router.Run(s.config.Addr)
}
