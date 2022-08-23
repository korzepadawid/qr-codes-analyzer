package qr_code

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	mockcache "github.com/korzepadawid/qr-codes-analyzer/cache/mock"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	mockdb "github.com/korzepadawid/qr-codes-analyzer/db/mock"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	mockencoder "github.com/korzepadawid/qr-codes-analyzer/encode/mock"
	mockstorage "github.com/korzepadawid/qr-codes-analyzer/storage/mock"
	mocktoken "github.com/korzepadawid/qr-codes-analyzer/token/mock"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type createQRCodeStabs struct {
	store         *mockdb.MockStore
	tokenProvider *mocktoken.MockProvider
	storage       *mockstorage.MockFileStorage
	cache         *mockcache.MockCache
	encoder       *mockencoder.MockEncoder
}

var wg sync.WaitGroup

func TestCreateQRCodeAPI(t *testing.T) {
	testCases := []struct {
		description   string
		requestBody   createQRCodeRequest
		buildStabs    func(stabs createQRCodeStabs)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			description: "should create new QRCode and cache it when ok",
			requestBody: mockCreateQRCodeRequestBody,
			buildStabs: func(stabs createQRCodeStabs) {
				stabs.tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				stabs.store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), gomock.Eq(db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				})).Times(1).Return(mockGroup, nil)
				stabs.encoder.EXPECT().Encode(gomock.Any()).Times(1).Return(mockGeneratedQRCode, nil)
				stabs.storage.EXPECT().PutFile(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				stabs.store.EXPECT().CreateQRCode(gomock.Any(), gomock.Any()).Times(1).Return(mockSavedQRCode, nil)
				stabs.cache.EXPECT().Set(gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			description: "should return an error when group does not exists",
			requestBody: mockCreateQRCodeRequestBody,
			buildStabs: func(stabs createQRCodeStabs) {
				stabs.tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				stabs.store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), gomock.Eq(db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				})).Times(1).Return(db.Group{}, sql.ErrNoRows)
				stabs.encoder.EXPECT().Encode(gomock.Any()).Times(0)
				stabs.storage.EXPECT().PutFile(gomock.Any(), gomock.Any()).Times(0)
				stabs.store.EXPECT().CreateQRCode(gomock.Any(), gomock.Any()).Times(0)
				stabs.cache.EXPECT().Set(gomock.Any()).Times(0.)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			description: "should return an error when fails during group fetching",
			requestBody: mockCreateQRCodeRequestBody,
			buildStabs: func(stabs createQRCodeStabs) {
				stabs.tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				stabs.store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), gomock.Eq(db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				})).Times(1).Return(db.Group{}, sql.ErrConnDone)
				stabs.encoder.EXPECT().Encode(gomock.Any()).Times(0)
				stabs.storage.EXPECT().PutFile(gomock.Any(), gomock.Any()).Times(0)
				stabs.store.EXPECT().CreateQRCode(gomock.Any(), gomock.Any()).Times(0)
				stabs.cache.EXPECT().Set(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "should return an error when db qr code generation fails",
			requestBody: mockCreateQRCodeRequestBody,
			buildStabs: func(stabs createQRCodeStabs) {
				stabs.tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				stabs.store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), gomock.Eq(db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				})).Times(1).Return(mockGroup, nil)
				stabs.encoder.EXPECT().Encode(gomock.Any()).Times(1).Return(make([]byte, 0), errors.New("fail"))
				stabs.storage.EXPECT().PutFile(gomock.Any(), gomock.Any()).Times(0)
				stabs.store.EXPECT().CreateQRCode(gomock.Any(), gomock.Any()).Times(0)
				stabs.cache.EXPECT().Set(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "should return an error when storage fails",
			requestBody: mockCreateQRCodeRequestBody,
			buildStabs: func(stabs createQRCodeStabs) {
				stabs.tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				stabs.store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), gomock.Eq(db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				})).Times(1).Return(mockGroup, nil)
				stabs.encoder.EXPECT().Encode(gomock.Any()).Times(1).Return(mockGeneratedQRCode, nil)
				stabs.storage.EXPECT().PutFile(gomock.Any(), gomock.Any()).Times(1).Return(errors.New("fail"))
				stabs.store.EXPECT().CreateQRCode(gomock.Any(), gomock.Any()).Times(0)
				stabs.cache.EXPECT().Set(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "should return an error when db fails during insertion and should delete file from the storage",
			requestBody: mockCreateQRCodeRequestBody,
			buildStabs: func(stabs createQRCodeStabs) {
				stabs.tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				stabs.store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), gomock.Eq(db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				})).Times(1).Return(mockGroup, nil)
				stabs.encoder.EXPECT().Encode(gomock.Any()).Times(1).Return(mockGeneratedQRCode, nil)
				stabs.storage.EXPECT().PutFile(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				stabs.store.EXPECT().CreateQRCode(gomock.Any(), gomock.Any()).Times(1).Return(db.QrCode{}, sql.ErrConnDone)
				stabs.storage.EXPECT().DeleteFile(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				stabs.cache.EXPECT().Set(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			description: "should return an error when url is not valid",
			requestBody: createQRCodeRequest{
				URL:   "root@localhost.com",
				Title: util.RandomString(255),
			},
			buildStabs: func(stabs createQRCodeStabs) {
				stabs.tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				stabs.store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), gomock.Eq(db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				})).Times(0)
				stabs.encoder.EXPECT().Encode(gomock.Any()).Times(0)
				stabs.storage.EXPECT().PutFile(gomock.Any(), gomock.Any()).Times(0)
				stabs.store.EXPECT().CreateQRCode(gomock.Any(), gomock.Any()).Times(0)
				stabs.cache.EXPECT().Set(gomock.Any()).Times(0.)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description: "should return an error when title is empty",
			requestBody: createQRCodeRequest{
				URL:   mockCreateQRCodeRequestBody.URL,
				Title: "",
			},
			buildStabs: func(stabs createQRCodeStabs) {
				stabs.tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				stabs.store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), gomock.Eq(db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				})).Times(0)
				stabs.encoder.EXPECT().Encode(gomock.Any()).Times(0)
				stabs.storage.EXPECT().PutFile(gomock.Any(), gomock.Any()).Times(0)
				stabs.store.EXPECT().CreateQRCode(gomock.Any(), gomock.Any()).Times(0)
				stabs.cache.EXPECT().Set(gomock.Any()).Times(0.)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			description: "should return an error when title is too long",
			requestBody: createQRCodeRequest{
				URL:         mockCreateQRCodeRequestBody.URL,
				Title:       util.RandomString(256),
				Description: util.RandomString(2),
			},
			buildStabs: func(stabs createQRCodeStabs) {
				stabs.tokenProvider.EXPECT().VerifyToken(gomock.Any()).Times(1).Return(mockPayload, nil)
				stabs.store.EXPECT().GetGroupByOwnerAndID(gomock.Any(), gomock.Eq(db.GetGroupByOwnerAndIDParams{
					Owner:   mockPayload.Username,
					GroupID: randomGroupID,
				})).Times(0)
				stabs.encoder.EXPECT().Encode(gomock.Any()).Times(0)
				stabs.storage.EXPECT().PutFile(gomock.Any(), gomock.Any()).Times(0)
				stabs.store.EXPECT().CreateQRCode(gomock.Any(), gomock.Any()).Times(0)
				stabs.cache.EXPECT().Set(gomock.Any()).Times(0.)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			// mocks
			cfg := config.Config{AppURL: util.RandomString(2), CDNAddress: util.RandomString(2)}
			ctrl := gomock.NewController(t)
			mockTokenProvider := mocktoken.NewMockProvider(ctrl)
			mockStore := mockdb.NewMockStore(ctrl)
			mockCache := mockcache.NewMockCache(ctrl)
			mockEncoder := mockencoder.NewMockEncoder(ctrl)
			mockFileStorage := mockstorage.NewMockFileStorage(ctrl)

			// stabs
			r := newMockQRCodeHandler(mockStore, cfg, mockFileStorage, mockEncoder, mockCache, mockTokenProvider)
			tC.buildStabs(createQRCodeStabs{
				store:         mockStore,
				tokenProvider: mockTokenProvider,
				storage:       mockFileStorage,
				cache:         mockCache,
				encoder:       mockEncoder,
			})
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/groups/%d/qr-codes", randomGroupID), util.MarshallBody(tC.requestBody))
			request.Header.Set("Authorization", "Bearer faketokenishere")
			require.NoError(t, err)
			r.ServeHTTP(recorder, request)
			// check response
			tC.checkResponse(t, recorder)
		})
	}
}
