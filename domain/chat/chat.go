package chat

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const messageMaxTime = 15_000

var (
	ErrMessageUpToDate = errors.New("message is up to date")
)

type Chat struct {
	id           string
	context      string
	conversation *conversation
}

type ErrChat struct {
	id string
}

func (ec ErrChat) Error() string {
	return fmt.Sprintf("error on chat with id %s", ec.id)
}

func NewChat(context string) Chat {
	return Chat{uuid.NewString(), context, &conversation{}}
}

func (c Chat) Size() int {
	return c.conversation.size
}

func (c Chat) SendMessage(msg Message) error {
	if err := msg.validateContent(); err != nil {
		return errors.Join(ErrChat{c.id}, err)
	}
	if err := c.conversation.enqueue(msg); err != nil {
		return errors.Join(ErrChat{c.id}, err)
	}
	return nil
}

func (c Chat) RemoveMessage(olderThan int64) (Message, error) {
	lastMessage, err := c.conversation.peek()
	if err != nil {
		return Message{}, errors.Join(ErrChat{c.id}, err)
	}
	if !lastMessage.isOutdated(time.Now().UnixMilli(), olderThan) {
		return Message{}, errors.Join(ErrChat{c.id}, ErrMessageUpToDate)
	}
	msg, err := c.conversation.dequeue()
	if err != nil {
		return Message{}, errors.Join(ErrChat{c.id}, err)
	}
	return msg, nil
}

func (c Chat) Conversation() ([]Message, error) {
	msgs, err := c.conversation.allMessages()
	if err != nil {
		return nil, errors.Join(ErrChat{c.id}, err)
	}
	return msgs, nil
}
