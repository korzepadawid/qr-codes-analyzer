package token

import "errors"

var (
	ErrExpiredToken = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

//Provider provides interface for Bearer token auth mechanism
//creating tokens and verifying them
type Provider interface {

	// CreateToken creates a new token,
	//otherwise returns an error
	CreateToken(username string) (string, error)

	// VerifyToken verifies the token and returns its payload,
	// otherwise returns an error
	VerifyToken(token string) (*Payload, error)
}
