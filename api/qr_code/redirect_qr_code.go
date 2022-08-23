package qr_code

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/ipapi"
	"log"
	"net/http"
)

const paramNameUUID = "uuid"

func (h *qrCodeHandler) qrCodeRedirect(ctx *gin.Context) {
	uuid, err := getQRCodeUUID(ctx)
	if err != nil {
		ctx.HTML(http.StatusNotFound, notFoundPageTemplateName, gin.H{})
		return
	}

	v, err := h.cache.Get(uuid)

	if err == nil {
		go h.createRedirectEntry(ctx, uuid)
		ctx.Header("Cache-Control", "no-store")
		ctx.Redirect(http.StatusPermanentRedirect, v)
		return
	}

	qrCode, err := h.store.GetQRCode(ctx, uuid)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.HTML(http.StatusNotFound, notFoundPageTemplateName, gin.H{})
			return
		}
		ctx.Error(err)
		return
	}

	go h.cacheQRCode(qrCode.Uuid, qrCode.RedirectionUrl)
	go h.createRedirectEntry(ctx, qrCode.Uuid)

	ctx.Header("Cache-Control", "no-store")
	ctx.Redirect(http.StatusPermanentRedirect, qrCode.RedirectionUrl)
}

func (h *qrCodeHandler) createRedirectEntry(ctx *gin.Context, uuid string) {
	c := ipapi.New()
	det, err := c.GetIPDetails("142.250.203.206")

	if err != nil {
		log.Println(err)
		return
	}

	if err := h.store.IncrementRedirectEntriesTx(ctx, db.IncrementRedirectEntriesTxParams{
		UUID:      uuid,
		IPv4:      "142.250.203.206",
		IPDetails: det,
	}); err != nil {
		log.Printf("%v", err)
	}
}
