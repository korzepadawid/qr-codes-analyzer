package qr_code

import (
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	mockstorage "github.com/korzepadawid/qr-codes-analyzer/storage/mock"
	mocktoken "github.com/korzepadawid/qr-codes-analyzer/token/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteQRCodeAPI(t *testing.T) {
	testCases := []struct {
		description   string
		qrCodeUUID    string
		buildStabs    func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, storage *mockstorage.MockFileStorage)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			description: "should delete qr code when ok",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, storage *mockstorage.MockFileStorage) {
				params := db.DeleteQRCodeParams{
					Uuid:  mockSavedQRCode.Uuid,
					Owner: mockPayload.Username,
				}
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Return(mockPayload, nil)
				store.EXPECT().DeleteQRCode(gomock.Any(), gomock.Eq(params)).Times(1).Return(nil)
				storage.EXPECT().DeleteFile(gomock.Any(), gomock.Eq(mockSavedQRCode.Uuid+".png")).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			description: "should return an error when db fails",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, storage *mockstorage.MockFileStorage) {
				params := db.DeleteQRCodeParams{
					Uuid:  mockSavedQRCode.Uuid,
					Owner: mockPayload.Username,
				}
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Return(mockPayload, nil)
				store.EXPECT().DeleteQRCode(gomock.Any(), gomock.Eq(params)).Times(1).Return(sql.ErrConnDone)
				storage.EXPECT().DeleteFile(gomock.Any(), gomock.Eq(mockSavedQRCode.Uuid+".png")).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "should return an error when empty uuid",
			qrCodeUUID:  "            ",
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, storage *mockstorage.MockFileStorage) {
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Return(mockPayload, nil)
				store.EXPECT().DeleteQRCode(gomock.Any(), gomock.Any()).Times(0)
				storage.EXPECT().DeleteFile(gomock.Any(), gomock.Eq(".png")).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			// mocks
			ctrl := gomock.NewController(t)
			mockStore := mockdb.NewMockStore(ctrl)
			mockTokenProvider := mocktoken.NewMockProvider(ctrl)
			mockFileStorage := mockstorage.NewMockFileStorage(ctrl)

			// stabs
			r := newMockQRCodeHandler(mockStore, config.Config{}, mockFileStorage, nil, nil, mockTokenProvider)
			tC.buildStabs(mockStore, mockTokenProvider, mockFileStorage)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/qr-codes/%s", tC.qrCodeUUID), nil)
			request.Header.Set("Authorization", "Bearer validtoken")
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}
