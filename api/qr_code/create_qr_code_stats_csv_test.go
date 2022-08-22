package qr_code

import (
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	mocktoken "github.com/korzepadawid/qr-codes-analyzer/token/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateQRCodeStatsCSV(t *testing.T) {
	testCases := []struct {
		description   string
		qrCodeUUID    string
		buildStabs    func(store *mockdb.MockStore, provider *mocktoken.MockProvider)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			description: "should set http headers when ok",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			buildStabs: func(store *mockdb.MockStore, provider *mocktoken.MockProvider) {
				provider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetQRCodeRedirectEntries(gomock.Any(), gomock.Eq(db.GetQRCodeRedirectEntriesParams{
					Uuid:  mockSavedQRCode.Uuid,
					Owner: mockPayload.Username,
				})).Times(1).Return(make([]db.GetQRCodeRedirectEntriesRow, 0), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, fmt.Sprintf("attachment; filename=stats-%s.csv", mockSavedQRCode.Uuid), recorder.Header().Get("Content-Disposition"))
				require.Equal(t, "application/octet-stream", recorder.Header().Get("Content-Type"))
				require.NotEmpty(t, recorder.Header().Get("Content-Length"))
			},
		},
		{
			description: "should return an error when db fails",
			qrCodeUUID:  mockSavedQRCode.Uuid,
			buildStabs: func(store *mockdb.MockStore, provider *mocktoken.MockProvider) {
				provider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetQRCodeRedirectEntries(gomock.Any(), gomock.Eq(db.GetQRCodeRedirectEntriesParams{
					Uuid:  mockSavedQRCode.Uuid,
					Owner: mockPayload.Username,
				})).Times(1).Return(make([]db.GetQRCodeRedirectEntriesRow, 0), sql.ErrConnDone)
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
			r := newMockQRCodeHandler(mockStore, config.Config{}, nil, nil, nil, mockTokenProvider)
			tC.buildStabs(mockStore, mockTokenProvider)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/qr-codes/%s/stats/csv", tC.qrCodeUUID), nil)
			request.Header.Set("Authorization", "Bearer validtoken")
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}
