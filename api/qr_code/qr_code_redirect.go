package qr_code

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

const paramNameUUID = "uuid"

func (h *qrCodeHandler) qrCodeRedirect(ctx *gin.Context) {
	uuid := ctx.Param(paramNameUUID)

	v, err := h.cache.Get(paramNameUUID)

	if err == nil {
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

	// todo: save redirect to db in a separated goroutine
	ctx.Redirect(http.StatusPermanentRedirect, qrCode.RedirectionUrl)
}
