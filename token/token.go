package token

import "errors"

var (
	ErrExpiredToken = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

type Provider interface {
	CreateToken(username string) (string, error)

	VerifyToken(token string) (*Payload, error)
}
