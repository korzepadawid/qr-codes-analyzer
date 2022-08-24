package token

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTTokenProvider struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	Duration   time.Duration
}

func NewJWTMaker(duration time.Duration) *JWTTokenProvider {
	privateKeyBytes, err := ioutil.ReadFile("./token/keys/rsa.private")
	if err != nil {
		log.Fatalln(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		log.Fatalln(err)
	}

	publicKeyBytes, err := ioutil.ReadFile("./token/keys/rsa.public")
	if err != nil {
		log.Fatalln(err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Fatalln(err)
	}

	return &JWTTokenProvider{
		privateKey: privateKey,
		publicKey:  publicKey,
		Duration:   duration,
	}
}

func (jwtMaker *JWTTokenProvider) CreateToken(username string) (string, error) {
	payload := NewPayload(username, jwtMaker.Duration)
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, payload).SignedString(jwtMaker.privateKey)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (jwtMaker *JWTTokenProvider) VerifyToken(t string) (*Payload, error) {

	jwtToken, err := jwt.ParseWithClaims(t, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			return jwtMaker.publicKey, nil
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
