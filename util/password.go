package util

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	HashPassword(rawPassword string) (string, error)

	VerifyPassword(hashedPassword, rawPassword string) error
}

type BCryptHasher struct{}

func NewBCryptHasher() Hasher {
	return &BCryptHasher{}
}

func (b *BCryptHasher) HashPassword(rawPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (b *BCryptHasher) VerifyPassword(hashedPassword, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
}
