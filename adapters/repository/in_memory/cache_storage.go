package inmemory_storage

import (
	"errors"
	"sync"

	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

var (
	errChatAlreadyExists = errors.New("chat already exists")
	errChatNotExists     = errors.New("chat not exists")
)

const (
	chatsCapacity = 50
)

type UserChats map[string]*chat.Chat

type Storage struct {
	chats UserChats
	mu    sync.Mutex
}

func New() *Storage {
	return &Storage{chats: make(UserChats, chatsCapacity)}
}

func (s *Storage) CreateChat(id, context string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.chat(id) != nil {
		return errChatAlreadyExists
	}
	s.chats[id] = chat.NewChat(context)
	return nil
}

func (s *Storage) SendMessage(chatId string, msg chat.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.chat(chatId) == nil {
		return errChatNotExists
	}
	if err := s.chats[chatId].SendMessage(msg); err != nil {
		return err
	}
	return nil
}

func (s *Storage) RemoveOutdatedMessages() chat.UserMessages {
	s.mu.Lock()
	defer s.mu.Unlock()
	usersMessages := make(chat.UserMessages, len(s.chats))
	for chatId, c := range s.chats {
		removedMsgs := s.removeOutdatedChatMessages(c)
		usersMessages[chatId] = removedMsgs
	}
	return usersMessages
}

func (s *Storage) RetrieveAllConversations() (chat.UserMessages, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	usersMessages := make(chat.UserMessages, len(s.chats))
	for chatId, chat := range s.chats {
		msgs, _ := chat.Conversation()
		if len(msgs) > 0 {
			usersMessages[chatId] = msgs
		}
	}
	return usersMessages, nil
}

func (s *Storage) chat(id string) *chat.Chat {
	chat, ok := s.chats[id]
	if !ok {
		return nil
	}
	return chat
}

func (s *Storage) removeOutdatedChatMessages(c *chat.Chat) []chat.Message {
	msgs := make([]chat.Message, 0, c.Size())
	for {
		msg, err := c.RemoveMessage(chat.MaxMessageTime)
		if err == chat.ErrEmptyConversation || err == chat.ErrMessageUpToDate || err != nil {
			break
		}
		msgs = append(msgs, msg)
	}
	return msgs
}
