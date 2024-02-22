// Code generated by MockGen. DO NOT EDIT.
// Source: ./ratelimiter/adapters/storage.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockRateLimitStorageAdapter is a mock of RateLimitStorageAdapter interface.
type MockRateLimitStorageAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockRateLimitStorageAdapterMockRecorder
}

// MockRateLimitStorageAdapterMockRecorder is the mock recorder for MockRateLimitStorageAdapter.
type MockRateLimitStorageAdapterMockRecorder struct {
	mock *MockRateLimitStorageAdapter
}

// NewMockRateLimitStorageAdapter creates a new mock instance.
func NewMockRateLimitStorageAdapter(ctrl *gomock.Controller) *MockRateLimitStorageAdapter {
	mock := &MockRateLimitStorageAdapter{ctrl: ctrl}
	mock.recorder = &MockRateLimitStorageAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateLimitStorageAdapter) EXPECT() *MockRateLimitStorageAdapterMockRecorder {
	return m.recorder
}

// AddBlock mocks base method.
func (m *MockRateLimitStorageAdapter) AddBlock(ctx context.Context, keyType, key string, milliseconds int64) (*time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBlock", ctx, keyType, key, milliseconds)
	ret0, _ := ret[0].(*time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddBlock indicates an expected call of AddBlock.
func (mr *MockRateLimitStorageAdapterMockRecorder) AddBlock(ctx, keyType, key, milliseconds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBlock", reflect.TypeOf((*MockRateLimitStorageAdapter)(nil).AddBlock), ctx, keyType, key, milliseconds)
}

// GetBlock mocks base method.
func (m *MockRateLimitStorageAdapter) GetBlock(ctx context.Context, keyType, key string) (*time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", ctx, keyType, key)
	ret0, _ := ret[0].(*time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlock indicates an expected call of GetBlock.
func (mr *MockRateLimitStorageAdapterMockRecorder) GetBlock(ctx, keyType, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockRateLimitStorageAdapter)(nil).GetBlock), ctx, keyType, key)
}

// IncrementAccesses mocks base method.
func (m *MockRateLimitStorageAdapter) IncrementAccesses(ctx context.Context, keyType, key string, maxAccesses int64) (bool, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementAccesses", ctx, keyType, key, maxAccesses)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// IncrementAccesses indicates an expected call of IncrementAccesses.
func (mr *MockRateLimitStorageAdapterMockRecorder) IncrementAccesses(ctx, keyType, key, maxAccesses interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementAccesses", reflect.TypeOf((*MockRateLimitStorageAdapter)(nil).IncrementAccesses), ctx, keyType, key, maxAccesses)
}