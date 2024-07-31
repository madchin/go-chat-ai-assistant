package service

import (
	"time"

	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

const (
	cleanUpInterval = 20_000
)

type StorageService interface {
	RemoveOutdatedMessages() chat.ChatMessages
}

type PeriodicCleanUp struct {
	storage  StorageService
	history  History
	interval int
}

type History interface {
	SaveHistory(chat.ChatMessages)
}

func newPeriodicCleanUpService(storage StorageService, history History) *PeriodicCleanUp {
	service := &PeriodicCleanUp{storage, history, cleanUpInterval}
	go func() {
		for {
			service.flushOutdatedMessages(history)
			time.Sleep(time.Millisecond * time.Duration(service.interval))
		}
	}()
	return service
}

func (p *PeriodicCleanUp) flushOutdatedMessages(history History) {
	if msgs := p.storage.RemoveOutdatedMessages(); len(msgs) > 0 {
		history.SaveHistory(msgs)
	}
}
