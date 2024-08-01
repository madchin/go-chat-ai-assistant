package cache

import (
	"errors"
	"sync"

	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

var (
	errChatAlreadyExists = errors.New("chat already exists")
)

const (
	ChatsCapacity = 50
)

type userChats map[string]*chat.Chat

type storage struct {
	chats userChats
	mu    sync.Mutex
}

func New() *storage {
	return &storage{chats: make(userChats, ChatsCapacity)}
}

func (s *storage) CreateChat(id, context string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.createChat(id, context)
}

func (s *storage) SendMessage(chatId string, msg chat.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.chatExists(chatId) {
		_ = s.createChat(chatId, "")
	}
	if err := s.chats[chatId].SendMessage(msg); err != nil {
		return err
	}
	return nil
}

func (s *storage) RemoveOutdatedMessages() chat.ChatMessages {
	s.mu.Lock()
	defer s.mu.Unlock()
	removedMessages := make(chat.ChatMessages, len(s.chats))
	for chatId, chat := range s.chats {
		removedMessages[chatId] = s.removeOutdatedChatMessages(chat)
		if chat.IsConversationEmpty() {
			delete(s.chats, chatId)
		}
	}
	return removedMessages
}

func (s *storage) createChat(id, context string) error {
	if s.chatExists(id) {
		return errChatAlreadyExists
	}
	s.chats[id] = chat.NewChat(context)
	return nil
}

func (s *storage) chatExists(id string) bool {
	_, ok := s.chats[id]
	return ok
}

func (s *storage) removeOutdatedChatMessages(c *chat.Chat) []chat.Message {
	msgs := make([]chat.Message, 0, c.Size())
	for {
		msg, err := c.RemoveMessage(chat.MaxMessageTime)
		if err != nil {
			break
		}
		msgs = append(msgs, msg)
	}
	return msgs
}
