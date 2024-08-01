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
	msg, err := chat.NewCustomerMessage(i.Content)
	if err != nil {
		return chat.Message{}, err
	}
	return msg, nil
}

func mapDomainMessageToOutcomingMessage(domainMsg chat.Message) OutcomingMessage {
	return OutcomingMessage{Author: domainMsg.Author().Role(), Content: domainMsg.Content()}
}
