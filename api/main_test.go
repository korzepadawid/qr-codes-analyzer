package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/stretchr/testify/require"
)

func newMockServer(t *testing.T, store db.Store) *Server {
	server, err := NewServer(config.Config{}, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
