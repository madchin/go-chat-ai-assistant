package http_server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	mock_chat_repository "github.com/madchin/go-chat-ai-assistant/mocks/adapters/repository/cache"
	mock_history "github.com/madchin/go-chat-ai-assistant/mocks/adapters/repository/history"
	mock_service "github.com/madchin/go-chat-ai-assistant/mocks/services"
	service "github.com/madchin/go-chat-ai-assistant/services"
	"go.uber.org/mock/gomock"
)

type mocks struct {
	ctrl             *gomock.Controller
	assistantService *mock_service.MockassistantService
	chatRepository   *mock_chat_repository.MockRepository
	chatHistory      *mock_history.MockHistoryRepository
}

func TestCreateChat(t *testing.T) {
	t.Run("Happy Path", func(t *testing.T) {
		mocks := setupMocks(t)
		server := setupHttpServer(mocks)
		incomingReq, resWriter := httptest.NewRequest("POST", chatCreate.r, nil), httptest.NewRecorder()

		mocks.chatRepository.EXPECT().CreateChat(gomock.Any(), gomock.Any())
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
		server := setupHttpServer(mocks)
		incomingReq, resWriter := httptest.NewRequest("POST", chatCreate.r, nil), httptest.NewRecorder()

		mocks.chatRepository.EXPECT().CreateChat(gomock.Any(), gomock.Any()).Return(errors.New(chatCreationError))
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

func setupMocks(t *testing.T) *mocks {
	ctrl := gomock.NewController(t)
	return &mocks{
		ctrl:             ctrl,
		assistantService: mock_service.NewMockassistantService(ctrl),
		chatRepository:   mock_chat_repository.NewMockRepository(ctrl),
		chatHistory:      mock_history.NewMockHistoryRepository(ctrl),
	}
}

func setupHttpServer(mocks *mocks) *HttpServer {
	chatService := service.NewChatService(mocks.assistantService, mocks.chatRepository)
	return &HttpServer{chatService, httprouter.New()}
}
