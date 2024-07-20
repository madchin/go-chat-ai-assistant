package services

import (
	"errors"
	"time"

	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

const (
	cleanUpInterval = 20_000
)

var (
	ErrUnableToCreateMultipleServices = errors.New("unable to create multiple services")
)

type PeriodicCleanUp struct {
	storage  chat.Repository
	history  History
	interval int
}

type History interface {
	SaveHistory(chat.UserMessages) error
}

func NewPeriodicCleanUpService(storage chat.Repository, history History) *PeriodicCleanUp {
	service := &PeriodicCleanUp{storage, history, cleanUpInterval}
	go func() {
		for {
			service.FlushOutdatedMessages(history)
			time.Sleep(time.Millisecond * cleanUpInterval)
		}
	}()
	return service
}

func (p *PeriodicCleanUp) FlushOutdatedMessages(history History) {
	if msgs := p.storage.RemoveOutdatedMessages(); len(msgs) > 0 {
		history.SaveHistory(msgs)
	}
}
