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

type descriptionErrors struct {
	d string
}

var (
	serverCodeError = codeErrors{"server"}
	clientCodeError = codeErrors{"client"}

	genericDescriptionError = descriptionErrors{"Oops! Something went wrong"}
	wrongContentTypeError   = descriptionErrors{"Content-Type header need to be application/json"}
	bodyParseError          = descriptionErrors{"JSON parse failed"}
	wrongCookieValue        = descriptionErrors{"Wrong cookie value"}
)

func mapDomainMessageToHttpAssistantMessage(domainMsg chat.Message) AssistantMessage {
	return AssistantMessage{Author: domainMsg.Author().Role(), Content: domainMsg.Content()}
}
