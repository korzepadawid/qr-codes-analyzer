package group

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/common"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const (
	validAuthorizationHeader = "Bearer tokengoeshere"
)

var (
	mockUsername               = util.RandomUsername()
	mockPayload                = &token.Payload{Username: mockUsername}
	mockCreateGroupRequestBody = createGroupRequest{
		Title:       util.RandomString(10),
		Description: util.RandomString(10),
	}
	mockGroup = db.Group{
		ID:          util.RandomInt64(1, 100),
		Owner:       mockUsername,
		Title:       mockCreateGroupRequestBody.Title,
		Description: mockCreateGroupRequestBody.Description,
		CreatedAt:   time.Now().Add(-time.Hour),
	}
)

func newMockGroupHandler(store db.Store, maker token.Maker) *gin.Engine {
	r := gin.Default()
	common.SetUpErrorHandler(r)
	newGroupHandler := NewGroupHandler(store, auth.SecureRoute(maker))
	newGroupHandler.RegisterRoutes(r)
	return r
}

func requireMatchGroup(t *testing.T, b *bytes.Buffer, expect db.Group) {
	data, err := ioutil.ReadAll(b)
	require.NoError(t, err)
	var got db.Group
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	require.Equal(t, expect.Description, got.Description)
	require.Equal(t, expect.Title, got.Title)
	require.Equal(t, expect.ID, got.ID)
	require.Equal(t, expect.Owner, got.Owner)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
