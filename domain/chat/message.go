package chat

import (
	"errors"
	"time"
)

const (
	MinCharLength = 5
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

func NewMessage(author Participant, content string, timestamp int64) Message {
	return Message{author, content, timestamp}
}

func NewCustomerMessage(content string) (Message, error) {
	msg := Message{Customer, content, time.Now().UnixMilli()}
	if err := msg.validateContent(); err != nil {
		return Message{}, err
	}
	return msg, nil
}

func NewAssistantMessage(content string) Message {
	msg := Message{Assistant, content, time.Now().UnixMilli()}
	return msg
}

func (m Message) Content() string {
	return m.content
}

func (m Message) Timestamp() int64 {
	return m.timestamp
}

func (m Message) Author() Participant {
	return m.author
}

func (m Message) isOutdated(currTime, outdateTime int64) bool {
	return currTime-m.timestamp >= outdateTime
}

func (m Message) validateContent() error {
	if len(m.content) <= MinCharLength {
		return ErrMessageShort
	}
	if len(m.content) > maxCharLength {
		return ErrMessageLong
	}
	return nil
}
