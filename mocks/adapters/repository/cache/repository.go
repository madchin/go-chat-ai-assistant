// Code generated by MockGen. DO NOT EDIT.
// Source: ./domain/chat/repository.go
//
// Generated by this command:
//
//	mockgen -source ./domain/chat/repository.go -destination ./mocks/adapters/repository/cache/repository.go -package mock_chat_repository -exclude_interfaces HistoryRepository
//

// Package mock_chat_repository is a generated GoMock package.
package mock_chat_repository

import (
	reflect "reflect"

	chat "github.com/madchin/go-chat-ai-assistant/domain/chat"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddMessage mocks base method.
func (m *MockRepository) AddMessage(chatId string, msg chat.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMessage", chatId, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMessage indicates an expected call of AddMessage.
func (mr *MockRepositoryMockRecorder) AddMessage(chatId, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMessage", reflect.TypeOf((*MockRepository)(nil).AddMessage), chatId, msg)
}

// CreateChat mocks base method.
func (m *MockRepository) CreateChat(id, context string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChat", id, context)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateChat indicates an expected call of CreateChat.
func (mr *MockRepositoryMockRecorder) CreateChat(id, context any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChat", reflect.TypeOf((*MockRepository)(nil).CreateChat), id, context)
}
