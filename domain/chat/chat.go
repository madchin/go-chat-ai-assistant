package chat

import (
	"errors"
	"time"
)

const MaxMessageTime = 15_000

var (
	ErrMessageUpToDate = errors.New("message is up to date")
)

type Chat struct {
	context      string
	conversation *conversation
}

func NewChat(context string) *Chat {
	return &Chat{context, &conversation{}}
}

func (c *Chat) Size() int {
	return c.conversation.size
}

func (c *Chat) SendMessage(msg Message) error {
	if err := msg.validateContent(); err != nil {
		return err
	}
	if err := c.conversation.enqueue(msg); err != nil {
		return err
	}
	return nil
}

func (c *Chat) RemoveMessage(olderThan int64) (Message, error) {
	lastMessage, err := c.conversation.peek()
	if err != nil {
		return Message{}, err
	}
	if !lastMessage.isOutdated(time.Now().UnixMilli(), olderThan) {
		return Message{}, ErrMessageUpToDate
	}
	msg, err := c.conversation.dequeue()
	if err != nil {
		return Message{}, err
	}
	return msg, nil
}

func (c *Chat) Conversation() ([]Message, error) {
	msgs, err := c.conversation.allMessages()
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
