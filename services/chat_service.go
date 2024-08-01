package service

import (
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

type AssistantService interface {
	SendMessage(content, chatId string) (chat.Message, error)
	SendMessageStream(response chan<- string, content, chatId string) (chat.Message, error)
}

type ChatService struct {
	assistant AssistantService
	storage   chat.Repository
}

func NewChatService(assistant AssistantService, storage chat.Repository) *ChatService {
	return &ChatService{assistant, storage}
}

func (c *ChatService) CreateChat(chatId, context string) error {
	return c.storage.CreateChat(chatId, context)
}

func (c *ChatService) SendMessage(chatId string, customerMsg chat.Message) (chat.Message, error) {
	if err := c.storage.SendMessage(chatId, customerMsg); err != nil {
		return chat.Message{}, err
	}
	assistantResponse, err := c.assistant.SendMessage(customerMsg.Content(), chatId)
	if err != nil {
		return chat.Message{}, err
	}
	if err := c.storage.SendMessage(chatId, assistantResponse); err != nil {
		return chat.Message{}, err
	}
	return assistantResponse, nil

}

func (c *ChatService) SendMessageStream(responseCh chan<- string, chatId string, customerMsg chat.Message) error {
	if err := c.storage.SendMessage(chatId, customerMsg); err != nil {
		return err
	}
	assistantResponse, err := c.assistant.SendMessageStream(responseCh, customerMsg.Content(), chatId)
	if err != nil {
		return err
	}
	if err := c.storage.SendMessage(chatId, assistantResponse); err != nil {
		return err
	}
	return nil
}
