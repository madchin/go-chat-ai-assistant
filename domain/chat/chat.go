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

type chat struct {
	id           string
	context      string
	conversation *conversation
}

func NewChat(context string) chat {
	return chat{uuid.NewString(), context, &conversation{}}
}

func (c chat) Size() int {
	return c.conversation.size
}

func (c chat) SendMessage(msg Message) error {
	if err := msg.validateContent(); err != nil {
		return errorChat(c, err)
	}
	if err := c.conversation.enqueue(msg); err != nil {
		return errorChat(c, err)
	}
	return nil
}

func (c chat) RemoveMessage() (Message, error) {
	lastMessage, err := c.conversation.peek()
	if err != nil {
		return Message{}, errorChat(c, err)
	}
	if !lastMessage.isOutdated(time.Now().UnixMilli(), messageMaxTime) {
		return Message{}, errorChat(c, ErrMessageUpToDate)
	}
	msg, err := c.conversation.dequeue()
	if err != nil {
		return Message{}, errorChat(c, err)
	}
	return msg, errorChat(c, err)
}

func (c chat) Conversation() ([]Message, error) {
	msgs, err := c.conversation.allMessages()
	if err != nil {
		return nil, errorChat(c, err)
	}
	return msgs, nil
}

func errorChat(chat chat, err error) error {
	return fmt.Errorf("chat with id %s: %w", chat.id, err)
}
