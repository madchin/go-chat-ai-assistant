package service

import "github.com/madchin/go-chat-ai-assistant/domain/chat"

type HistoryRetrieveService struct {
	history chat.HistoryRetriever
}

func NewHistoryRetrieveService(historyRetriever chat.HistoryRetriever) *HistoryRetrieveService {
	return &HistoryRetrieveService{historyRetriever}
}

func (h *HistoryRetrieveService) RetrieveAllChatsHistoryStream(partialResponseCh chan<- chat.ChatMessages) error {
	return h.history.RetrieveAllChatsHistoryStream(partialResponseCh)
}

func (h *HistoryRetrieveService) RetrieveAllChatsHistory() (chat.ChatMessages, error) {
	return h.history.RetrieveAllChatsHistory()
}
