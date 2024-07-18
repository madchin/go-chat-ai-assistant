package inmemory_storage

import (
	"errors"
	"sync"

	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

var (
	ErrChatAlreadyExists = errors.New("chat already exists")
	ErrChatNotExists     = errors.New("chat not exists")
)

const (
	chatsCapacity = 50
)

type UserChats map[string]*chat.Chat

type Storage struct {
	chats UserChats
	mu    sync.Mutex
	r     chat.Repository
}

func New() *Storage {
	return &Storage{chats: make(UserChats, chatsCapacity)}
}

func (s *Storage) AddNewChat(id, context string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.chatExists(id) {
		return ErrChatAlreadyExists
	}
	s.chats[id] = chat.NewChat(context)
	return nil
}

func (s *Storage) SendMessage(chatId string, msg chat.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.chatExists(chatId) {
		return ErrChatNotExists
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
	for cId, c := range s.chats {
		msgs := make([]chat.Message, 0, c.Size())
		for {
			msg, err := c.RemoveMessage(chat.MaxMessageTime)
			msgs = append(msgs, msg)
			if err == chat.ErrEmptyConversation || err == chat.ErrMessageUpToDate || err != nil {
				usersMessages[cId] = msgs
				break
			}
		}
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

func (s *Storage) chatExists(id string) bool {
	_, ok := s.chats[id]
	return ok
}
