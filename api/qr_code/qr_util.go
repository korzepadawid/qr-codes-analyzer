package qr_code

import (
	"github.com/korzepadawid/qr-codes-analyzer/cache"
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
