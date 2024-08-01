package service

import (
	"time"

	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

const (
	cleanUpInterval = 20_000
)

type storageRemover interface {
	RemoveOutdatedMessages() chat.ChatMessages
}

type periodicCleanUp struct {
	storage  storageRemover
	history  historySaver
	interval int
}

type historySaver interface {
	SaveHistory(chat.ChatMessages) error
}

func newPeriodicCleanUpService(storage storageRemover, history historySaver) *periodicCleanUp {
	service := &periodicCleanUp{storage, history, cleanUpInterval}
	go func() {
		for {
			service.flushOutdatedMessages(history)
			time.Sleep(time.Millisecond * time.Duration(service.interval))
		}
	}()
	return service
}

func (p *periodicCleanUp) flushOutdatedMessages(history historySaver) {
	if msgs := p.storage.RemoveOutdatedMessages(); len(msgs) > 0 {
		_ = history.SaveHistory(msgs)
	}
}
