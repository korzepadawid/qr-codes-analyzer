package qr_code

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const paramNameUUID = "uuid"

func (h *qrCodeHandler) qrCodeRedirect(ctx *gin.Context) {
	uuid := strings.TrimSpace(ctx.Param(paramNameUUID))

	if len(uuid) == 0 {
		ctx.HTML(http.StatusNotFound, notFoundPageTemplateName, gin.H{})
		return
	}

	v, err := h.cache.Get(uuid)

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

	// todo : save redirect to db in a separated goroutine and increase
	ctx.Redirect(http.StatusPermanentRedirect, qrCode.RedirectionUrl)
}
