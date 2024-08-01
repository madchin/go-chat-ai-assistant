package service

import "github.com/madchin/go-chat-ai-assistant/domain/chat"

type historyService struct {
	history historyRetriever
}

type historyRetriever interface {
	RetrieveAllChatsHistory() (chat.ChatMessages, error)
}

func NewHistoryService(historyRetriever historyRetriever) historyService {
	return historyService{historyRetriever}
}

func (h *historyService) RetrieveAllChatsHistory() (chat.ChatMessages, error) {
	return h.history.RetrieveAllChatsHistory()
}
