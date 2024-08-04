package history

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
	"google.golang.org/api/iterator"
)

const (
	firebase_write_delay = 1000
	max_retries          = 5
)

var (
	errSaveHistoryRetriesExceeded = errors.New("save history retries count has been exceeded")
)

type history struct {
	client *firestore.Client
}

type Message struct {
	Author  string `firestore:"author"`
	Content string `firestore:"content"`
}

func newStorageMessage(author chat.Participant, content string) Message {
	return Message{author.Role(), content}
}

// timestampToMessage
type messages map[string]Message

// chatIdsToMessages
type chatHistories map[string]messages

func New(cfg *firebase.Config) *history {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, cfg)
	if err != nil {
		panic(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}
	return &history{client}
}

func (h *history) SaveHistory(chatMessages chat.ChatMessages) error {
	storageChats := mapDomainChatMessagesToStorageChats(chatMessages)
	retries := 0
	return h.trySaveHistory(storageChats, retries)
}

func (h *history) RetrieveAllChatsHistoryStream(responseCh chan<- chat.ChatMessages) error {
	chatsCollection := h.chatCollection()
	iter := chatsCollection.Documents(context.Background())
	defer iter.Stop()
	for {
		chatMessages := make(chat.ChatMessages, 1)
		doc, err := iter.Next()
		if err == iterator.Done {
			close(responseCh)
			break
		}
		if err != nil {
			close(responseCh)
			return err
		}
		var messages messages
		err = doc.DataTo(&messages)
		if err != nil {
			return err
		}
		chatMessages[doc.Ref.ID] = mapStorageMessagesToDomainMessages(messages)
		responseCh <- chatMessages
	}

	return nil
}

func (h *history) RetrieveAllChatsHistory() (chat.ChatMessages, error) {
	chatsCollection := h.chatCollection()
	iter := chatsCollection.Documents(context.Background())
	chatMessages := make(chat.ChatMessages, 100)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return chatMessages, err
		}
		var messages messages
		err = doc.DataTo(&messages)
		if err != nil {
			return chatMessages, err
		}
		chatMessages[doc.Ref.ID] = mapStorageMessagesToDomainMessages(messages)
	}
	return chatMessages, nil
}

func (h *history) trySaveHistory(histories chatHistories, retries int) error {
	if retries >= max_retries {
		return errSaveHistoryRetriesExceeded
	}
	retries++
	unsaved := make(chatHistories, len(histories))
	for chatId, messages := range histories {
		err := h.addMessages(chatId, messages)
		if err != nil {
			unsaved[chatId] = histories[chatId]
		}
	}
	if len(unsaved) > 0 {
		<-time.After(time.Duration(time.Millisecond * firebase_write_delay))
		return h.trySaveHistory(unsaved, retries)
	}
	return nil
}

func (h *history) addMessages(chatId string, messages messages) error {
	if _, err := h.chatCollection().Doc(chatId).Set(context.Background(), messages, firestore.MergeAll); err != nil {
		return err
	}
	return nil
}

func (h *history) chatCollection() *firestore.CollectionRef {
	return h.client.Collection("chats")
}

func (m Message) authorToParticipant() chat.Participant {
	if m.Author == chat.Assistant.Role() {
		return chat.Assistant
	}
	return chat.Customer
}

func mapDomainChatMessagesToStorageChats(domainChatMessages chat.ChatMessages) chatHistories {
	chats := make(chatHistories, len(domainChatMessages))
	for chatId, domainMessages := range domainChatMessages {
		storageMessages := make(messages, len(domainMessages))
		for i := 0; i < len(domainMessages); i++ {
			author, content, timestamp := domainMessages[i].Author(), domainMessages[i].Content(), fmt.Sprintf("%v", domainMessages[i].Timestamp())
			storageMessages[timestamp] = newStorageMessage(author, content)
		}
		chats[chatId] = storageMessages
	}
	return chats
}

func mapStorageMessagesToDomainMessages(storageMessages messages) []chat.Message {
	chatMessages := make([]chat.Message, 0, len(storageMessages))
	for timestamp, storageMessage := range storageMessages {
		parsedTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			parsedTimestamp = 0
		}
		chatMessages = append(chatMessages, chat.NewMessage(storageMessage.authorToParticipant(), storageMessage.Content, parsedTimestamp))
	}
	return chatMessages
}
