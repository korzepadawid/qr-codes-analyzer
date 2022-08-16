package group

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/korzepadawid/qr-codes-analyzer/api/common"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	mockmaker "github.com/korzepadawid/qr-codes-analyzer/token/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	groups = getGroups(10)
)

func TestGetGroupsByOwnerAPI(t *testing.T) {
	testCases := []struct {
		description         string
		authorizationHeader string
		queryParams         string
		buildStabs          func(*mockdb.MockStore, *mockmaker.MockMaker)
		checkResponse       func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			description:         "should return first page when first page requested",
			authorizationHeader: validAuthorizationHeader,
			queryParams:         "?page_size=10&page_number=1",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				arg := db.GetGroupsByOwnerParams{
					Limit:  10,
					Offset: 0,
					Owner:  mockPayload.Username,
				}
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupsCountByOwner(gomock.Any(), mockPayload.Username).Times(1).Return(groupsCount, nil)
				store.EXPECT().GetGroupsByOwner(gomock.Any(), arg).Times(1).Return(groups, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				pageResponse := parseGroupPageResponse(t, recorder.Body)
				response := common.NewPageResponse(1, 10, 3, groups)
				require.Equal(t, pageResponse.PageNumber, response.PageNumber)
				require.Equal(t, pageResponse.ItemsPerPage, response.ItemsPerPage)
				require.Equal(t, pageResponse.LastPage, response.LastPage)
				require.Equal(t, len(pageResponse.Data), 10)
			},
		},
		{
			description:         "should return first page and valid last page number when first page requested",
			authorizationHeader: validAuthorizationHeader,
			queryParams:         "?page_size=10&page_number=1",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				arg := db.GetGroupsByOwnerParams{
					Limit:  10,
					Offset: 0,
					Owner:  mockPayload.Username,
				}
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupsCountByOwner(gomock.Any(), mockPayload.Username).Times(1).Return(int64(len(groups)), nil)
				store.EXPECT().GetGroupsByOwner(gomock.Any(), arg).Times(1).Return(groups, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				pageResponse := parseGroupPageResponse(t, recorder.Body)
				response := common.NewPageResponse(1, 10, 1, groups)
				require.Equal(t, pageResponse.PageNumber, response.PageNumber)
				require.Equal(t, pageResponse.ItemsPerPage, response.ItemsPerPage)
				require.Equal(t, pageResponse.LastPage, response.LastPage)
				require.Equal(t, len(pageResponse.Data), 10)
			},
		},
		{
			description:         "should return second page and valid last page number when second page requested",
			authorizationHeader: validAuthorizationHeader,
			queryParams:         "?page_size=10&page_number=2",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				arg := db.GetGroupsByOwnerParams{
					Limit:  10,
					Offset: 10,
					Owner:  mockPayload.Username,
				}
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupsCountByOwner(gomock.Any(), mockPayload.Username).Times(1).Return(groupsCount, nil)
				store.EXPECT().GetGroupsByOwner(gomock.Any(), arg).Times(1).Return(groups, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				pageResponse := parseGroupPageResponse(t, recorder.Body)
				response := common.NewPageResponse(2, 10, 3, groups)
				require.Equal(t, pageResponse.PageNumber, response.PageNumber)
				require.Equal(t, pageResponse.ItemsPerPage, response.ItemsPerPage)
				require.Equal(t, pageResponse.LastPage, response.LastPage)
				require.Equal(t, len(pageResponse.Data), 10)
			},
		},
		{
			description:         "should return error when db error while getting page",
			authorizationHeader: validAuthorizationHeader,
			queryParams:         "?page_size=10&page_number=2",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				arg := db.GetGroupsByOwnerParams{
					Limit:  10,
					Offset: 10,
					Owner:  mockPayload.Username,
				}
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupsCountByOwner(gomock.Any(), mockPayload.Username).Times(1).Return(groupsCount, nil)
				store.EXPECT().GetGroupsByOwner(gomock.Any(), arg).Times(1).Return([]db.Group{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description:         "should return error when db error while counting rows",
			authorizationHeader: validAuthorizationHeader,
			queryParams:         "?page_size=10&page_number=2",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				arg := db.GetGroupsByOwnerParams{
					Limit:  10,
					Offset: 10,
					Owner:  mockPayload.Username,
				}
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupsCountByOwner(gomock.Any(), mockPayload.Username).Times(1).Return(int64(0), sql.ErrConnDone)
				store.EXPECT().GetGroupsByOwner(gomock.Any(), arg).Times(1).Return(groups, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description:         "should return error when errors in both goroutines",
			authorizationHeader: validAuthorizationHeader,
			queryParams:         "?page_size=10&page_number=2",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				arg := db.GetGroupsByOwnerParams{
					Limit:  10,
					Offset: 10,
					Owner:  mockPayload.Username,
				}
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupsCountByOwner(gomock.Any(), mockPayload.Username).Times(1).Return(int64(0), sql.ErrConnDone)
				store.EXPECT().GetGroupsByOwner(gomock.Any(), arg).Times(1).Return([]db.Group{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description:         "should return error when missing required param",
			authorizationHeader: validAuthorizationHeader,
			queryParams:         "?page_number=2",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupsCountByOwner(gomock.Any(), mockPayload.Username).Times(0)
				store.EXPECT().GetGroupsByOwner(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description:         "should return error when negative page number",
			authorizationHeader: validAuthorizationHeader,
			queryParams:         "?page_number=-1&page_size=10",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupsCountByOwner(gomock.Any(), mockPayload.Username).Times(0)
				store.EXPECT().GetGroupsByOwner(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description:         "should return error when negative page size",
			authorizationHeader: validAuthorizationHeader,
			queryParams:         "?page_number=2&page_size=-10",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupsCountByOwner(gomock.Any(), mockPayload.Username).Times(0)
				store.EXPECT().GetGroupsByOwner(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description:         "should return error when page size is equal zero",
			authorizationHeader: validAuthorizationHeader,
			queryParams:         "?page_number=2&page_size=0",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupsCountByOwner(gomock.Any(), mockPayload.Username).Times(0)
				store.EXPECT().GetGroupsByOwner(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description:         "should return error when page number is equal zero",
			authorizationHeader: validAuthorizationHeader,
			queryParams:         "?page_number=0&page_size=11",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupsCountByOwner(gomock.Any(), mockPayload.Username).Times(0)
				store.EXPECT().GetGroupsByOwner(gomock.Any(), gomock.Any()).Times(0)
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
			mockMaker := mockmaker.NewMockMaker(ctrl)
			mockStore := mockdb.NewMockStore(ctrl)

			// stabs
			r := newMockGroupHandler(mockStore, mockMaker)
			tC.buildStabs(mockStore, mockMaker)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, routerGroupPrefix+tC.queryParams, nil)
			request.Header.Set("Authorization", tC.authorizationHeader)
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}

func getGroups(n int) []db.Group {
	g := make([]db.Group, n)
	for i := 0; i < n; i++ {
		g[i] = mockGroup
	}
	return g
}
