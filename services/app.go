package service

import (
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

type Application struct {
	periodicStorageCleanUp *PeriodicCleanUp
}

func NewApplication(
	storage chat.Repository,
	history chat.HistoryRepository,
	storageService StorageService,
) *Application {
	return &Application{
		periodicStorageCleanUp: NewPeriodicCleanUpService(storageService, history),
	}
}
