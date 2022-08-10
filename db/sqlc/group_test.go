package db

import (
	"context"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomGroup(t *testing.T, user User) Group {
	arg := CreateGroupParams{
		Owner:       user.Username,
		Title:       util.RandomString(6),
		Description: util.RandomString(20),
	}

	group, err := testQueries.CreateGroup(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, group)
	require.Equal(t, group.Owner, user.Username)
	require.Equal(t, group.Title, arg.Title)
	require.Equal(t, group.Description, arg.Description)
	require.NotZero(t, group.CreatedAt)
	require.NotZero(t, group.ID)

	return group
}

func TestCreateGroup(t *testing.T) {
	user := createRandomUser(t)
	createRandomGroup(t, user)
}

func TestGetGroupsByOwner(t *testing.T) {
	owner := createRandomUser(t)
	otherUser := createRandomUser(t)

	for i := 0; i < 10; i++ {
		createRandomGroup(t, owner)
	}

	arg := GetGroupsByOwnerParams{
		Limit:  20,
		Offset: 0,
		Owner:  owner.Username,
	}

	groups, err := testQueries.GetGroupsByOwner(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, groups, 10)

	for _, g := range groups {
		require.NotEmpty(t, g)
		require.Equal(t, g.Owner, owner.Username)
		require.NotEqual(t, g.Owner, otherUser.Username)
		require.NotZero(t, g.CreatedAt)
		require.NotZero(t, g.ID)
	}
}
