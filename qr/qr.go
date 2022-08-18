package qr

import "github.com/skip2/go-qrcode"

const qrCodeSize = 256

type Encoder interface {
	Encode(string) ([]byte, error)
}

type qrCodeEncoder struct{}

func NewQRCodeHandler() Encoder {
	return qrCodeEncoder{}
}

func (q qrCodeEncoder) Encode(url string) ([]byte, error) {
	return qrcode.Encode(url, qrcode.Medium, qrCodeSize)
}
