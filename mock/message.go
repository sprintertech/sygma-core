// Code generated by MockGen. DO NOT EDIT.
// Source: ./relayer/message/handler.go
//
// Generated by this command:
//
//	mockgen -source=./relayer/message/handler.go -destination=./mock/message.go -package mock
//
// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	message "github.com/sygmaprotocol/sygma-core/relayer/message"
	proposal "github.com/sygmaprotocol/sygma-core/relayer/proposal"
	gomock "go.uber.org/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler[T any] struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder[T]
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder[T any] struct {
	mock *MockHandler[T]
}

// NewMockHandler creates a new mock instance.
func NewMockHandler[T any](ctrl *gomock.Controller) *MockHandler[T] {
	mock := &MockHandler[T]{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler[T]) EXPECT() *MockHandlerMockRecorder[T] {
	return m.recorder
}

// HandleMessage mocks base method.
func (m_2 *MockHandler[T]) HandleMessage(m *message.Message[T]) (*proposal.Proposal[T], error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "HandleMessage", m)
	ret0, _ := ret[0].(*proposal.Proposal[T])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleMessage indicates an expected call of HandleMessage.
func (mr *MockHandlerMockRecorder[T]) HandleMessage(m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleMessage", reflect.TypeOf((*MockHandler[T])(nil).HandleMessage), m)
}
