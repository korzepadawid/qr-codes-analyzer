package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"github.com/stretchr/testify/require"
)

func TestSignUpAPI(t *testing.T) {
	testCases := []struct {
		descrption    string
		requestBody   signUpRequest
		buildStabs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			descrption: "should save user and return its representation",
			requestBody: signUpRequest{
				username: util.RandomUsername(),
				email:    util.RandomMail(),
				fullName: util.RandomUsername(),
				password: util.RandomString(8),
			},
			buildStabs:    func(store *mockdb.MockStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.descrption, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)

			// build stabs
			tc.buildStabs(mockStore)

			server := newMockServer(t, mockStore)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/auth/signup", util.MarshallBody(tc.requestBody))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			// check response(recorder)
			tc.checkResponse(recorder)
		})
	}
}
