package qr_code

import (
	"github.com/korzepadawid/qr-codes-analyzer/cache"
	"time"
)

func (h *qrCodeHandler) cacheQRCode(uuid, url string) {
	params := cache.SetParams{
		Key:      uuid,
		Value:    url,
		Duration: time.Minute * 2,
	}

	if err := h.cache.Set(&params); err != nil {
		panic(err)
	}
}
