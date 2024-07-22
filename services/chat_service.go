package service

import (
	"bytes"
	"io"

	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

const maxResponseLength = 1024

type AssistantService interface {
	SendMessage(w io.Writer, content, chatId string) error
	SendMessageStream(w io.Writer, content, chatId string) error
}

type ChatService struct {
	assistant AssistantService
	storage   chat.Repository
}

func NewChatService(assistant AssistantService, storage chat.Repository) ChatService {
	return ChatService{assistant, storage}
}

func (c *ChatService) CreateChat(chatId, context string) error {
	return c.storage.CreateChat(chatId, context)
}

func (c *ChatService) SendMessage(chatId string, msg chat.Message) (chat.Message, error) {
	if err := c.storage.SendMessage(chatId, msg); err != nil {
		return chat.Message{}, err
	}
	buf := make([]byte, 0, maxResponseLength)
	buffer := bytes.NewBuffer(buf)
	if err := c.assistant.SendMessage(buffer, msg.Content(), chatId); err != nil {
		return chat.Message{}, err
	}
	content, err := io.ReadAll(buffer)
	if err != nil {
		return chat.Message{}, err
	}
	assistantMsg, err := chat.NewMessage(chat.Assistant, string(content))
	if err != nil {
		return chat.Message{}, err
	}
	if err := c.storage.SendMessage(chatId, assistantMsg); err != nil {
		return chat.Message{}, err
	}
	return assistantMsg, nil
}

func (c *ChatService) SendMessageStream(w io.Writer, chatId string, msg chat.Message) error {
	if err := c.storage.SendMessage(chatId, msg); err != nil {
		return err
	}
	sBuf := make([]byte, 0, maxResponseLength)
	storageBuf := bytes.NewBuffer(sBuf)
	multiWriter := io.MultiWriter(w, storageBuf)
	if err := c.assistant.SendMessageStream(multiWriter, msg.Content(), chatId); err != nil {
		return err
	}
	content, err := io.ReadAll(storageBuf)
	if err != nil {
		return err
	}
	assistantMsg, err := chat.NewMessage(chat.Assistant, string(content))
	if err != nil {
		return err
	}
	if err := c.storage.SendMessage(chatId, assistantMsg); err != nil {
		return err
	}
	return nil
}
