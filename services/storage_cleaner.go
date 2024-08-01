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
	history  historySaver
	interval int
}

type historySaver interface {
	SaveHistory(chat.ChatMessages) error
}

func newPeriodicCleanUpService(storage StorageService, history historySaver) *PeriodicCleanUp {
	service := &PeriodicCleanUp{storage, history, cleanUpInterval}
	go func() {
		for {
			service.flushOutdatedMessages(history)
			time.Sleep(time.Millisecond * time.Duration(service.interval))
		}
	}()
	return service
}

func (p *PeriodicCleanUp) flushOutdatedMessages(history historySaver) {
	if msgs := p.storage.RemoveOutdatedMessages(); len(msgs) > 0 {
		_ = history.SaveHistory(msgs)
	}
}
