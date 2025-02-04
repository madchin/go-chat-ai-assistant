package chat

import (
	"errors"
)

const MaxConversationSize = 5 * 2

var (
	ErrEmptyConversation           = errors.New("conversation is empty")
	ErrConversationMaxSizeExceeded = errors.New("conversation size is exceeded")
)

type node struct {
	message Message
	next    *node
}

type conversation struct {
	size  int
	first *node
	last  *node
}

func (m *conversation) enqueue(message Message) error {
	if m.size == MaxConversationSize {
		return ErrConversationMaxSizeExceeded
	}
	if m.size == 0 {
		m.first = &node{message, nil}
		m.last = m.first
	} else {
		m.last.next = &node{message, nil}
		m.last = m.last.next
	}
	m.size++
	return nil
}

func (m *conversation) dequeue() (Message, error) {
	if m.size == 0 {
		return Message{}, ErrEmptyConversation
	}
	m.size--
	msg := m.first.message
	m.first = m.first.next
	if m.size == 0 {
		m.last = nil
	}
	return msg, nil
}

func (m *conversation) peek() (Message, error) {
	if m.size == 0 {
		return Message{}, ErrEmptyConversation
	}
	return m.first.message, nil
}
