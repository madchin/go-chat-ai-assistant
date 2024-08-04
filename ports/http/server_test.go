package http_server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
	mock_chat_service "github.com/madchin/go-chat-ai-assistant/mocks/ports/chat"
	mock_history_service "github.com/madchin/go-chat-ai-assistant/mocks/ports/history"
	"go.uber.org/mock/gomock"
)

type mocks struct {
	ctrl           *gomock.Controller
	chatService    *mock_chat_service.MockChatService
	historyService *mock_history_service.MockHistoryService
}

func TestCreateChat(t *testing.T) {
	t.Run("Happy Path", func(t *testing.T) {
		mocks := setupMocks(t)
		server := &HttpServer{mocks.chatService, httprouter.New()}
		incomingReq, resWriter := httptest.NewRequest("POST", chatCreate.r, nil), httptest.NewRecorder()

		mocks.chatService.EXPECT().CreateChat(gomock.Any(), gomock.Any())
		server.createChatHandler(resWriter, incomingReq, nil)

		result := resWriter.Result()
		if result.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status code: %d, Actual: %d", http.StatusCreated, result.StatusCode)
		}
		cookie := result.Cookies()[0]
		if cookie.Name != "chatId" {
			t.Fatalf("chatId cookie not found")
		}
		if !isValidUUID(cookie.Value) {
			t.Fatalf("chatId cookie is not valid. Actual value: %s", cookie.Value)
		}
	})

	t.Run("Chat creation fails", func(t *testing.T) {
		const chatCreationError = "error"

		mocks := setupMocks(t)
		server := &HttpServer{mocks.chatService, httprouter.New()}
		incomingReq, resWriter := httptest.NewRequest("POST", chatCreate.r, nil), httptest.NewRecorder()

		mocks.chatService.EXPECT().CreateChat(gomock.Any(), gomock.Any()).Return(errors.New(chatCreationError))
		server.createChatHandler(resWriter, incomingReq, nil)

		result := resWriter.Result()
		if result.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status code: %d, Actual: %d", http.StatusBadRequest, result.StatusCode)
		}
		if len(result.Cookies()) > 0 {
			t.Fatalf("Cookies should not be set. Actual count of cookies set: %d", len(result.Cookies()))
		}
		var httpErr *HttpError
		err := json.NewDecoder(result.Body).Decode(&httpErr)
		if err != nil {
			t.Fatalf("error during decoding response body")
		}
		if httpErr.Code != serverCodeError.c {
			t.Fatalf("Expected response with Code field: %s, Actual %s", serverCodeError.c, httpErr.Code)
		}
		if httpErr.Description != chatCreationError {
			t.Fatalf("Expected response with Description field: %s, Actual %s", chatCreationError, httpErr.Description)
		}
	})
}

