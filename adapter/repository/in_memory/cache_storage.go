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
	initialChatsCapacity = 50
)

type UserChats map[string]*chat.Chat

type Storage struct {
	chats UserChats
	mu    sync.Mutex
	chat.Repository
}

func NewStorage() *Storage {
	return &Storage{chats: make(UserChats, initialChatsCapacity)}
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

func (s *Storage) RemoveMessage(chatId string, onRemoveMessage chat.OnRemoveMessage) (chat.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.chatExists(chatId) {
		return chat.Message{}, ErrChatNotExists
	}
	msg, err := s.chats[chatId].RemoveMessage(chat.MaxMessageTime)
	if err != nil {
		return chat.Message{}, err
	}
	onRemoveMessage(msg)
	return msg, nil
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
