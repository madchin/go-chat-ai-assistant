package service

import (
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

type Application struct {
	ChatService                      ChatService
	periodicStorageCleanUp           *PeriodicCleanUp
	periodicClientConnectionsCleanUp *periodicClientConnectionsCleanUp
}

func NewApplication(
	storage chat.Repository,
	history chat.HistoryRepository,
	clientConnectionsStorageService clientConnectionsStorageService,
	storageService StorageService,
	assistant AssistantService,
) *Application {
	return &Application{
		ChatService:                      NewChatService(assistant, storage),
		periodicStorageCleanUp:           newPeriodicCleanUpService(storageService, history),
		periodicClientConnectionsCleanUp: newClientConnectionsCleanUpService(clientConnectionsStorageService),
	}
}
