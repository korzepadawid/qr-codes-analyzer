package qr_code

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/golang/mock/gomock"
	c "github.com/korzepadawid/qr-codes-analyzer/cache"
	mockcache "github.com/korzepadawid/qr-codes-analyzer/cache/mock"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRedirectQRCode(t *testing.T) {
	testCases := []struct {
		description   string
		qrCodeUUID    string
		buildStabs    func(store *mockdb.MockStore, cache *mockcache.MockCache)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			description: "should redirect when using url from cache",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			buildStabs: func(_ *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().Get(gomock.Eq(mockSavedQRCode.Uuid)).Times(1).Return(mockSavedQRCode.StorageUrl, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusPermanentRedirect, recorder.Code)
			},
		},
		{
			description: "should cache url and redirect when using db instead of cache",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			buildStabs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().Get(gomock.Eq(mockSavedQRCode.Uuid)).Times(1).Return("", redis.Nil)
				store.EXPECT().GetQRCode(gomock.Any(), gomock.Eq(mockSavedQRCode.Uuid)).Times(1).Return(mockSavedQRCode, nil)
				cache.EXPECT().Set(gomock.Eq(&c.SetParams{
					Key:      mockSavedQRCode.Uuid,
					Value:    mockSavedQRCode.RedirectionUrl,
					Duration: time.Minute * 2,
				})).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusPermanentRedirect, recorder.Code)
			},
		},
		{
			description: "should return an error when url not found either in cache or db",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			buildStabs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().Get(gomock.Eq(mockSavedQRCode.Uuid)).Times(1).Return("", redis.Nil)
				store.EXPECT().GetQRCode(gomock.Any(), gomock.Eq(mockSavedQRCode.Uuid)).Times(1).Return(db.QrCode{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			description: "should return an error when db failed to query",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			buildStabs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().Get(gomock.Eq(mockSavedQRCode.Uuid)).Times(1).Return("", redis.Nil)
				store.EXPECT().GetQRCode(gomock.Any(), gomock.Eq(mockSavedQRCode.Uuid)).Times(1).Return(db.QrCode{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "should return an error when empty uuid",
			qrCodeUUID:  "",
			buildStabs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().Get(gomock.Eq("")).Times(0)
				store.EXPECT().GetQRCode(gomock.Any(), gomock.Eq("")).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			description: "should return an error when trimmed uuid is empty)",
			qrCodeUUID:  "       ",
			buildStabs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
				cache.EXPECT().Get(gomock.Eq("")).Times(0)
				store.EXPECT().GetQRCode(gomock.Any(), gomock.Eq("")).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			// mocks
			ctrl := gomock.NewController(t)
			mockStore := mockdb.NewMockStore(ctrl)
			mockCache := mockcache.NewMockCache(ctrl)

			// stabs
			r := newMockQRCodeHandler(mockStore, config.Config{}, nil, nil, mockCache, nil)
			tC.buildStabs(mockStore, mockCache)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf(
					"/qr-codes/%s/redirect",
					tC.qrCodeUUID,
				),
				nil,
			)
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}
