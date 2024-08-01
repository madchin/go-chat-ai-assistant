package cache

import (
	"errors"
	"sync"

	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

var (
	ErrChatAlreadyExists = errors.New("chat already exists")
	errChatNotExists     = errors.New("chat not exists")
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
	if s.chat(chatId) == nil {
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
	usersMessages := make(chat.ChatMessages, len(s.chats))
	for chatId, c := range s.chats {
		removedMsgs := s.removeOutdatedChatMessages(c)
		usersMessages[chatId] = removedMsgs
		if c.IsConversationEmpty() {
			delete(s.chats, chatId)
		}
	}
	return usersMessages
}

func (s *storage) RetrieveAllConversations() (chat.ChatMessages, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	usersMessages := make(chat.ChatMessages, len(s.chats))
	for chatId, chat := range s.chats {
		msgs, _ := chat.Conversation()
		if len(msgs) > 0 {
			usersMessages[chatId] = msgs
		}
	}
	return usersMessages, nil
}

func (s *storage) createChat(id, context string) error {
	if s.chat(id) != nil {
		return ErrChatAlreadyExists
	}
	s.chats[id] = chat.NewChat(context)
	return nil
}

func (s *storage) chat(id string) *chat.Chat {
	chat, ok := s.chats[id]
	if !ok {
		return nil
	}
	return chat
}

func (s *storage) removeOutdatedChatMessages(c *chat.Chat) []chat.Message {
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
