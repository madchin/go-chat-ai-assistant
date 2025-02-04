// Code generated by MockGen. DO NOT EDIT.
// Source: ./services/chat_service.go
//
// Generated by this command:
//
//	mockgen -source ./services/chat_service.go -destination ./mocks/services/assistant_service_mock.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	chat "github.com/madchin/go-chat-ai-assistant/domain/chat"
	gomock "go.uber.org/mock/gomock"
)

// MockAssistantService is a mock of AssistantService interface.
type MockAssistantService struct {
	ctrl     *gomock.Controller
	recorder *MockAssistantServiceMockRecorder
}

// MockAssistantServiceMockRecorder is the mock recorder for MockAssistantService.
type MockAssistantServiceMockRecorder struct {
	mock *MockAssistantService
}

// NewMockAssistantService creates a new mock instance.
func NewMockAssistantService(ctrl *gomock.Controller) *MockAssistantService {
	mock := &MockAssistantService{ctrl: ctrl}
	mock.recorder = &MockAssistantServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAssistantService) EXPECT() *MockAssistantServiceMockRecorder {
	return m.recorder
}

// Ask mocks base method.
func (m *MockAssistantService) Ask(content, chatId string) (chat.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ask", content, chatId)
	ret0, _ := ret[0].(chat.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Ask indicates an expected call of Ask.
func (mr *MockAssistantServiceMockRecorder) Ask(content, chatId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ask", reflect.TypeOf((*MockAssistantService)(nil).Ask), content, chatId)
}

// AskStream mocks base method.
func (m *MockAssistantService) AskStream(response chan<- string, content, chatId string) (chat.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AskStream", response, content, chatId)
	ret0, _ := ret[0].(chat.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AskStream indicates an expected call of AskStream.
func (mr *MockAssistantServiceMockRecorder) AskStream(response, content, chatId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AskStream", reflect.TypeOf((*MockAssistantService)(nil).AskStream), response, content, chatId)
}
