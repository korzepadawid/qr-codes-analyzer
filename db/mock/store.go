// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/korzepadawid/qr-codes-analyzer/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	sqlc "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	reflect "reflect"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateGroup mocks base method
func (m *MockStore) CreateGroup(arg0 context.Context, arg1 sqlc.CreateGroupParams) (sqlc.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup
func (mr *MockStoreMockRecorder) CreateGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockStore)(nil).CreateGroup), arg0, arg1)
}

// CreateQRCode mocks base method
func (m *MockStore) CreateQRCode(arg0 context.Context, arg1 sqlc.CreateQRCodeParams) (sqlc.QrCode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateQRCode", arg0, arg1)
	ret0, _ := ret[0].(sqlc.QrCode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateQRCode indicates an expected call of CreateQRCode
func (mr *MockStoreMockRecorder) CreateQRCode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQRCode", reflect.TypeOf((*MockStore)(nil).CreateQRCode), arg0, arg1)
}

// CreateRedirect mocks base method
func (m *MockStore) CreateRedirect(arg0 context.Context, arg1 sqlc.CreateRedirectParams) (sqlc.Redirect, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRedirect", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Redirect)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRedirect indicates an expected call of CreateRedirect
func (mr *MockStoreMockRecorder) CreateRedirect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRedirect", reflect.TypeOf((*MockStore)(nil).CreateRedirect), arg0, arg1)
}

// CreateUser mocks base method
func (m *MockStore) CreateUser(arg0 context.Context, arg1 sqlc.CreateUserParams) (sqlc.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(sqlc.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// DeleteGroupByOwnerAndID mocks base method
func (m *MockStore) DeleteGroupByOwnerAndID(arg0 context.Context, arg1 sqlc.DeleteGroupByOwnerAndIDParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroupByOwnerAndID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroupByOwnerAndID indicates an expected call of DeleteGroupByOwnerAndID
func (mr *MockStoreMockRecorder) DeleteGroupByOwnerAndID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroupByOwnerAndID", reflect.TypeOf((*MockStore)(nil).DeleteGroupByOwnerAndID), arg0, arg1)
}

// GetGroupByOwnerAndID mocks base method
func (m *MockStore) GetGroupByOwnerAndID(arg0 context.Context, arg1 sqlc.GetGroupByOwnerAndIDParams) (sqlc.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupByOwnerAndID", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupByOwnerAndID indicates an expected call of GetGroupByOwnerAndID
func (mr *MockStoreMockRecorder) GetGroupByOwnerAndID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupByOwnerAndID", reflect.TypeOf((*MockStore)(nil).GetGroupByOwnerAndID), arg0, arg1)
}

// GetGroupByOwnerAndIDForUpdate mocks base method
func (m *MockStore) GetGroupByOwnerAndIDForUpdate(arg0 context.Context, arg1 sqlc.GetGroupByOwnerAndIDForUpdateParams) (sqlc.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupByOwnerAndIDForUpdate", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupByOwnerAndIDForUpdate indicates an expected call of GetGroupByOwnerAndIDForUpdate
func (mr *MockStoreMockRecorder) GetGroupByOwnerAndIDForUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupByOwnerAndIDForUpdate", reflect.TypeOf((*MockStore)(nil).GetGroupByOwnerAndIDForUpdate), arg0, arg1)
}

// GetGroupsByOwner mocks base method
func (m *MockStore) GetGroupsByOwner(arg0 context.Context, arg1 sqlc.GetGroupsByOwnerParams) ([]sqlc.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupsByOwner", arg0, arg1)
	ret0, _ := ret[0].([]sqlc.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupsByOwner indicates an expected call of GetGroupsByOwner
func (mr *MockStoreMockRecorder) GetGroupsByOwner(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupsByOwner", reflect.TypeOf((*MockStore)(nil).GetGroupsByOwner), arg0, arg1)
}

// GetGroupsCountByOwner mocks base method
func (m *MockStore) GetGroupsCountByOwner(arg0 context.Context, arg1 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupsCountByOwner", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupsCountByOwner indicates an expected call of GetGroupsCountByOwner
func (mr *MockStoreMockRecorder) GetGroupsCountByOwner(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupsCountByOwner", reflect.TypeOf((*MockStore)(nil).GetGroupsCountByOwner), arg0, arg1)
}

// GetQRCode mocks base method
func (m *MockStore) GetQRCode(arg0 context.Context, arg1 string) (sqlc.QrCode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQRCode", arg0, arg1)
	ret0, _ := ret[0].(sqlc.QrCode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQRCode indicates an expected call of GetQRCode
func (mr *MockStoreMockRecorder) GetQRCode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQRCode", reflect.TypeOf((*MockStore)(nil).GetQRCode), arg0, arg1)
}

// GetQRCodeForUpdate mocks base method
func (m *MockStore) GetQRCodeForUpdate(arg0 context.Context, arg1 string) (sqlc.QrCode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQRCodeForUpdate", arg0, arg1)
	ret0, _ := ret[0].(sqlc.QrCode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQRCodeForUpdate indicates an expected call of GetQRCodeForUpdate
func (mr *MockStoreMockRecorder) GetQRCodeForUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQRCodeForUpdate", reflect.TypeOf((*MockStore)(nil).GetQRCodeForUpdate), arg0, arg1)
}

// GetUserByEmail mocks base method
func (m *MockStore) GetUserByEmail(arg0 context.Context, arg1 string) (sqlc.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(sqlc.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail
func (mr *MockStoreMockRecorder) GetUserByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockStore)(nil).GetUserByEmail), arg0, arg1)
}

// GetUserByUsername mocks base method
func (m *MockStore) GetUserByUsername(arg0 context.Context, arg1 string) (sqlc.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", arg0, arg1)
	ret0, _ := ret[0].(sqlc.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername
func (mr *MockStoreMockRecorder) GetUserByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockStore)(nil).GetUserByUsername), arg0, arg1)
}

// GetUserByUsernameOrEmail mocks base method
func (m *MockStore) GetUserByUsernameOrEmail(arg0 context.Context, arg1 sqlc.GetUserByUsernameOrEmailParams) (sqlc.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsernameOrEmail", arg0, arg1)
	ret0, _ := ret[0].(sqlc.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsernameOrEmail indicates an expected call of GetUserByUsernameOrEmail
func (mr *MockStoreMockRecorder) GetUserByUsernameOrEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsernameOrEmail", reflect.TypeOf((*MockStore)(nil).GetUserByUsernameOrEmail), arg0, arg1)
}

// IncrementQRCodeEntries mocks base method
func (m *MockStore) IncrementQRCodeEntries(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementQRCodeEntries", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementQRCodeEntries indicates an expected call of IncrementQRCodeEntries
func (mr *MockStoreMockRecorder) IncrementQRCodeEntries(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementQRCodeEntries", reflect.TypeOf((*MockStore)(nil).IncrementQRCodeEntries), arg0, arg1)
}

// IncrementRedirectEntriesTx mocks base method
func (m *MockStore) IncrementRedirectEntriesTx(arg0 context.Context, arg1 sqlc.IncrementRedirectEntriesTxParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementRedirectEntriesTx", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementRedirectEntriesTx indicates an expected call of IncrementRedirectEntriesTx
func (mr *MockStoreMockRecorder) IncrementRedirectEntriesTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementRedirectEntriesTx", reflect.TypeOf((*MockStore)(nil).IncrementRedirectEntriesTx), arg0, arg1)
}

// UpdateGroupByOwnerAndID mocks base method
func (m *MockStore) UpdateGroupByOwnerAndID(arg0 context.Context, arg1 sqlc.UpdateGroupByOwnerAndIDParams) (sqlc.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroupByOwnerAndID", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateGroupByOwnerAndID indicates an expected call of UpdateGroupByOwnerAndID
func (mr *MockStoreMockRecorder) UpdateGroupByOwnerAndID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroupByOwnerAndID", reflect.TypeOf((*MockStore)(nil).UpdateGroupByOwnerAndID), arg0, arg1)
}

// UpdateGroupTx mocks base method
func (m *MockStore) UpdateGroupTx(arg0 context.Context, arg1 sqlc.UpdateGroupTxParams) (sqlc.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroupTx", arg0, arg1)
	ret0, _ := ret[0].(sqlc.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateGroupTx indicates an expected call of UpdateGroupTx
func (mr *MockStoreMockRecorder) UpdateGroupTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroupTx", reflect.TypeOf((*MockStore)(nil).UpdateGroupTx), arg0, arg1)
}
