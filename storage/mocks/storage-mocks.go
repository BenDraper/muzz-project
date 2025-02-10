// Code generated by MockGen. DO NOT EDIT.
// Source: ../storage.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	storage "muzz-project/storage"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddDecision mocks base method.
func (m *MockStorage) AddDecision(ctx context.Context, actorId, recipientId string, liked bool) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDecision", ctx, actorId, recipientId, liked)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDecision indicates an expected call of AddDecision.
func (mr *MockStorageMockRecorder) AddDecision(ctx, actorId, recipientId, liked interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDecision", reflect.TypeOf((*MockStorage)(nil).AddDecision), ctx, actorId, recipientId, liked)
}

// GetLikesCountForUser mocks base method.
func (m *MockStorage) GetLikesCountForUser(ctx context.Context, userId string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikesCountForUser", ctx, userId)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikesCountForUser indicates an expected call of GetLikesCountForUser.
func (mr *MockStorageMockRecorder) GetLikesCountForUser(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikesCountForUser", reflect.TypeOf((*MockStorage)(nil).GetLikesCountForUser), ctx, userId)
}

// GetLikesForUser mocks base method.
func (m *MockStorage) GetLikesForUser(ctx context.Context, userId string, paginationToken int) ([]*storage.Decision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikesForUser", ctx, userId, paginationToken)
	ret0, _ := ret[0].([]*storage.Decision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikesForUser indicates an expected call of GetLikesForUser.
func (mr *MockStorageMockRecorder) GetLikesForUser(ctx, userId, paginationToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikesForUser", reflect.TypeOf((*MockStorage)(nil).GetLikesForUser), ctx, userId, paginationToken)
}

// GetNewLikesForUser mocks base method.
func (m *MockStorage) GetNewLikesForUser(ctx context.Context, userId string, paginationToken int) ([]*storage.Decision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewLikesForUser", ctx, userId, paginationToken)
	ret0, _ := ret[0].([]*storage.Decision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNewLikesForUser indicates an expected call of GetNewLikesForUser.
func (mr *MockStorageMockRecorder) GetNewLikesForUser(ctx, userId, paginationToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewLikesForUser", reflect.TypeOf((*MockStorage)(nil).GetNewLikesForUser), ctx, userId, paginationToken)
}
