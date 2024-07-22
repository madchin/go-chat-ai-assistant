package http_server

import "github.com/madchin/go-chat-ai-assistant/domain/chat"

type IncomingMessage struct {
	Content string `json:"content"`
}

type OutcomingMessage struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type HttpError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (i IncomingMessage) toDomainMessage() (chat.Message, error) {
	msg, err := chat.NewMessage(chat.Customer, i.Content)
	if err != nil {
		return chat.Message{}, err
	}
	return msg, nil
}
