package qr_code

import (
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/korzepadawid/qr-codes-analyzer/api/errors"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	mocktoken "github.com/korzepadawid/qr-codes-analyzer/token/mock"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateQRCodeAPI(t *testing.T) {

	testCases := []struct {
		description   string
		qrCodeUUID    string
		requestBody   UpdateQRCodeRequest
		buildStabs    func(store *mockdb.MockStore, provider *mocktoken.MockProvider)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			description: "should return no error when ok",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			requestBody: mockUpdateQRCodeRequest,
			buildStabs: func(store *mockdb.MockStore, provider *mocktoken.MockProvider) {
				provider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().UpdateQRCodeTx(gomock.Any(), gomock.Eq(
					db.UpdateQRCodeTxParams{
						UUID:        mockSavedQRCode.Uuid,
						Owner:       mockPayload.Username,
						Title:       mockUpdateQRCodeRequest.Title,
						Description: mockUpdateQRCodeRequest.Description,
					})).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			description: "should return no error when when need to trim fields",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			requestBody: UpdateQRCodeRequest{
				Title:       "           a          ",
				Description: "          ",
			},
			buildStabs: func(store *mockdb.MockStore, provider *mocktoken.MockProvider) {
				provider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().UpdateQRCodeTx(gomock.Any(), gomock.Eq(
					db.UpdateQRCodeTxParams{
						UUID:        mockSavedQRCode.Uuid,
						Owner:       mockPayload.Username,
						Title:       "a",
						Description: "",
					})).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			description: "should return an error when qr code not found",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			requestBody: mockUpdateQRCodeRequest,
			buildStabs: func(store *mockdb.MockStore, provider *mocktoken.MockProvider) {
				provider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().UpdateQRCodeTx(gomock.Any(), gomock.Eq(
					db.UpdateQRCodeTxParams{
						UUID:        mockSavedQRCode.Uuid,
						Owner:       mockPayload.Username,
						Title:       mockUpdateQRCodeRequest.Title,
						Description: mockUpdateQRCodeRequest.Description,
					})).Times(1).Return(errors.ErrQRCodeNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			description: "should return an error when db failed",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			requestBody: mockUpdateQRCodeRequest,
			buildStabs: func(store *mockdb.MockStore, provider *mocktoken.MockProvider) {
				provider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().UpdateQRCodeTx(gomock.Any(), gomock.Eq(
					db.UpdateQRCodeTxParams{
						UUID:        mockSavedQRCode.Uuid,
						Owner:       mockPayload.Username,
						Title:       mockUpdateQRCodeRequest.Title,
						Description: mockUpdateQRCodeRequest.Description,
					})).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			// mocks
			ctrl := gomock.NewController(t)
			mockStore := mockdb.NewMockStore(ctrl)
			mockTokenProvider := mocktoken.NewMockProvider(ctrl)

			// stabs
			r := newMockQRCodeHandler(mockStore, config.Config{}, nil, nil, nil, mockTokenProvider, nil)
			tC.buildStabs(mockStore, mockTokenProvider)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/qr-codes/%s", tC.qrCodeUUID), util.MarshallBody(tC.requestBody))
			request.Header.Set("Authorization", "Bearer validtoken")
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}
