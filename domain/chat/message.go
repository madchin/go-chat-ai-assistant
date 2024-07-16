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
	timestamp time.Time
}

func NewMessage(content string) Message {
	return Message{content, time.Now()}
}

func (m Message) isOutdated(t int64) bool {
	return time.Since(m.timestamp).Milliseconds() >= t
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
