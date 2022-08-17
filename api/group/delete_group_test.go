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

func TestDeleteGroupAPI(t *testing.T) {
	testCases := []struct {
		description         string
		authorizationHeader string
		routeParam          string
		buildStabs          func(*mockdb.MockStore, *mockmaker.MockMaker)
		checkResponse       func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			description:         "should perform db query when ok",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          fmt.Sprintf("/%d", randomGroupID),
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				arg := db.DeleteGroupByOwnerAndIDParams{
					GroupID: randomGroupID,
					Owner:   mockPayload.Username,
				}
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().DeleteGroupByOwnerAndID(gomock.Any(), arg).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			description:         "should perform db query when ok",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          fmt.Sprintf("/%d", randomGroupID),
			buildStabs: func(store *mockdb.MockStore, maker *mockmaker.MockMaker) {
				arg := db.DeleteGroupByOwnerAndIDParams{
					GroupID: randomGroupID,
					Owner:   mockPayload.Username,
				}
				maker.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().DeleteGroupByOwnerAndID(gomock.Any(), arg).Times(1).Return(sql.ErrConnDone)
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
			mockMaker := mockmaker.NewMockMaker(ctrl)
			mockStore := mockdb.NewMockStore(ctrl)

			// stabs
			r := newMockGroupHandler(mockStore, mockMaker)
			tC.buildStabs(mockStore, mockMaker)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodDelete, routerGroupPrefix+tC.routeParam, nil)
			request.Header.Set("Authorization", tC.authorizationHeader)
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}
