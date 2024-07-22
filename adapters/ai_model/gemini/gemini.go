package gemini

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/vertexai/genai"
	inmemory_storage "github.com/madchin/go-chat-ai-assistant/adapters/repository/in_memory"
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

func (g Model) SendMessageStream(w io.Writer, content, chatId string) error {
	ctx := context.Background()
	client, err := g.createClientConnection(ctx, chatId)
	if err != nil {
		return errClientGeneration
	}
	model := client.GenerativeModel(g.model)
	msg := genai.Text(content)
	streamIterator := model.GenerateContentStream(ctx, msg)
	for {
		resp, err := streamIterator.Next()
		if err == iterator.Done {
			return nil
		}
		if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
			return errEmptyResponse
		}
		if err != nil {
			return errUnexpected
		}
		for _, candidate := range resp.Candidates {
			for _, part := range candidate.Content.Parts {
				fmt.Fprintf(w, "%s ", part)
			}
		}
	}
}

func (g Model) SendMessage(w io.Writer, content, chatId string) error {
	ctx := context.Background()
	client, err := g.createClientConnection(ctx, chatId)
	if err != nil {
		return errClientGeneration
	}
	Model := client.GenerativeModel(g.model)
	msg := genai.Text(content)
	resp, err := Model.GenerateContent(ctx, msg)
	if err != nil {
		return err
	}
	for _, c := range resp.Candidates {
		for _, p := range c.Content.Parts {
			fmt.Fprintf(w, "%s ", p)
		}
	}
	return nil
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
