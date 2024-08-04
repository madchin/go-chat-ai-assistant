package http_server

import "github.com/madchin/go-chat-ai-assistant/domain/chat"

type CustomerMessage struct {
	Content string `json:"content"`
}

type AssistantMessage struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type HttpError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Message struct {
	Author    string `json:"author"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type ChatsHistory struct {
	Chats map[string][]Message `json:"chats"`
}

type codeErrors struct {
	c string
}

type descriptionErrors struct {
	d string
}

var (
	serverCodeError = codeErrors{"server"}
	clientCodeError = codeErrors{"client"}

	genericDescriptionError = descriptionErrors{"Oops! Something went wrong"}
	wrongContentTypeError   = descriptionErrors{"Content-Type header need to be application/json"}
	bodyParseError          = descriptionErrors{"JSON parse failed"}
	wrongCookieValue        = descriptionErrors{"Wrong cookie value"}
)

func mapDomainMessageToHttpAssistantMessage(domainMsg chat.Message) AssistantMessage {
	return AssistantMessage{Author: domainMsg.Author().Role(), Content: domainMsg.Content()}
}

func mapDomainChatMessagesToHttpChatsHistory(domainChatMessages chat.ChatMessages) ChatsHistory {
	chatsHistory := make(map[string][]Message, len(domainChatMessages))
	for chatId, domainMessages := range domainChatMessages {
		httpMessages := make([]Message, len(domainMessages))
		for i := 0; i < len(domainMessages); i++ {
			httpMessages[i] = Message{Author: domainMessages[i].Author().Role(), Content: domainMessages[i].Content(), Timestamp: domainMessages[i].Timestamp()}
		}
		chatsHistory[chatId] = httpMessages
	}
	return ChatsHistory{Chats: chatsHistory}
}
