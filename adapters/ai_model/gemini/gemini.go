package gemini

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/vertexai/genai"
	inmemory_storage "github.com/madchin/go-chat-ai-assistant/adapters/repository/in_memory"
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
	"google.golang.org/api/iterator"
)

var (
	errClientGeneration = errors.New("client generation error")
	errEmptyResponse    = errors.New("empty response from model")
	errUnexpected       = errors.New("unexpected")
)

type Model struct {
	model, projectId, location string
	Connections                connections
}

func NewModel(model, projectId, location string) Model {
	return Model{model, projectId, location, make(map[string]*connection, inmemory_storage.ChatsCapacity)}
}

func (g Model) SendMessageStream(responseCh chan<- string, content, chatId string) (msg chat.Message, err error) {
	ctx := context.Background()
	client, err := g.createClientConnection(ctx, chatId)
	if err != nil {
		err = errClientGeneration
		return
	}
	model := client.GenerativeModel(g.model)
	phrase := genai.Text(content)
	streamIterator := model.GenerateContentStream(ctx, phrase)
	responseMsgContent := ""
	for {
		resp, err := streamIterator.Next()
		if resp == nil {
			break
		}
		if err == iterator.Done {
			break
		}
		if len(resp.Candidates) == 0 {
			err = errEmptyResponse
			break
		}
		if err != nil {
			err = errUnexpected
			break
		}
		for _, candidate := range resp.Candidates {
			for _, part := range candidate.Content.Parts {
				trimmedPart := strings.TrimSpace(fmt.Sprintf("%s", part))
				responseMsgContent += trimmedPart
				responseCh <- trimmedPart
			}
		}
	}
	close(responseCh)
	msg = chat.NewAssistantMessage(responseMsgContent)
	return
}

func (g Model) SendMessage(content, chatId string) (msg chat.Message, err error) {
	ctx := context.Background()
	client, err := g.createClientConnection(ctx, chatId)
	if err != nil {
		err = errClientGeneration
		return
	}
	model := client.GenerativeModel(g.model)
	phrase := genai.Text(content)
	resp, err := model.GenerateContent(ctx, phrase)
	if err != nil {
		return
	}
	responseMsgContent := ""
	for _, c := range resp.Candidates {
		for _, part := range c.Content.Parts {
			trimmedPart := strings.TrimSpace(fmt.Sprintf("%s", part))
			responseMsgContent += trimmedPart
		}
	}
	msg = chat.NewAssistantMessage(responseMsgContent)
	return
}

func (g Model) createClientConnection(ctx context.Context, chatId string) (*connection, error) {
	client, err := genai.NewClient(ctx, g.projectId, g.location)
	if err != nil {
		return nil, errClientGeneration
	}
	connection := NewConnection(client, time.Now().Unix())
	err = g.Connections.bindChatToClientConnection(chatId, connection)
	if err != errChatConnectionAlreadyEstablished && err != nil {
		return nil, err
	}
	return connection, nil
}
