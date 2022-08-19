package auth

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	mocktoken "github.com/korzepadawid/qr-codes-analyzer/token/mock"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testToken = "djkfsgjdfshgkjdshfgjkhdfjghsdfkjghsdfkjghsdfjkghsdjkfhgdjfkghdfjkghdsfjgkhjdf"
)

var (
	testUsername = util.RandomUsername()
)

type testHeader struct {
	key   string
	value string
}

func TestAuthorizationMiddleware(t *testing.T) {
	testCases := []struct {
		description   string
		header        testHeader
		buildStabs    func(t *testing.T, maker *mocktoken.MockProvider)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			description: "should return 401 when no Authorization header",
			buildStabs: func(t *testing.T, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Eq("")).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return 401 when Authorization header with empty string",
			header: testHeader{
				key:   authorizationHeaderKey,
				value: "",
			},
			buildStabs: func(t *testing.T, tokenProvider *mocktoken.MockProvider) {
				tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return 401 when invalid auth type",
			header: testHeader{
				key:   authorizationHeaderKey,
				value: "BAarer",
			},
			buildStabs: func(t *testing.T, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return 401 when invalid auth type (with space)",
			header: testHeader{
				key:   authorizationHeaderKey,
				value: "BAarer ",
			},
			buildStabs: func(t *testing.T, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return 401 when invalid auth type (with 2 spaces)",
			header: testHeader{
				key:   authorizationHeaderKey,
				value: "BAer  ",
			},
			buildStabs: func(t *testing.T, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return 401 when invalid Bearer token format",
			header: testHeader{
				key:   authorizationHeaderKey,
				value: fmt.Sprintf("%s %s %s", authorizationType, "somerandomtext", "anotherrandom text"),
			},
			buildStabs: func(t *testing.T, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return 401 when invalid Bearer token format",
			header: testHeader{
				key:   authorizationHeaderKey,
				value: fmt.Sprintf("%s  %s %s", authorizationType, "somerandomtext", "anotherrandom text"),
			},
			buildStabs: func(t *testing.T, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return 401 when invalid Bearer token format",
			header: testHeader{
				key:   authorizationHeaderKey,
				value: fmt.Sprintf("%s  %s", authorizationType, "theregoestoken"),
			},
			buildStabs: func(t *testing.T, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return 401 when invalid token",
			header: testHeader{
				key:   authorizationHeaderKey,
				value: fmt.Sprintf("%s %s", authorizationType, testToken),
			},
			buildStabs: func(t *testing.T, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Eq(testToken)).Times(1).Return(&token.Payload{}, token.ErrInvalidToken)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return 401 when expired token",
			header: testHeader{
				key:   authorizationHeaderKey,
				value: fmt.Sprintf("%s %s", authorizationType, testToken),
			},
			buildStabs: func(t *testing.T, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Eq(testToken)).Times(1).Return(&token.Payload{}, token.ErrExpiredToken)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return 200 and username when valid token",
			header: testHeader{
				key:   authorizationHeaderKey,
				value: fmt.Sprintf("%s %s", authorizationType, testToken),
			},
			buildStabs: func(t *testing.T, maker *mocktoken.MockProvider) {
				maker.EXPECT().VerifyToken(gomock.Eq(testToken)).Times(1).Return(&token.Payload{
					Username: testUsername,
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				util.RequireBodyMatchObject(t, recorder.Body, testResponseForSecuredRoute{testUsername})
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			// mocks
			ctrl := gomock.NewController(t)
			mockMaker := mocktoken.NewMockProvider(ctrl)

			// build stabs
			tC.buildStabs(t, mockMaker)
			r := setUpAuthMiddleware(mockMaker)
			request, err := http.NewRequest(http.MethodGet, securedTestRouteUrl, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			request.Header.Set(tC.header.key, tC.header.value)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}
