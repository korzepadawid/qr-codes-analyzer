package auth

import (
	"database/sql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	mocktoken "github.com/korzepadawid/qr-codes-analyzer/token/mock"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	mockpassword "github.com/korzepadawid/qr-codes-analyzer/util/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSignInAPI(t *testing.T) {
	token := util.RandomString(60)
	user := db.User{
		Username:  util.RandomUsername(),
		Email:     util.RandomMail(),
		FullName:  util.RandomUsername(),
		Password:  hashedPassword,
		CreatedAt: time.Now().Add(-time.Minute),
	}

	requestEmail := signInRequest{
		Username: util.RandomMail(),
		Password: util.RandomString(30),
	}

	requestUsername := signInRequest{
		Username: util.RandomUsername(),
		Password: util.RandomString(30),
	}

	testCases := []struct {
		description   string
		requestBody   signInRequest
		buildStabs    func(t *testing.T, store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, passwordService *mockpassword.MockPasswordService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			description: "should sign in when email provided instead of username",
			requestBody: requestEmail,
			buildStabs: func(t *testing.T, store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, passwordService *mockpassword.MockPasswordService) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(0)
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(1).Return(user, nil)
				passwordService.EXPECT().VerifyPassword(gomock.Eq(user.Password), gomock.Eq(requestEmail.Password)).Times(1).Return(nil)
				tokenProvider.EXPECT().CreateToken(gomock.Eq(user.Username)).Times(1).Return(token, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				util.RequireBodyMatchObject(t, recorder.Body, newSingInResponse(token))
			},
		},
		{
			description: "should sign in when username provided",
			requestBody: requestUsername,
			buildStabs: func(t *testing.T, store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, passwordService *mockpassword.MockPasswordService) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(requestUsername.Username)).Times(0)
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(requestUsername.Username)).Times(1).Return(user, nil)
				passwordService.EXPECT().VerifyPassword(gomock.Eq(user.Password), gomock.Eq(requestUsername.Password)).Times(1).Return(nil)
				tokenProvider.EXPECT().CreateToken(gomock.Eq(user.Username)).Times(1).Return(token, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				util.RequireBodyMatchObject(t, recorder.Body, newSingInResponse(token))
			},
		},
		{
			description: "should return unauthorized when username not found",
			requestBody: requestUsername,
			buildStabs: func(t *testing.T, store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, passwordService *mockpassword.MockPasswordService) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(requestUsername.Username)).Times(0)
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(requestUsername.Username)).Times(1).Return(db.User{}, sql.ErrNoRows)
				passwordService.EXPECT().VerifyPassword(gomock.Eq(user.Password), gomock.Eq(requestUsername.Password)).Times(0)
				tokenProvider.EXPECT().CreateToken(gomock.Eq(user.Username)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return unauthorized when email not found",
			requestBody: requestEmail,
			buildStabs: func(t *testing.T, store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, passwordService *mockpassword.MockPasswordService) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(0)
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(1).Return(db.User{}, sql.ErrNoRows)
				passwordService.EXPECT().VerifyPassword(gomock.Eq(user.Password), gomock.Eq(requestEmail.Password)).Times(0)
				tokenProvider.EXPECT().CreateToken(gomock.Eq(user.Username)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return unauthorized when password mismatched",
			requestBody: requestEmail,
			buildStabs: func(t *testing.T, store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, passwordService *mockpassword.MockPasswordService) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(0)
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(1).Return(user, nil)
				passwordService.EXPECT().VerifyPassword(gomock.Eq(user.Password), gomock.Eq(requestEmail.Password)).Times(1).Return(bcrypt.ErrMismatchedHashAndPassword)
				tokenProvider.EXPECT().CreateToken(gomock.Eq(user.Username)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			description: "should return server internal error when token creation error",
			requestBody: requestEmail,
			buildStabs: func(t *testing.T, store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, passwordService *mockpassword.MockPasswordService) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(0)
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(1).Return(user, nil)
				passwordService.EXPECT().VerifyPassword(gomock.Eq(user.Password), gomock.Eq(requestEmail.Password)).Times(1).Return(nil)
				tokenProvider.EXPECT().CreateToken(gomock.Eq(user.Username)).Times(1).Return("", jwt.ErrInvalidKey)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "should return bad request when missing password",
			requestBody: signInRequest{
				Username: user.Username,
			},
			buildStabs: func(t *testing.T, store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, passwordService *mockpassword.MockPasswordService) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(0)
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(0)
				passwordService.EXPECT().VerifyPassword(gomock.Eq(user.Password), gomock.Eq(requestEmail.Password)).Times(0)
				tokenProvider.EXPECT().CreateToken(gomock.Eq(user.Username)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description: "should return bad request when missing username",
			requestBody: signInRequest{
				Password: requestUsername.Password,
			},
			buildStabs: func(t *testing.T, store *mockdb.MockStore, tokenProvider *mocktoken.MockProvider, passwordService *mockpassword.MockPasswordService) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(0)
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(requestEmail.Username)).Times(0)
				passwordService.EXPECT().VerifyPassword(gomock.Eq(user.Password), gomock.Eq(requestEmail.Password)).Times(0)
				tokenProvider.EXPECT().CreateToken(gomock.Eq(user.Username)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			passwordService := mockpassword.NewMockPasswordService(ctrl)
			tokenProvider := mocktoken.NewMockProvider(ctrl)
			r := setUpHandler(store, tokenProvider, passwordService)

			// building stabs
			tC.buildStabs(t, store, tokenProvider, passwordService)

			// server & response
			request, err := http.NewRequest(http.MethodPost, routerGroupPrefix+signInUrl, util.MarshallBody(tC.requestBody))
			recorder := httptest.NewRecorder()
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}
