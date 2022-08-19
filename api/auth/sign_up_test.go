package auth

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	mockpassword "github.com/korzepadawid/qr-codes-analyzer/util/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	hashedPassword = "hashedPassword"
)

func TestSignUpAPI(t *testing.T) {

	user := db.User{
		Username:  util.RandomUsername(),
		Email:     util.RandomMail(),
		FullName:  util.RandomUsername(),
		Password:  hashedPassword,
		CreatedAt: time.Now().Add(-time.Minute),
	}

	body := signUpRequest{
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Password: "str@on@gPa@sdword123",
	}

	createUserParams := db.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Password: hashedPassword,
	}

	testCases := []struct {
		description   string
		requestBody   signUpRequest
		buildStabs    func(t *testing.T, store *mockdb.MockStore, passwordService *mockpassword.MockPasswordService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			description: "creates a new user when user doesn't exist",
			requestBody: body,
			buildStabs: func(t *testing.T, store *mockdb.MockStore, passwordService *mockpassword.MockPasswordService) {
				arg := mapRequestToGetUserByUsernameOrEmailParams(body)
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.User{}, sql.ErrNoRows)
				passwordService.EXPECT().HashPassword(gomock.Eq(body.Password)).Times(1).Return(hashedPassword, nil)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(createUserParams)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				util.RequireBodyMatchObject(t, recorder.Body, mapUserToResponse(user))
			},
		},
		{
			description: "returns an error when user exists",
			requestBody: body,
			buildStabs: func(t *testing.T, store *mockdb.MockStore, passwordService *mockpassword.MockPasswordService) {
				arg := mapRequestToGetUserByUsernameOrEmailParams(body)
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Eq(arg)).Times(1).Return(user, nil)
				passwordService.EXPECT().HashPassword(gomock.Eq(body.Password)).Times(0)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(createUserParams)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description: "returns an error when db connection interrupted when checking existence of user",
			requestBody: body,
			buildStabs: func(t *testing.T, store *mockdb.MockStore, passwordService *mockpassword.MockPasswordService) {
				arg := mapRequestToGetUserByUsernameOrEmailParams(body)
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.User{}, sql.ErrConnDone)
				passwordService.EXPECT().HashPassword(gomock.Eq(body.Password)).Times(0)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(createUserParams)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "returns an error when connection interrupted while inserting a new record",
			requestBody: body,
			buildStabs: func(t *testing.T, store *mockdb.MockStore, passwordService *mockpassword.MockPasswordService) {
				arg := mapRequestToGetUserByUsernameOrEmailParams(body)
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.User{}, sql.ErrNoRows)
				passwordService.EXPECT().HashPassword(gomock.Eq(body.Password)).Times(1).Return(hashedPassword, nil)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(createUserParams)).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "returns an error when bcrypt returns an error",
			requestBody: body,
			buildStabs: func(t *testing.T, store *mockdb.MockStore, passwordService *mockpassword.MockPasswordService) {
				arg := mapRequestToGetUserByUsernameOrEmailParams(body)
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.User{}, sql.ErrNoRows)
				passwordService.EXPECT().HashPassword(gomock.Eq(body.Password)).Times(1).Return(hashedPassword, bcrypt.ErrHashTooShort)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(createUserParams)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "returns an error when email validation fails",
			requestBody: signUpRequest{
				Username: user.Username,
				Email:    "notanemail.com",
				FullName: user.FullName,
				Password: user.Password,
			},
			buildStabs: func(t *testing.T, store *mockdb.MockStore, passwordService *mockpassword.MockPasswordService) {
				arg := mapRequestToGetUserByUsernameOrEmailParams(body)
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Eq(arg)).Times(0)
				passwordService.EXPECT().HashPassword(gomock.Eq(body.Password)).Times(0)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(createUserParams)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description: "returns an error when missing required fields",
			requestBody: signUpRequest{
				Username: user.Username,
				Password: user.Password,
			},
			buildStabs: func(t *testing.T, store *mockdb.MockStore, passwordService *mockpassword.MockPasswordService) {
				arg := mapRequestToGetUserByUsernameOrEmailParams(body)
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Eq(arg)).Times(0)
				passwordService.EXPECT().HashPassword(gomock.Eq(body.Password)).Times(0)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(createUserParams)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description: "returns an error when missing required fields",
			requestBody: signUpRequest{
				Username: user.Username,
				Password: user.Password,
			},
			buildStabs: func(t *testing.T, store *mockdb.MockStore, passwordService *mockpassword.MockPasswordService) {
				arg := mapRequestToGetUserByUsernameOrEmailParams(body)
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Eq(arg)).Times(0)
				passwordService.EXPECT().HashPassword(gomock.Eq(body.Password)).Times(0)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(createUserParams)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description: "returns an error when extra characters",
			requestBody: signUpRequest{
				Username: "f@dsafasdf",
				Email:    user.Email,
				FullName: "f@dsafasdf",
				Password: user.Password,
			},
			buildStabs: func(t *testing.T, store *mockdb.MockStore, passwordService *mockpassword.MockPasswordService) {
				arg := mapRequestToGetUserByUsernameOrEmailParams(body)
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Eq(arg)).Times(0)
				passwordService.EXPECT().HashPassword(gomock.Eq(body.Password)).Times(0)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(createUserParams)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	// creating mocks
	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)
	passwordService := mockpassword.NewMockPasswordService(ctrl)
	r := setUpHandler(store, nil, passwordService)

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {

			// building stabs
			tC.buildStabs(t, store, passwordService)

			// server & response
			request, err := http.NewRequest(http.MethodPost, routerGroupPrefix+signUpUrl, util.MarshallBody(tC.requestBody))
			recorder := httptest.NewRecorder()
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)

			// check response
			tC.checkResponse(t, recorder)
		})
	}
}

func mapRequestToGetUserByUsernameOrEmailParams(request signUpRequest) db.GetUserByUsernameOrEmailParams {
	return db.GetUserByUsernameOrEmailParams{
		Username: request.Username,
		Email:    request.Email,
	}
}
