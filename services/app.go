package services

import (
	inmemory_storage "github.com/madchin/go-chat-ai-assistant/adapters/repository/in_memory"
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

type application struct {
	periodicStorageCleanUp *inmemory_storage.PeriodicCleanUp
}

func Application(storage chat.Repository, history chat.HistoryRepository) *application {
	return &application{
		periodicStorageCleanUp: inmemory_storage.NewPeriodicCleanUpService(storage, history),
	}
}
