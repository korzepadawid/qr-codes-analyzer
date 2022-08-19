package group

import (
	"fmt"
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

func TestUpdateGroupAPI(t *testing.T) {
	testCases := []struct {
		description         string
		authorizationHeader string
		routeParam          string
		requestBody         updateGroupRequest
		buildStabs          func(*mockdb.MockStore, *mocktoken.MockProvider)
		checkResponse       func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			description:         "should create transaction with valid args when title and description both given",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          fmt.Sprintf("/%d", randomGroupID),
			requestBody: updateGroupRequest{
				Title:       mockGroupUpdateTitle,       // not empty
				Description: mockGroupUpdateDescription, // not empty
			},
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				arg := db.UpdateGroupTxParams{
					Title:       mockGroupUpdateTitle,
					Description: mockGroupUpdateDescription,
					Owner:       mockPayload.Username,
					ID:          randomGroupID,
				}
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().UpdateGroupTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(mockGroup, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireMatchGroup(t, recorder.Body, mockGroup)
			},
		},
		{
			description:         "should create transaction with empty title when title not given",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          fmt.Sprintf("/%d", randomGroupID),
			requestBody: updateGroupRequest{
				Title:       "          ",
				Description: mockGroupUpdateDescription, // not empty
			},
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				arg := db.UpdateGroupTxParams{
					Title:       "",
					Description: mockGroupUpdateDescription,
					Owner:       mockPayload.Username,
					ID:          randomGroupID,
				}
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().UpdateGroupTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(mockGroup, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireMatchGroup(t, recorder.Body, mockGroup)
			},
		},
		{
			description:         "should create transaction with empty description when description is blank",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          fmt.Sprintf("/%d", randomGroupID),
			requestBody: updateGroupRequest{
				Title:       mockGroupUpdateTitle, // not empty
				Description: "                ",
			},
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				arg := db.UpdateGroupTxParams{
					Title:       mockGroupUpdateTitle,
					Description: "",
					Owner:       mockPayload.Username,
					ID:          randomGroupID,
				}
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().UpdateGroupTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(mockGroup, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireMatchGroup(t, recorder.Body, mockGroup)
			},
		},
		{
			description:         "should create transaction with empty description when description is blank",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          fmt.Sprintf("/%d", randomGroupID),
			requestBody: updateGroupRequest{
				Title:       mockGroupUpdateTitle, // not empty
				Description: "                ",
			},
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				arg := db.UpdateGroupTxParams{
					Title:       mockGroupUpdateTitle,
					Description: "",
					Owner:       mockPayload.Username,
					ID:          randomGroupID,
				}
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().UpdateGroupTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(mockGroup, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireMatchGroup(t, recorder.Body, mockGroup)
			},
		},
		{
			description:         "should return an error when too long description",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          fmt.Sprintf("/%d", randomGroupID),
			requestBody: updateGroupRequest{
				Title:       mockGroupUpdateTitle, // not empty
				Description: util.RandomString(256),
			},
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().UpdateGroupTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description:         "should return an error when too long title",
			authorizationHeader: validAuthorizationHeader,
			routeParam:          fmt.Sprintf("/%d", randomGroupID),
			requestBody: updateGroupRequest{
				Title:       util.RandomString(256),
				Description: mockGroupUpdateDescription, // not empty
			},
			buildStabs: func(store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider) {
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				store.EXPECT().UpdateGroupTx(gomock.Any(), gomock.Any()).Times(0)
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
			request, err := http.NewRequest(http.MethodPatch, routerGroupPrefix+tC.routeParam, util.MarshallBody(tC.requestBody))
			request.Header.Set("Authorization", tC.authorizationHeader)
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}
