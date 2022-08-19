package qr_code

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

const paramNameUUID = "uuid"

func (h *qrCodeHandler) qrCodeRedirect(ctx *gin.Context) {
	uuid := ctx.Param(paramNameUUID)
	qrCode, err := h.store.GetQRCode(ctx, uuid)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.HTML(http.StatusNotFound, notFoundPageTemplateName, gin.H{})
			return
		}
		ctx.Error(err)
		return
	}

	// todo: save redirect to db in a separated goroutine
	ctx.Redirect(http.StatusPermanentRedirect, qrCode.RedirectionUrl)
}
