package chat

import (
	"errors"
	"time"
)

const (
	minCharLength = 5
	maxCharLength = 200
)

var (
	ErrMessageShort = errors.New("message is too short")
	ErrMessageLong  = errors.New("message is too long")
)

type Message struct {
	content   string
	timestamp int64
}

func NewMessage(content string) Message {
	return Message{content, time.Now().UnixMilli()}
}

func (m Message) isOutdated(currTime, outdateTime int64) bool {
	return currTime-m.timestamp >= outdateTime
}

func (m Message) validateContent() error {
	if len(m.content) <= minCharLength {
		return ErrMessageShort
	}
	if len(m.content) > maxCharLength {
		return ErrMessageLong
	}
	return nil
}
