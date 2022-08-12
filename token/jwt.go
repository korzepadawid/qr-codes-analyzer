package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTTokenizer struct {
	SymmetricKey string
	Duration     time.Duration
}

func NewJWTTokenizer(symmetricKey string, duration time.Duration) *JWTTokenizer {
	return &JWTTokenizer{
		SymmetricKey: symmetricKey,
		Duration:     duration,
	}
}

func (tokenizer *JWTTokenizer) CreateToken(username string) (string, error) {
	payload := NewPayload(username, tokenizer.Duration)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(tokenizer.SymmetricKey))
}

func (tokenizer *JWTTokenizer) VerifyToken(t string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(t, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(tokenizer.SymmetricKey), nil
		}
		return nil, ErrInvalidToken
	})

	if err != nil {
		jerr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(jerr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
