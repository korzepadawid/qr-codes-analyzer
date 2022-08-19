package group

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	mocktoken "github.com/korzepadawid/qr-codes-analyzer/token/mock"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateGroupAPI(t *testing.T) {
	testCases := []struct {
		description         string
		authorizationHeader string
		requestBody         createGroupRequest
		buildStabs          func(*mockdb.MockStore, *mocktoken.MockProvider)
		checkResponse       func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			description:         "should create a new group for encode-codes when ok",
			authorizationHeader: validAuthorizationHeader,
			requestBody:         mockCreateGroupRequestBody,
			buildStabs: func(store *mockdb.MockStore, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().CreateGroup(gomock.Any(), gomock.Eq(mapRequestToParams(mockPayload.Username, mockCreateGroupRequestBody))).Times(1).Return(mockGroup, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireMatchGroup(t, recorder.Body, mockGroup)
			},
		},
		{
			description:         "should return an internal error when error during db insertion",
			authorizationHeader: validAuthorizationHeader,
			requestBody:         mockCreateGroupRequestBody,
			buildStabs: func(store *mockdb.MockStore, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().CreateGroup(gomock.Any(), gomock.Eq(mapRequestToParams(mockPayload.Username, mockCreateGroupRequestBody))).Times(1).Return(db.Group{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description:         "should return an client error when missing title of the group (the only one required field)",
			authorizationHeader: validAuthorizationHeader,
			requestBody:         createGroupRequest{},
			buildStabs: func(store *mockdb.MockStore, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().CreateGroup(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description:         "should return an client error when too long title",
			authorizationHeader: validAuthorizationHeader,
			requestBody: createGroupRequest{
				Title: util.RandomString(256),
			},
			buildStabs: func(store *mockdb.MockStore, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().CreateGroup(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description:         "should return an client error when too long description",
			authorizationHeader: validAuthorizationHeader,
			requestBody: createGroupRequest{
				Title:       util.RandomString(255),
				Description: util.RandomString(256),
			},
			buildStabs: func(store *mockdb.MockStore, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().CreateGroup(gomock.Any(), gomock.Any()).Times(0)
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
			mockMaker := mocktoken.NewMockProvider(ctrl)
			mockStore := mockdb.NewMockStore(ctrl)

			// stabs
			r := newMockGroupHandler(mockStore, mockMaker)
			tC.buildStabs(mockStore, mockMaker)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, routerGroupPrefix, util.MarshallBody(tC.requestBody))
			request.Header.Set("Authorization", tC.authorizationHeader)
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)
			//
			tC.checkResponse(t, recorder)
		})
	}
}
