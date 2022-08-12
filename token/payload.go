package token

import (
	"time"

	"github.com/google/uuid"
)

// Payload implements Claims interface
// from jwt package (https://pkg.go.dev/github.com/golang-jwt/jwt#Claims)
type Payload struct {
	ID         string
	Username   string
	ExpiringAt time.Time
	IssuedAt   time.Time
}

func NewPayload(username string, duration time.Duration) *Payload {
	return &Payload{
		ID:         uuid.NewString(),
		Username:   username,
		ExpiringAt: time.Now().Add(duration),
		IssuedAt:   time.Now(),
	}
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiringAt) {
		return ErrExpiredToken
	}
	return nil
}