func TestSendMessage(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		mocks := setupMocks(t)
		server := &HttpServer{mocks.chatService, httprouter.New()}
		customerMsgContent := "my message"
		httpCustomerMsg := CustomerMessage{Content: customerMsgContent}
		marshaledCustomerMsg, _ := json.Marshal(httpCustomerMsg)
		incomingReq, resWriter := httptest.NewRequest("POST", "/chat/send", bytes.NewBuffer(marshaledCustomerMsg)), httptest.NewRecorder()
		chatId := uuid.NewString()
		cookie := chatIdCookie(chatId)
		incomingReq.AddCookie(cookie)
		incomingReq.Header.Add("Content-Type", "application/json")
		domainAssistantMsg := chat.NewAssistantMessage("Hello im here to help you")
		mocks.chatService.EXPECT().SendMessage(chatId, httpCustomerMsg.Content).Return(domainAssistantMsg, nil)

		server.sendMessageHandler(resWriter, incomingReq, nil)

		result := resWriter.Result()

		if result.StatusCode != http.StatusOK {
			t.Fatalf("expected response status %d, actual: %d", http.StatusOK, result.StatusCode)
		}
		var httpAssistantMsg AssistantMessage
		resBody, _ := io.ReadAll(result.Body)
		if err := json.Unmarshal(resBody, &httpAssistantMsg); err != nil {
			t.Fatalf("error occured during unmarshaling response body, %v", err)
		}
		if httpAssistantMsg.Author != domainAssistantMsg.Author().Role() {
			t.Fatalf("response Author field contain wrong value. expected: %s, actual: %s", chat.Assistant.Role(), httpAssistantMsg.Author)
		}
		if httpAssistantMsg.Content != domainAssistantMsg.Content() {
			t.Fatalf("response Content field contain wrong value. expected: %s, actual: %s", domainAssistantMsg.Content(), httpAssistantMsg.Content)
		}
	})
	t.Run("Header Content-Type is empty", func(t *testing.T) {
		mocks := setupMocks(t)
		server := &HttpServer{mocks.chatService, httprouter.New()}
		customerMsgContent := "my message"
		httpCustomerMsg := CustomerMessage{Content: customerMsgContent}
		marshaledCustomerMsg, _ := json.Marshal(httpCustomerMsg)
		incomingReq, resWriter := httptest.NewRequest("POST", "/chat/send", bytes.NewBuffer(marshaledCustomerMsg)), httptest.NewRecorder()

		server.sendMessageHandler(resWriter, incomingReq, nil)

		result := resWriter.Result()

		if result.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected response status %d, actual: %d", http.StatusOK, result.StatusCode)
		}
		var httpErr HttpError
		resBody, _ := io.ReadAll(result.Body)
		if err := json.Unmarshal(resBody, &httpErr); err != nil {
			t.Fatalf("error occured during unmarshaling response body, %v", err)
		}
		if httpErr.Code != clientCodeError.c {
			t.Fatalf("wrong HttpErr Code field: expected: %v, actual: %v", clientCodeError.c, httpErr.Code)
		}
		if httpErr.Description != wrongContentTypeError.d {
			t.Fatalf("wrong HttpErr Code field: expected: %v, actual: %v", wrongContentTypeError.d, httpErr.Description)
		}
	})
	t.Run("Header Content-Type != \"application/json\"", func(t *testing.T) {
		mocks := setupMocks(t)
		server := &HttpServer{mocks.chatService, httprouter.New()}
		customerMsgContent := "my message"
		httpCustomerMsg := CustomerMessage{Content: customerMsgContent}
		marshaledCustomerMsg, _ := json.Marshal(httpCustomerMsg)
		incomingReq, resWriter := httptest.NewRequest("POST", "/chat/send", bytes.NewBuffer(marshaledCustomerMsg)), httptest.NewRecorder()
		incomingReq.Header.Add("Content-Type", "text/plain")

		server.sendMessageHandler(resWriter, incomingReq, nil)

		result := resWriter.Result()

		if result.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected response status %d, actual: %d", http.StatusOK, result.StatusCode)
		}
		var httpErr HttpError
		resBody, _ := io.ReadAll(result.Body)
		if err := json.Unmarshal(resBody, &httpErr); err != nil {
			t.Fatalf("error occured during unmarshaling response body, %v", err)
		}
		if httpErr.Code != clientCodeError.c {
			t.Fatalf("wrong HttpErr Code field: expected: %v, actual: %v", clientCodeError.c, httpErr.Code)
		}
		if httpErr.Description != wrongContentTypeError.d {
			t.Fatalf("wrong HttpErr Code field: expected: %v, actual: %v", wrongContentTypeError.d, httpErr.Description)
		}
	})
	t.Run("Wrong request body", func(t *testing.T) {
		mocks := setupMocks(t)
		server := &HttpServer{mocks.chatService, httprouter.New()}
		marshaledCustomerMsg, _ := json.Marshal([]byte("{\"Some\"}"))
		incomingReq, resWriter := httptest.NewRequest("POST", "/chat/send", bytes.NewBuffer(marshaledCustomerMsg)), httptest.NewRecorder()
		incomingReq.Header.Add("Content-Type", "application/json")

		server.sendMessageHandler(resWriter, incomingReq, nil)

		result := resWriter.Result()

		if result.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected response status %d, actual: %d", http.StatusOK, result.StatusCode)
		}
		var httpErr HttpError
		resBody, _ := io.ReadAll(result.Body)
		if err := json.Unmarshal(resBody, &httpErr); err != nil {
			t.Fatalf("error occured during unmarshaling response body, %v", err)
		}
		if httpErr.Code != clientCodeError.c {
			t.Fatalf("wrong HttpErr Code field: expected: %v, actual: %v", clientCodeError.c, httpErr.Code)
		}
		if httpErr.Description != bodyParseError.d {
			t.Fatalf("wrong HttpErr Code field: expected: %v, actual: %v", bodyParseError.d, httpErr.Description)
		}
	})
	t.Run("ChatId cookie is missing", func(t *testing.T) {
		mocks := setupMocks(t)
		server := &HttpServer{mocks.chatService, httprouter.New()}
		customerMsgContent := "my message"
		httpCustomerMsg := CustomerMessage{Content: customerMsgContent}
		marshaledCustomerMsg, _ := json.Marshal(httpCustomerMsg)
		incomingReq, resWriter := httptest.NewRequest("POST", "/chat/send", bytes.NewBuffer(marshaledCustomerMsg)), httptest.NewRecorder()
		incomingReq.Header.Add("Content-Type", "application/json")

		server.sendMessageHandler(resWriter, incomingReq, nil)

		result := resWriter.Result()

		if result.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected response status %d, actual: %d", http.StatusOK, result.StatusCode)
		}
		var httpErr HttpError
		resBody, _ := io.ReadAll(result.Body)
		if err := json.Unmarshal(resBody, &httpErr); err != nil {
			t.Fatalf("error occured during unmarshaling response body, %v", err)
		}
		if httpErr.Code != clientCodeError.c {
			t.Fatalf("wrong HttpErr Code field: expected: %v, actual: %v", clientCodeError.c, httpErr.Code)
		}
		if httpErr.Description == "" {
			t.Fatalf("wrong HttpErr Description field: expected to be not empty, actual: %v", httpErr.Description)
		}
	})

	t.Run("ChatId cookie value is wrong", func(t *testing.T) {
		mocks := setupMocks(t)
		server := &HttpServer{mocks.chatService, httprouter.New()}
		customerMsgContent := "my message"
		httpCustomerMsg := CustomerMessage{Content: customerMsgContent}
		marshaledCustomerMsg, _ := json.Marshal(httpCustomerMsg)
		incomingReq, resWriter := httptest.NewRequest("POST", "/chat/send", bytes.NewBuffer(marshaledCustomerMsg)), httptest.NewRecorder()
		incomingReq.Header.Add("Content-Type", "application/json")
		chatId := "notuuid"
		cookie := chatIdCookie(chatId)
		incomingReq.AddCookie(cookie)
		server.sendMessageHandler(resWriter, incomingReq, nil)

		result := resWriter.Result()

		if result.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected response status %d, actual: %d", http.StatusOK, result.StatusCode)
		}
		var httpErr HttpError
		resBody, _ := io.ReadAll(result.Body)
		if err := json.Unmarshal(resBody, &httpErr); err != nil {
			t.Fatalf("error occured during unmarshaling response body, %v", err)
		}
		if httpErr.Code != clientCodeError.c {
			t.Fatalf("wrong HttpErr Code field: expected: %s, actual: %s", clientCodeError.c, httpErr.Code)
		}
		if httpErr.Description != wrongCookieValue.d {
			t.Fatalf("wrong HttpErr Description field: expected: %s, actual: %s", wrongCookieValue.d, httpErr.Description)
		}
	})

	t.Run("chat service SendMessage returns error", func(t *testing.T) {
		const sendMessageError = "send message error"
		mocks := setupMocks(t)
		server := &HttpServer{mocks.chatService, httprouter.New()}
		customerMsgContent := "my message"
		httpCustomerMsg := CustomerMessage{Content: customerMsgContent}
		marshaledCustomerMsg, _ := json.Marshal(httpCustomerMsg)
		incomingReq, resWriter := httptest.NewRequest("POST", "/chat/send", bytes.NewBuffer(marshaledCustomerMsg)), httptest.NewRecorder()
		incomingReq.Header.Add("Content-Type", "application/json")
		chatId := uuid.NewString()
		cookie := chatIdCookie(chatId)
		incomingReq.AddCookie(cookie)

		mocks.chatService.EXPECT().SendMessage(chatId, customerMsgContent).Return(chat.Message{}, errors.New(sendMessageError))
		server.sendMessageHandler(resWriter, incomingReq, nil)

		result := resWriter.Result()

		if result.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected response status %d, actual: %d", http.StatusOK, result.StatusCode)
		}
		var httpErr HttpError
		resBody, _ := io.ReadAll(result.Body)
		if err := json.Unmarshal(resBody, &httpErr); err != nil {
			t.Fatalf("error occured during unmarshaling response body, %v", err)
		}
		if httpErr.Code != clientCodeError.c {
			t.Fatalf("wrong HttpErr Code field: expected: %s, actual: %s", clientCodeError.c, httpErr.Code)
		}
		if httpErr.Description != sendMessageError {
			t.Fatalf("wrong HttpErr Description field: expected: %s, actual: %s", sendMessageError, httpErr.Description)
		}
	})
}

func setupMocks(t *testing.T) *mocks {
	ctrl := gomock.NewController(t)
	return &mocks{
		ctrl:           ctrl,
		chatService:    mock_chat_service.NewMockChatService(ctrl),
		historyService: mock_history_service.NewMockHistoryService(ctrl),
	}
}
