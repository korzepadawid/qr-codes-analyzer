package qr_code

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/qr"
	"github.com/korzepadawid/qr-codes-analyzer/storage"
)

const (
	routePrefix          = "/qr-codes"
	routePrefixWithGroup = "/groups/:group_id/qr-codes"
)

type qrCodeHandler struct {
	config        config.Config
	store         db.Store
	middlewares   []gin.HandlerFunc
	storage       storage.FileStorage
	qrCodeEncoder qr.Encoder
}

func NewQRCodeHandler(store db.Store, config config.Config, fileStorage storage.FileStorage, qrCodeEncoder qr.Encoder, middlewares ...gin.HandlerFunc) *qrCodeHandler {
	return &qrCodeHandler{
		config:        config,
		store:         store,
		middlewares:   middlewares,
		storage:       fileStorage,
		qrCodeEncoder: qrCodeEncoder,
	}
}

func (h qrCodeHandler) RegisterRoutes(r *gin.Engine) {
	r.Use(h.middlewares...)
	r.POST(routePrefixWithGroup, h.createQRCode)
}
