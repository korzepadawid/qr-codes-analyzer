package group

import (
	"github.com/golang/mock/gomock"
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
