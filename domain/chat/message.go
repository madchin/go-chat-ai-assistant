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

func NewMessage(author Participant, content string) (Message, error) {
	msg := Message{author, content, time.Now().UnixMilli()}
	if msg.author == Customer {
		if err := msg.validateContent(); err != nil {
			return Message{}, err
		}
	}
	return msg, nil
}

func NewValidMessageWithTimestamp(timestamp int64) Message {
	return Message{Customer, strings.Repeat("a", minCharLength+1), timestamp}
}

func NewValidMessage() Message {
	msg, _ := NewMessage(Customer, strings.Repeat("a", minCharLength+1))
	return msg
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
