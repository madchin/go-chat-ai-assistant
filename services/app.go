package services

import (
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

type Application struct {
	storage                chat.Repository
	history                chat.HistoryRepository
	periodicStorageCleanUp *PeriodicCleanUp
}

func NewApplication(storage chat.Repository, history chat.HistoryRepository) *Application {
	return &Application{
		periodicStorageCleanUp: NewPeriodicCleanUpService(storage, history),
	}
}
