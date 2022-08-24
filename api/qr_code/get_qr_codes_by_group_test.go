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

func TestGetQRCodesAPI(t *testing.T) {
	testCases := []struct {
		description   string
		groupID       int64
		queryString   string
		buildStabs    func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			description: "should return a valid page response when ok",
			groupID:     mockGroup.ID,
			queryString: "?page_size=11&page_number=1",
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetQRCodesPageByGroupAndOwner(gomock.Any(), gomock.Eq(db.GetQRCodesPageByGroupAndOwnerParams{
					Limit:   11,
					Offset:  0,
					GroupID: mockGroup.ID,
					Owner:   mockPayload.Username,
				})).Times(1).Return(getQRCodesSlice(11), nil)
				store.EXPECT().GetQRCodesCountByGroupAndOwner(gomock.Any(), gomock.Eq(db.GetQRCodesCountByGroupAndOwnerParams{
					GroupID: mockGroup.ID,
					Owner:   mockPayload.Username,
				})).Times(1).Return(int64(12), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				pageResponse := parseQRCodePageResponse(t, recorder.Body)
				require.Equal(t, int32(11), pageResponse.ItemsPerPage)
				require.Equal(t, int32(1), pageResponse.PageNumber)
				require.Equal(t, int32(2), pageResponse.LastPage)
				require.Equal(t, 11, len(pageResponse.Data))
			},
		},
		{
			description: "should return an error when db fails while getting qr codes page",
			groupID:     mockGroup.ID,
			queryString: "?page_size=11&page_number=1",
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetQRCodesPageByGroupAndOwner(gomock.Any(), gomock.Eq(db.GetQRCodesPageByGroupAndOwnerParams{
					Limit:   11,
					Offset:  0,
					GroupID: mockGroup.ID,
					Owner:   mockPayload.Username,
				})).Times(1).Return(make([]db.QrCode, 0), sql.ErrConnDone)
				store.EXPECT().GetQRCodesCountByGroupAndOwner(gomock.Any(), gomock.Eq(db.GetQRCodesCountByGroupAndOwnerParams{
					GroupID: mockGroup.ID,
					Owner:   mockPayload.Username,
				})).Times(1).Return(int64(12), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "should return an error when db fails during both queries",
			groupID:     mockGroup.ID,
			queryString: "?page_size=11&page_number=1",
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetQRCodesPageByGroupAndOwner(gomock.Any(), gomock.Eq(db.GetQRCodesPageByGroupAndOwnerParams{
					Limit:   11,
					Offset:  0,
					GroupID: mockGroup.ID,
					Owner:   mockPayload.Username,
				})).Times(1).Return(make([]db.QrCode, 0), sql.ErrConnDone)
				store.EXPECT().GetQRCodesCountByGroupAndOwner(gomock.Any(), gomock.Eq(db.GetQRCodesCountByGroupAndOwnerParams{
					GroupID: mockGroup.ID,
					Owner:   mockPayload.Username,
				})).Times(1).Return(int64(0), sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "should return an error when non-positive page number",
			groupID:     mockGroup.ID,
			queryString: "?page_size=11&page_number=0",
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetQRCodesPageByGroupAndOwner(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetQRCodesCountByGroupAndOwner(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description: "should return an error when page number less than 10",
			groupID:     mockGroup.ID,
			queryString: "?page_size=9&page_number=1",
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetQRCodesPageByGroupAndOwner(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetQRCodesCountByGroupAndOwner(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description: "should return an error when page number negative",
			groupID:     mockGroup.ID,
			queryString: "?page_size=-1&page_number=1",
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetQRCodesPageByGroupAndOwner(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetQRCodesCountByGroupAndOwner(gomock.Any(), gomock.Any()).Times(0)
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

			// stabs
			r := newMockQRCodeHandler(mockStore, config.Config{}, nil, nil, nil, mockTokenProvider, nil)
			tC.buildStabs(mockStore, mockTokenProvider)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/groups/%d/qr-codes%s", tC.groupID, tC.queryString), nil)
			request.Header.Set("Authorization", "Bearer validtoken")
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}

func getQRCodesSlice(n int) []db.QrCode {
	res := make([]db.QrCode, n)
	for i := 0; i < n; i++ {
		res[i] = mockSavedQRCode
	}
	return res
}
