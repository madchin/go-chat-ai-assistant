package assistant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"github.com/madchin/go-chat-ai-assistant/adapters/repository/cache"
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
	"google.golang.org/api/iterator"
)

type geminiAssistant struct {
	model, projectId, location string
	Connections                connections
}

func NewGemini(model, projectId, location string) *geminiAssistant {
	return &geminiAssistant{model, projectId, location, make(map[string]*connection, cache.ChatsCapacity)}
}

func (g *geminiAssistant) AskStream(responseCh chan<- string, content, chatId string) (chat.Message, error) {
	ctx := context.Background()
	client, err := g.createClientConnection(ctx, chatId)
	if err != nil {
		return chat.Message{}, err
	}
	model := client.GenerativeModel(g.model)
	phrase := genai.Text(content)
	streamIterator := model.GenerateContentStream(ctx, phrase)
	responseMsgContent := ""

	for {
		resp, err := streamIterator.Next()
		if err == iterator.Done {
			close(responseCh)
			break
		}
		if err != nil {
			close(responseCh)
			return chat.Message{}, err
		}
		for _, candidate := range resp.Candidates {
			for _, part := range candidate.Content.Parts {
				trimmedPart := strings.TrimSpace(fmt.Sprintf("%s", part))
				responseMsgContent += trimmedPart
				responseCh <- trimmedPart
			}
		}
	}

	return chat.NewAssistantMessage(responseMsgContent), nil
}

func (g *geminiAssistant) Ask(content, chatId string) (chat.Message, error) {
	ctx := context.Background()
	client, err := g.createClientConnection(ctx, chatId)
	if err != nil {
		return chat.Message{}, err
	}
	model := client.GenerativeModel(g.model)
	phrase := genai.Text(content)
	resp, err := model.GenerateContent(ctx, phrase)
	if err != nil {
		return chat.Message{}, err
	}
	responseMsgContent := ""
	for _, c := range resp.Candidates {
		for _, part := range c.Content.Parts {
			trimmedPart := strings.TrimSpace(fmt.Sprintf("%s", part))
			responseMsgContent += trimmedPart
		}
	}
	return chat.NewAssistantMessage(responseMsgContent), nil
}

func (g *geminiAssistant) createClientConnection(ctx context.Context, chatId string) (*connection, error) {
	if isChatConnectionEstablished(chatId, g.Connections) {
		return g.Connections[chatId], nil
	}
	client, err := genai.NewClient(ctx, g.projectId, g.location)
	if err != nil {
		return nil, err
	}
	connection := NewConnection(client, time.Now().Unix())
	g.Connections[chatId] = connection

	return connection, nil
}
