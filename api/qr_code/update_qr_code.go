package qr_code

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"net/http"
	"strings"
)

type UpdateQRCodeRequest struct {
	Title       string `json:"title,omitempty" binding:"max=255"`
	Description string `json:"description,omitempty" binding:"max=255"`
}

func (h *qrCodeHandler) updateQRCode(ctx *gin.Context) {
	var request UpdateQRCodeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		return
	}

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

	arg := db.UpdateQRCodeTxParams{
		UUID:        uuid,
		Owner:       owner,
		Title:       strings.TrimSpace(request.Title),
		Description: strings.TrimSpace(request.Description),
	}
	err = h.store.UpdateQRCodeTx(ctx, arg)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
