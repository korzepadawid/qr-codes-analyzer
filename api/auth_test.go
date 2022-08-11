package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)

	if !ok {
		return false
	}

	err := util.VerifyPassword(arg.Password, e.password)

	if err != nil {
		return false
	}

	e.arg.Password = arg.Password

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{
		arg:      arg,
		password: password,
	}
}

func TestSignUpAPI(t *testing.T) {
	password := "secret1234"
	hashedPassword, _ := util.HashPassword(password)
	user := db.User{
		Username:  util.RandomUsername(),
		Email:     util.RandomMail(),
		FullName:  util.RandomUsername(),
		Password:  hashedPassword,
		CreatedAt: time.Now().Add(-time.Minute),
	}

	testCases := []struct {
		descrption    string
		requestBody   gin.H
		buildStabs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			descrption: "should save user and return its representation",
			requestBody: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  password,
			},
			buildStabs: func(store *mockdb.MockStore) {
				argGet := db.GetUserByUsernameOrEmailParams{
					Username: user.Username,
					Email:    user.Email,
				}
				argCreate := db.CreateUserParams{
					Username: user.Username,
					Email:    user.Email,
					FullName: user.FullName,
				}
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Eq(argGet)).Times(1).Return(db.User{}, sql.ErrNoRows)
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(argCreate, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				respBody := userResponse{
					Username: user.Username,
					Email:    user.Email,
					FullName: user.FullName,
				}
				require.Equal(t, recorder.Code, http.StatusCreated)
				util.RequireBodyMatchObject(t, recorder.Body, respBody)
			},
		},
		{
			descrption: "should return an error when invalid email",
			requestBody: gin.H{
				"username":  user.Username,
				"email":     "thisisnotanemail",
				"full_name": user.FullName,
				"password":  password,
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			descrption: "should return an error when too short password",
			requestBody: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  "xdxdxd2",
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			descrption: "should return an error when not alphanumeric username",
			requestBody: gin.H{
				"username":  "test#1",
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  password,
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			descrption: "should return an error when required field not included",
			requestBody: gin.H{
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  password,
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			descrption: "should return an error when too long input",
			requestBody: gin.H{
				"username":  util.RandomString(256),
				"email":     util.RandomString(256) + util.RandomMail(),
				"full_name": util.RandomString(256),
				"password":  util.RandomString(256),
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			descrption: "should return bad request when user has already existed",
			requestBody: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"full_name": user.FullName,
				"password":  password,
			},
			buildStabs: func(store *mockdb.MockStore) {
				arg := db.GetUserByUsernameOrEmailParams{
					Username: user.Username,
					Email:    user.Email,
				}
				store.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Eq(arg)).Times(1).Return(user, nil)
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
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
