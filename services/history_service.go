package service

import "github.com/madchin/go-chat-ai-assistant/domain/chat"

type historyService struct {
	history historyRetriever
}

type historyRetriever interface {
	RetrieveAllChatsHistory(partialResponseCh chan<- chat.ChatMessages) error
}

func NewHistoryService(historyRetriever historyRetriever) historyService {
	return historyService{historyRetriever}
}

func (h *historyService) RetrieveAllChatsHistory(partialResponseCh chan<- chat.ChatMessages) error {
	return h.history.RetrieveAllChatsHistory(partialResponseCh)
}
