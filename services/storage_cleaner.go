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
	history  chat.HistorySaver
	interval int
}

func newPeriodicCleanUpService(storage storageRemover, history chat.HistorySaver) *periodicCleanUp {
	service := &periodicCleanUp{storage, history, cleanUpInterval}
	go func() {
		for {
			service.flushOutdatedMessages(history)
			time.Sleep(time.Millisecond * time.Duration(service.interval))
		}
	}()
	return service
}

func (p *periodicCleanUp) flushOutdatedMessages(history chat.HistorySaver) {
	if msgs := p.storage.RemoveOutdatedMessages(); len(msgs) > 0 {
		_ = history.SaveHistory(msgs)
	}
}
