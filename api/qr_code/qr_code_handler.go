package qr_code

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/cache"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/encode"
	"github.com/korzepadawid/qr-codes-analyzer/storage"
)

const (
	notFoundPageTemplateName = "not_found.html"
	routePrefix              = "/qr-codes/:uuid"
	routeRedirect            = "/qr-codes/:uuid/redirect"
	routePrefixWithGroup     = "/groups/:group_id/qr-codes"
)

type qrCodeHandler struct {
	config        config.Config
	store         db.Store
	middlewares   []gin.HandlerFunc
	storage       storage.FileStorage
	qrCodeEncoder encode.Encoder
	cache         cache.Cache
}

func NewQRCodeHandler(
	store db.Store,
	config config.Config,
	fileStorage storage.FileStorage,
	qrCodeEncoder encode.Encoder,
	cache cache.Cache,
	middlewares ...gin.HandlerFunc,
) *qrCodeHandler {
	return &qrCodeHandler{
		config:        config,
		store:         store,
		middlewares:   middlewares,
		storage:       fileStorage,
		qrCodeEncoder: qrCodeEncoder,
		cache:         cache,
	}
}

func (h qrCodeHandler) RegisterRoutes(r *gin.Engine) {
	r.GET(routeRedirect, h.qrCodeRedirect) // the redirect route is publicly accessible
	r.Use(h.middlewares...)
	r.POST(routePrefixWithGroup, h.createQRCode)
	r.DELETE(routePrefix, h.deleteQRCode)
}
