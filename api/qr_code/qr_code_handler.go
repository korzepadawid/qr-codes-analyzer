package qr_code

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/encode"
	"github.com/korzepadawid/qr-codes-analyzer/storage"
)

const (
	routePrefix          = "/encode-codes"
	routePrefixWithGroup = "/groups/:group_id/encode-codes"
)

type qrCodeHandler struct {
	config        config.Config
	store         db.Store
	middlewares   []gin.HandlerFunc
	storage       storage.FileStorage
	qrCodeEncoder encode.Encoder
}

func NewQRCodeHandler(store db.Store, config config.Config, fileStorage storage.FileStorage, qrCodeEncoder encode.Encoder, middlewares ...gin.HandlerFunc) *qrCodeHandler {
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