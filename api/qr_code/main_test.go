package qr_code

import (
	"github.com/gin-gonic/gin"
	"github.com/korzepadawid/qr-codes-analyzer/api/auth"
	"github.com/korzepadawid/qr-codes-analyzer/api/common"
	"github.com/korzepadawid/qr-codes-analyzer/cache"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/encode"
	"github.com/korzepadawid/qr-codes-analyzer/storage"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"os"
	"testing"
	"time"
)

var (
	mockUsername  = util.RandomUsername()
	mockPayload   = &token.Payload{Username: mockUsername}
	randomGroupID = util.RandomInt64(1, 1000)
	mockGroup     = db.Group{
		ID:          randomGroupID,
		Owner:       mockUsername,
		Title:       util.RandomString(25),
		Description: util.RandomString(25),
		CreatedAt:   time.Now().Add(-time.Hour),
	}
	mockCreateQRCodeRequestBody = createQRCodeRequest{
		URL:   "https://stackoverflow.com/",
		Title: util.RandomString(5),
	}
	mockGeneratedQRCode = make([]byte, 0)
	mockSavedQRCode     = db.QrCode{
		Uuid:           util.RandomString(20),
		Owner:          mockUsername,
		GroupID:        randomGroupID,
		UsagesCount:    0,
		RedirectionUrl: mockCreateQRCodeRequestBody.URL,
		Title:          mockCreateQRCodeRequestBody.Title,
		Description:    mockCreateQRCodeRequestBody.Description,
		StorageUrl:     util.RandomString(20),
		CreatedAt:      time.Time{},
	}
)

func newMockQRCodeHandler(
	store db.Store,
	config config.Config,
	fileStorage storage.FileStorage,
	qrCodeEncoder encode.Encoder,
	cache cache.Cache,
	tokenProvider token.Provider,
) *gin.Engine {
	r := gin.Default()
	common.SetUpErrorHandler(r)
	h := NewQRCodeHandler(store, config, fileStorage, qrCodeEncoder, cache, auth.SecureRoute(tokenProvider))
	h.RegisterRoutes(r)
	return r
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
