// Code generated by MockGen. DO NOT EDIT.
// Source: s3.go

// Package mock_s3 is a generated GoMock package.
package mock_s3

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockService) Delete(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockServiceMockRecorder) Delete(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService)(nil).Delete), key)
}

// Download mocks base method.
func (m *MockService) Download(key string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockServiceMockRecorder) Download(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockService)(nil).Download), key)
}

// Upload mocks base method.
func (m *MockService) Upload(key string, b []byte, expiry time.Time, contentType string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", key, b, expiry, contentType)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upload indicates an expected call of Upload.
func (mr *MockServiceMockRecorder) Upload(key, b, expiry, contentType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockService)(nil).Upload), key, b, expiry, contentType)
}