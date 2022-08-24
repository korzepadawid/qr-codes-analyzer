package qr_code

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/cache"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/encode"
	"github.com/korzepadawid/qr-codes-analyzer/ipapi"
	"github.com/korzepadawid/qr-codes-analyzer/storage"
)

const (
	notFoundPageTemplateName = "not_found.html"
	routePrefix              = "/qr-codes/:uuid"
	routeRedirect            = "/qr-codes/:uuid/redirect"
	routePrefixWithGroup     = "/groups/:group_id/qr-codes"
)

type qrCodeHandler struct {
	config            config.Config
	store             db.Store
	middlewares       []gin.HandlerFunc
	storage           storage.FileStorage
	qrCodeEncoder     encode.Encoder
	cache             cache.Cache
	clientIP          ipapi.Client
	cacheWorker       chan cacheQRCodeJob
	redirectionWorker chan saveRedirectJob
}

func NewQRCodeHandler(
	store db.Store,
	config config.Config,
	fileStorage storage.FileStorage,
	qrCodeEncoder encode.Encoder,
	cache cache.Cache,
	clientIP ipapi.Client,
	middlewares ...gin.HandlerFunc,
) *qrCodeHandler {
	return &qrCodeHandler{
		config:            config,
		store:             store,
		middlewares:       middlewares,
		storage:           fileStorage,
		qrCodeEncoder:     qrCodeEncoder,
		cache:             cache,
		clientIP:          clientIP,
		redirectionWorker: make(chan saveRedirectJob),
		cacheWorker:       make(chan cacheQRCodeJob),
	}
}

func (h qrCodeHandler) RegisterRoutes(r *gin.Engine) {
	go h.saveRedirectWorker()
	go h.cacheQRCodesWorker()
	r.GET(routeRedirect, h.qrCodeRedirect) // the redirect route is publicly accessible
	r.Use(h.middlewares...)
	r.POST(routePrefixWithGroup, h.createQRCode)
	r.DELETE(routePrefix, h.deleteQRCode)
	r.GET(routePrefixWithGroup, h.getQRCodes)
	r.GET(routePrefix+"/stats/csv", h.createQRCodeStatsCSV)
	r.PATCH(routePrefix, h.updateQRCode)
}
