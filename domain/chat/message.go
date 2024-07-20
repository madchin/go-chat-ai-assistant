package chat

import (
	"errors"
	"strings"
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
	author    Participant
	content   string
	timestamp int64
}

func NewMessage(author Participant, content string) Message {
	return Message{author, content, time.Now().UnixMilli()}
}

func NewValidMessageWithTimestamp(timestamp int64) Message {
	return Message{Customer, strings.Repeat("a", minCharLength+1), timestamp}
}

func NewValidMessage() Message {
	return NewMessage(Customer, strings.Repeat("a", minCharLength+1))
}

func (m Message) Content() string {
	return m.content
}

func (m Message) Author() Participant {
	return m.author
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
