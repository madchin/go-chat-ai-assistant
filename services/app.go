package service

import (
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

type services struct {
	Chat    *ChatService
	History *HistoryRetrieveService
}

type cleanUpWorkers struct {
	cache                   *periodicCleanUp
	geminiClientConnections *periodicClientConnectionsCleanUp
}

type Application struct {
	Service        *services
	cleanUpWorkers *cleanUpWorkers
}

type storageService interface {
	chat.Repository
	storageRemover
}

func NewApplication(
	storage storageService,
	history chat.HistoryRepository,
	clientConnectionsStorageService connectionsStorageRemover,
	assistant AssistantService,
) *Application {
	return &Application{
		Service: &services{
			Chat:    NewChatService(assistant, storage),
			History: NewHistoryRetrieveService(history),
		},
		cleanUpWorkers: &cleanUpWorkers{
			cache:                   newPeriodicCleanUpService(storage, history),
			geminiClientConnections: newClientConnectionsCleanUpService(clientConnectionsStorageService),
		},
	}
}
