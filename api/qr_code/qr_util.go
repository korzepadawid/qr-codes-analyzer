package qr_code

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"github.com/korzepadawid/qr-codes-analyzer/cache"
	"strings"
	"time"
)

func (h *qrCodeHandler) cacheQRCode(key, value string) {
	params := cache.SetParams{
		Key:      key,
		Value:    value,
		Duration: time.Minute * 2,
	}

	if err := h.cache.Set(&params); err != nil {
		panic(err)
	}
}

func getQRCodeUUID(ctx *gin.Context) (string, error) {
	param := ctx.Param(paramNameUUID)
	param = strings.TrimSpace(param)

	if len(param) == 0 {
		return "", errors.ErrQRCodeInvalidUUID
	}

	return param, nil
}
