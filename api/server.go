package api

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/common"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"github.com/korzepadawid/qr-codes-analyzer/api/group"
	qrcode "github.com/korzepadawid/qr-codes-analyzer/api/qr_code"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/qr"
	"github.com/korzepadawid/qr-codes-analyzer/storage"
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
	storage        storage.FileStorage
	encoder        qr.Encoder
}

func NewServer(
	config config.Config,
	store db.Store,
	maker token.Maker,
	hasher util.Hasher,
	storage storage.FileStorage,
	encoder qr.Encoder,
) (*Server, error) {
	server := Server{
		Config:         config,
		Store:          store,
		Router:         gin.Default(),
		TokenMaker:     maker,
		PasswordHasher: hasher,
		storage:        storage,
		encoder:        encoder,
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
	groupHandler := group.NewGroupHandler(server.Store, auth.SecureRoute(server.TokenMaker))
	qrCodeHandler := qrcode.NewQRCodeHandler(server.Store, server.Config, server.storage, server.encoder, auth.SecureRoute(server.TokenMaker))

	server.Handlers = append(server.Handlers, authHandler, groupHandler, qrCodeHandler)

	for _, h := range server.Handlers {
		h.RegisterRoutes(server.Router)
	}

	return &server, nil
}

func (s *Server) Run() error {
	return s.Router.Run(s.Config.Addr)
}
