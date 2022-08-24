package token

import (
	"testing"
	"time"

	"github.com/korzepadawid/qr-codes-analyzer/util"
	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	testCases := []struct {
		description string
		testCase    func(t *testing.T)
	}{
		{
			description: "should create a new token with valid payload",
			testCase: func(t *testing.T) {
				jwtTokenizer := NewJWTMaker(time.Minute)
				username := util.RandomString(8)
				token, err := jwtTokenizer.CreateToken(username)

				require.NoError(t, err)
				require.NotEmpty(t, token)

				payload, err := jwtTokenizer.VerifyToken(token)

				require.NoError(t, err)
				require.NotEmpty(t, payload)
				require.NotEmpty(t, payload.ID)
				require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second*2)
				require.WithinDuration(t, time.Now().Add(time.Minute), payload.ExpiringAt, time.Second*2)
				require.Equal(t, username, payload.Username)
			},
		},
		{
			description: "should return an error when token is expired",
			testCase: func(t *testing.T) {
				jwtTokenizer := NewJWTMaker(-time.Second)
				username := util.RandomString(8)
				token, err := jwtTokenizer.CreateToken(username)

				require.NoError(t, err)
				require.NotEmpty(t, token)

				payload, err := jwtTokenizer.VerifyToken(token)
				require.Empty(t, payload)
				require.EqualError(t, err, ErrExpiredToken.Error())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, tc.testCase)
	}
}
