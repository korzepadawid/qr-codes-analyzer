package qr_code

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/storage"
	"log"
	"net/http"
)

func (h *qrCodeHandler) deleteQRCode(ctx *gin.Context) {
	uuid, err := getQRCodeUUID(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	owner, err := auth.GetCurrentUserUsername(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	arg := db.DeleteQRCodeParams{
		Uuid:  uuid,
		Owner: owner,
	}
	if err = h.store.DeleteQRCode(ctx, arg); err != nil {
		ctx.Error(err)
		return
	}

	go func() {
		if err := h.storage.DeleteFile(ctx, uuid+storage.ImageExt); err != nil {
			log.Print(err)
		}
	}()

	ctx.Status(http.StatusNoContent)
}
