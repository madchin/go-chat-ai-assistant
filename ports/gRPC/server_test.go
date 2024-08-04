package grpc_server

import (
	"context"
	"errors"
	"testing"

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
	t.Run("Happy path", func(t *testing.T) {
		mocks := setupMocks(t)
		grpcServer := GrpcServer{mocks.chatService, mocks.historyService}
		mocks.chatService.EXPECT().CreateChat(gomock.Any(), gomock.Any())
		resp, err := grpcServer.CreateChat(context.Background(), &ChatRequest{})
		if err != nil {
			t.Fatalf("error occured during creating chat, %v", err)
		}
		if resp.ChatId == "" {
			t.Fatalf("chat id has not been returned")
		}
	})

	t.Run("Chat service creating chat error", func(t *testing.T) {
		const createChatError = "error"
		mocks := setupMocks(t)
		grpcServer := GrpcServer{mocks.chatService, mocks.historyService}
		mocks.chatService.EXPECT().CreateChat(gomock.Any(), gomock.Any()).Return(errors.New(createChatError))
		_, err := grpcServer.CreateChat(context.Background(), &ChatRequest{})
		if err == nil {
			t.Fatalf("error should occur during creating chat")
		}
	})
}

func TestSendMessage(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		mocks := setupMocks(t)
		grpcServer := GrpcServer{mocks.chatService, mocks.historyService}
		chatId, content := "chatId", "Content"
		mocks.chatService.EXPECT().SendMessage(chatId, content).Return(chat.NewAssistantMessage(content), nil)
		resp, err := grpcServer.SendMessage(context.Background(), &MessageRequest{ChatId: chatId, Content: content})
		if err != nil {
			t.Fatalf("error occured during sending message, %v", err)
		}
		if resp.GetAuthor() != chat.Assistant.Role() {
			t.Fatalf("response Author field is wrong. Expected: %s, Actual: %s", chat.Assistant.Role(), resp.GetAuthor())
		}
		if resp.GetContent() != content {
			t.Fatalf("response Content field is wrong. Expected %s, Actual: %s", content, resp.GetContent())
		}
	})
	t.Run("Error in sending message", func(t *testing.T) {
		const sendMessageError = "send message error"
		mocks := setupMocks(t)
		chatId, content := "chatId", "Content"
		grpcServer := GrpcServer{mocks.chatService, mocks.historyService}
		mocks.chatService.EXPECT().SendMessage(chatId, content).Return(chat.Message{}, errors.New(sendMessageError))
		_, err := grpcServer.SendMessage(context.Background(), &MessageRequest{ChatId: chatId, Content: content})
		if err == nil {
			t.Fatalf("error has not occured during sending message")
		}
	})
}

func TestSendMessageStream(t *testing.T) {
	t.Run("Happy Path", func(t *testing.T) {
		mocks := setupMocks(t)
		chatId, content := "chatId", "Content"
		grpcServer := GrpcServer{mocks.chatService, mocks.historyService}
		mocks.chatService.
			EXPECT().
			SendMessageStream(gomock.Any(), chatId, content).
			DoAndReturn(func(responseCh chan<- string, chatId, content string) error {
				close(responseCh)
				return nil
			})
		err := grpcServer.SendMessageStream(&MessageRequest{ChatId: chatId, Content: content}, &chatSendMessageStreamServer{})
		if err != nil {
			t.Fatalf("Error occured during sending message stream, err: %v", err)
		}
	})
	t.Run("error on sending message stream", func(t *testing.T) {
		mocks := setupMocks(t)
		chatId, content := "chatId", "Content"
		grpcServer := GrpcServer{mocks.chatService, mocks.historyService}
		mocks.chatService.EXPECT().SendMessageStream(gomock.Any(), chatId, content).Return(errors.New("error"))
		err := grpcServer.SendMessageStream(&MessageRequest{ChatId: chatId, Content: content}, &chatSendMessageStreamServer{})
		if err == nil {
			t.Fatalf("Error should occur during sending message stream")
		}
	})
}

func TestRetrieveHistory(t *testing.T) {
	t.Run("Happy Path", func(t *testing.T) {
		mocks := setupMocks(t)
		grpcServer := GrpcServer{mocks.chatService, mocks.historyService}
		mocks.historyService.
			EXPECT().
			RetrieveAllChatsHistoryStream(gomock.Any()).DoAndReturn(func(partialResponseCh chan<- chat.ChatMessages) error {
			close(partialResponseCh)
			return nil
		})
		err := grpcServer.RetrieveHistory(nil, &chatRetrieveHistoryServer{})
		if err != nil {
			t.Fatalf("Error occured during sending message stream, err: %v", err)
		}
	})
	t.Run("error on sending message stream", func(t *testing.T) {
		const sendingMessageError = "send msesage error"
		mocks := setupMocks(t)
		grpcServer := GrpcServer{mocks.chatService, mocks.historyService}
		mocks.historyService.
			EXPECT().
			RetrieveAllChatsHistoryStream(gomock.Any()).DoAndReturn(func(partialResponseCh chan<- chat.ChatMessages) error {
			close(partialResponseCh)
			return errors.New(sendingMessageError)
		})
		err := grpcServer.RetrieveHistory(nil, &chatRetrieveHistoryServer{})
		if err == nil {
			t.Fatalf("Error should occur during sending message stream")
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
