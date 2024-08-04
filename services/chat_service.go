package service

import (
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

type ChatService struct {
	assistant chat.AssistantService
	storage   chat.StorageRepository
}

func NewChatService(assistant chat.AssistantService, storage chat.StorageRepository) *ChatService {
	return &ChatService{assistant, storage}
}

func (c *ChatService) CreateChat(chatId, context string) error {
	return c.storage.CreateChat(chatId, context)
}

func (c *ChatService) SendMessage(chatId, question string) (chat.Message, error) {
	customerMsg, err := chat.NewCustomerMessage(question)
	if err != nil {
		return chat.Message{}, err
	}
	if err := c.storage.SendMessage(chatId, customerMsg); err != nil {
		return chat.Message{}, err
	}
	assistantResponse, err := c.assistant.Ask(customerMsg.Content(), chatId)
	if err != nil {
		return chat.Message{}, err
	}
	if err := c.storage.SendMessage(chatId, assistantResponse); err != nil {
		return chat.Message{}, err
	}
	return assistantResponse, nil

}

func (c *ChatService) SendMessageStream(responseCh chan<- string, chatId, question string) error {
	customerMsg, err := chat.NewCustomerMessage(question)
	if err != nil {
		return err
	}
	if err := c.storage.SendMessage(chatId, customerMsg); err != nil {
		return err
	}
	assistantResponse, err := c.assistant.AskStream(responseCh, customerMsg.Content(), chatId)
	if err != nil {
		return err
	}
	if err := c.storage.SendMessage(chatId, assistantResponse); err != nil {
		return err
	}
	return nil
}
