package qr_code

import (
	"database/sql"
	"github.com/gin-gonic/gin"
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
		h.redirectionWorker <- saveRedirectJob{
			UUID: uuid,
			IPv4: ctx.ClientIP(),
		}
		ctx.Header("Cache-Control", "no-store")
		ctx.Redirect(http.StatusPermanentRedirect, v)
		return
	}

	qr, err := h.store.GetQRCode(ctx, uuid)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.HTML(http.StatusNotFound, notFoundPageTemplateName, gin.H{})
			return
		}
		ctx.Error(err)
		return
	}

	h.cacheWorker <- cacheQRCodeJob{
		Key:   qr.Uuid,
		Value: qr.RedirectionUrl,
	}
	h.redirectionWorker <- saveRedirectJob{
		UUID: uuid,
		IPv4: ctx.ClientIP(),
	}
	ctx.Header("Cache-Control", "no-store")
	ctx.Redirect(http.StatusPermanentRedirect, qr.RedirectionUrl)
}
