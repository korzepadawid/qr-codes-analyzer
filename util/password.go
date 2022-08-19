package util

import "golang.org/x/crypto/bcrypt"

type PasswordService interface {
	HashPassword(rawPassword string) (string, error)

	VerifyPassword(hashedPassword, rawPassword string) error
}

type BCryptPasswordService struct{}

func NewBCryptHasher() *BCryptPasswordService {
	return &BCryptPasswordService{}
}

func (b *BCryptPasswordService) HashPassword(rawPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (b *BCryptPasswordService) VerifyPassword(hashedPassword, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
}
