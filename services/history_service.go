package service

import "github.com/madchin/go-chat-ai-assistant/domain/chat"

type HistoryRetrieveService struct {
	history historyRetriever
}

type historyRetriever interface {
	RetrieveAllChatsHistory(partialResponseCh chan<- chat.ChatMessages) error
}

func NewHistoryRetrieveService(historyRetriever historyRetriever) *HistoryRetrieveService {
	return &HistoryRetrieveService{historyRetriever}
}

func (h *HistoryRetrieveService) RetrieveAllChatsHistory(partialResponseCh chan<- chat.ChatMessages) error {
	return h.history.RetrieveAllChatsHistory(partialResponseCh)
}
