package ports

import (
	"io"

	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

type ChatService interface {
	CreateChat(chatId, context string) error
	SendMessage(chatId string, msg chat.Message) (chat.Message, error)
	SendMessageStream(w io.Writer, chatId string, msg chat.Message) error
}
