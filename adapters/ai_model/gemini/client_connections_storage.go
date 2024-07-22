package gemini

import (
	"errors"
	"time"

	"cloud.google.com/go/vertexai/genai"
)

var (
	errChatConnectionAlreadyEstablished = errors.New("chat is already connected with model")
)

const expiryTime = 240_000

type connection struct {
	establishTime int64
	*genai.Client
}

type connections map[string]*connection

func NewConnection(client *genai.Client, establishTime int64) *connection {
	return &connection{establishTime, client}
}

func (connections connections) bindChatToClientConnection(chatId string, connection *connection) error {
	if isChatConnectionEstablished(chatId, connections) {
		return errChatConnectionAlreadyEstablished
	}
	connections[chatId] = connection
	return nil
}

func (connections connections) RemoveOutdatedConnections() {
	for id, connection := range connections {
		if connection.isOutdated() {
			connection.Close()
			connections.removeConnection(id)
		}
	}
}

func (connections connections) removeConnection(chatId string) {
	if !isChatConnectionEstablished(chatId, connections) {
		return
	}
	delete(connections, chatId)
}

func isChatConnectionEstablished(chatId string, connections connections) bool {
	_, ok := connections[chatId]
	return ok
}

func (connection connection) isOutdated() bool {
	return time.Now().Unix()-connection.establishTime+expiryTime > 0
}
