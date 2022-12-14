package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/common"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"github.com/korzepadawid/qr-codes-analyzer/api/group"
	qrcode "github.com/korzepadawid/qr-codes-analyzer/api/qr_code"
	"github.com/korzepadawid/qr-codes-analyzer/cache"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/encode"
	"github.com/korzepadawid/qr-codes-analyzer/ipapi"
	"github.com/korzepadawid/qr-codes-analyzer/storage"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"github.com/korzepadawid/qr-codes-analyzer/valid"
	"go.uber.org/zap"
	"log"
)

type Server struct {
	Config          config.Config
	Store           db.Store
	Router          *gin.Engine
	TokenProvider   token.Provider
	PasswordService util.PasswordService
	Handlers        []common.Handler
	storage         storage.FileStorage
	qrCodeEncoder   encode.Encoder
	cache           cache.Cache
	clientIP        ipapi.Client
}

func NewServer(
	config config.Config,
	store db.Store,
	tokenProvider token.Provider,
	passwordService util.PasswordService,
	storage storage.FileStorage,
	encoder encode.Encoder,
	cache cache.Cache,
	clientIP ipapi.Client,
) (*Server, error) {
	server := Server{
		Config:          config,
		Store:           store,
		Router:          gin.Default(),
		TokenProvider:   tokenProvider,
		PasswordService: passwordService,
		storage:         storage,
		cache:           cache,
		qrCodeEncoder:   encoder,
		clientIP:        clientIP,
		Handlers:        make([]common.Handler, 0),
	}

	// setup gin
	gin.SetMode(gin.DebugMode)

	// html templates
	server.Router.LoadHTMLGlob("./templates/*.html")

	// additional validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("url", valid.URL)
		if err != nil {
			log.Fatal(err)
		}
	}

	// setup logger
	logger, _ := zap.NewProduction()

	// setup middlewares
	server.Router.Use(errors.HandleErrors(logger))

	// route Handlers
	authHandler := auth.NewAuthHandler(server.Store, server.TokenProvider, server.PasswordService)
	groupHandler := group.NewGroupHandler(server.Store, auth.SecureRoute(server.TokenProvider))
	qrCodeHandler := qrcode.NewQRCodeHandler(server.Store, server.Config, server.storage, server.qrCodeEncoder, server.cache, server.clientIP, auth.SecureRoute(server.TokenProvider))

	server.Handlers = append(server.Handlers, authHandler, groupHandler, qrCodeHandler)

	for _, h := range server.Handlers {
		h.RegisterRoutes(server.Router)
	}

	return &server, nil
}

func (s *Server) Run() error {
	return s.Router.Run(s.Config.Addr)
}
