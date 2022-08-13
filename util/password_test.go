package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func randomHash(t *testing.T, rawPassword string) string {
	bcryptHasher := NewBCryptHasher()
	hashedPassword, err := bcryptHasher.HashPassword(rawPassword)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	require.NotEqual(t, hashedPassword, rawPassword)

	return hashedPassword
}

func TestHashPassword(t *testing.T) {
	testCases := []struct {
		description string
		testCase    func(t *testing.T)
	}{
		{
			description: "hashed password should be different than the raw",
			testCase: func(t *testing.T) {
				randomHash(t, RandomString(8))
			},
		},
		{
			description: "password hashed for the second time should be different than the previous has",
			testCase: func(t *testing.T) {
				rawPassword := RandomString(8)

				hashedPassword1 := randomHash(t, rawPassword)

				hashedPassword2 := randomHash(t, rawPassword)

				require.NotEqual(t, hashedPassword1, hashedPassword2)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, tc.testCase)
	}
}

func TestVerifyPassword(t *testing.T) {
	bcryptHasher := NewBCryptHasher()

	testCases := []struct {
		description string
		testCase    func(t *testing.T)
	}{
		{
			description: "should not return any error when match",
			testCase: func(t *testing.T) {
				rawPassword := RandomString(8)
				hashedPassword := randomHash(t, rawPassword)

				err := bcryptHasher.VerifyPassword(hashedPassword, rawPassword)
				require.NoError(t, err)
			},
		},
		{
			description: "should return password missmatch error when no match",
			testCase: func(t *testing.T) {
				rawPassword := RandomString(8)
				hashedPassword := randomHash(t, rawPassword)

				err := bcryptHasher.VerifyPassword(hashedPassword, RandomString(7))
				require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, tc.testCase)
	}
}
