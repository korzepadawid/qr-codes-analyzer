package db

import (
	"context"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: util.RandomUsername(),
		Email:    util.RandomMail(),
		FullName: util.RandomUsername(),
		Password: util.RandomString(7),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Password, user.Password)
	require.NotEmpty(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	u, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, u.Username, user.Username)
	require.Equal(t, u.Email, user.Email)
	require.Equal(t, u.FullName, user.FullName)
	require.Equal(t, u.Password, user.Password)
}
