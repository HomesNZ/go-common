package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockConsumer is a mock of Consumer interface.
type MockConsumer struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerMockRecorder
}

// MockConsumerMockRecorder is the mock recorder for MockConsumer.
type MockConsumerMockRecorder struct {
	mock *MockConsumer
}

// NewMockConsumer creates a new mock instance.
func NewMockConsumer(ctrl *gomock.Controller) *MockConsumer {
	mock := &MockConsumer{ctrl: ctrl}
	mock.recorder = &MockConsumerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumer) EXPECT() *MockConsumerMockRecorder {
	return m.recorder
}

// BatchSize mocks base method.
func (m *MockConsumer) BatchSize(size int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchSize", size)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchSize indicates an expected call of BatchSize.
func (mr *MockConsumerMockRecorder) BatchSize(size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSize", reflect.TypeOf((*MockConsumer)(nil).BatchSize), size)
}

// Start mocks base method.
func (m *MockConsumer) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockConsumerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockConsumer)(nil).Start))
}

// Stop mocks base method.
func (m *MockConsumer) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockConsumerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockConsumer)(nil).Stop))
}

// WaitForCompletion mocks base method.
func (m *MockConsumer) WaitForCompletion(b bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WaitForCompletion", b)
}

// WaitForCompletion indicates an expected call of WaitForCompletion.
func (mr *MockConsumerMockRecorder) WaitForCompletion(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForCompletion", reflect.TypeOf((*MockConsumer)(nil).WaitForCompletion), b)
}
