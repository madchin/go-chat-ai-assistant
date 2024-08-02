package ports

import (
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

type ChatService interface {
	CreateChat(chatId, context string) error
	SendMessage(chatId string, customerMsg chat.Message) (chat.Message, error)
	SendMessageStream(responseCh chan<- string, chatId string, customerMsg chat.Message) error
}

type HistoryService interface {
	RetrieveAllChatsHistory(partialResponseCh chan<- chat.ChatMessages) error
}
