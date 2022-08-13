package token

import "errors"

var (
	ErrExpiredToken = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

type Maker interface {
	CreateToken(username string) (string, error)

	VerifyToken(token string) (*Payload, error)
}
