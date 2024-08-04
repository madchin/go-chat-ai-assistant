package service

import "github.com/madchin/go-chat-ai-assistant/domain/chat"

type HistoryRetrieveService struct {
	history chat.HistoryRetriever
}

func NewHistoryRetrieveService(historyRetriever chat.HistoryRetriever) *HistoryRetrieveService {
	return &HistoryRetrieveService{historyRetriever}
}

func (h *HistoryRetrieveService) RetrieveAllChatsHistory(partialResponseCh chan<- chat.ChatMessages) error {
	return h.history.RetrieveAllChatsHistory(partialResponseCh)
}
