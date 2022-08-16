package group

import (
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	mockmaker "github.com/korzepadawid/qr-codes-analyzer/token/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetGroupAPI(t *testing.T) {
	testCases := []struct {
		description         string
		authorizationHeader string
		routeParam          string
		buildStabs          func(*mockdb.MockStore, *mockmaker.MockMaker)
		checkResponse       func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			description:         "should return group when exists",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          fmt.Sprintf("/%d", randomGroupID),
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				arg := db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				}
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), arg).Times(1).Return(mockGroup, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireMatchGroup(t, recorder.Body, mockGroup)
			},
		},
		{
			description:         "should return not found error when not exists",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          fmt.Sprintf("/%d", randomGroupID),
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				arg := db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				}
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), arg).Times(1).Return(db.Group{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			description:         "should return internal error when db failed",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          fmt.Sprintf("/%d", randomGroupID),
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				arg := db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				}
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), arg).Times(1).Return(db.Group{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description:         "should return client error when given string instead of int",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          "/asdfa",
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), gomock.Any()).Times(0)
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
			request, err := http.NewRequest(http.MethodGet, routerGroupPrefix+tC.routeParam, nil)
			request.Header.Set("Authorization", tC.authorizationHeader)
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}
