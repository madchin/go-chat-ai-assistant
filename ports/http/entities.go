package http_server

import "github.com/madchin/go-chat-ai-assistant/domain/chat"

type CustomerMessage struct {
	Content string `json:"content"`
}

type AssistantMessage struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type HttpError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type codeErrors struct {
	c string
}

var (
	serverCodeError = codeErrors{"server"}
	clientCodeError = codeErrors{"client"}
)

func (i CustomerMessage) toDomainMessage() (chat.Message, error) {
	msg, err := chat.NewCustomerMessage(i.Content)
	if err != nil {
		return chat.Message{}, err
	}
	return msg, nil
}

func mapDomainMessageToHttpAssistantMessage(domainMsg chat.Message) AssistantMessage {
	return AssistantMessage{Author: domainMsg.Author().Role(), Content: domainMsg.Content()}
}
