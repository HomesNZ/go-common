package mock

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockCache is a mock of Cache interface.
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance.
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockCache) Delete(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockCacheMockRecorder) Delete(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCache)(nil).Delete), key)
}

// Get mocks base method.
func (m *MockCache) Get(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCacheMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCache)(nil).Get), key)
}

// Exists mocks base method.
func (m *MockCache) Exists(key string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockCacheMockRecorder) Exists(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockCache)(nil).Exists), key)
}

// Set mocks base method.
func (m *MockCache) Set(key, val string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, val)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCacheMockRecorder) Set(key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCache)(nil).Set), key, val)
}

// SetExpiry mocks base method.
func (m *MockCache) SetExpiry(key, val string, expireTime int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetExpiry", key, val, expireTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetExpiry indicates an expected call of SetExpiry.
func (mr *MockCacheMockRecorder) SetExpiry(key, val, expireTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetExpiry", reflect.TypeOf((*MockCache)(nil).SetExpiry), key, val, expireTime)
}

// SetExpiryTime mocks base method.
func (m *MockCache) SetExpiryTime(key, val string, expireTime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetExpiryTime", key, val, expireTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetExpiryTime indicates an expected call of SetExpiryTime.
func (mr *MockCacheMockRecorder) SetExpiryTime(key, val, expireTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetExpiryTime", reflect.TypeOf((*MockCache)(nil).SetExpiryTime), key, val, expireTime)
}

// Subscribe mocks base method.
func (m *MockCache) Subscribe(subscription string, handleResponse func(interface{})) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Subscribe", subscription, handleResponse)
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockCacheMockRecorder) Subscribe(subscription, handleResponse interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockCache)(nil).Subscribe), subscription, handleResponse)
}
