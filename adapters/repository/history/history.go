package history

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
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
type chats map[string]messages

func New(cfg *firebase.Config) *history {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return &history{client}
}

func (h *history) SaveHistory(chatMessages chat.ChatMessages) error {
	storageChats := mapDomainChatMessagesToStorageChats(chatMessages)
	retries := 0
	return h.trySaveHistory(storageChats, retries)
}

func (h *history) RetrieveAllChatsHistory() (chat.ChatMessages, error) {
	return chat.ChatMessages{}, nil
	// chats := h.chatCollection()
	// retrieved := make(chat.ChatMessages, 100)
	// iter := chats.Documents(context.Background())
	// for {
	// 	doc, err := iter.Next()
	// 	if err == iterator.Done {
	// 		return retrieved, nil
	// 	}
	// 	if err != nil {
	// 		return retrieved, err
	// 	}
	// 	err = doc.DataTo(&retrieved)
	// 	if err != nil {
	// 		return retrieved, err
	// 	}
	// }
}

func (h *history) trySaveHistory(chatHistories chats, retries int) error {
	if retries >= max_retries {
		return errSaveHistoryRetriesExceeded
	}
	retries++
	unsaved := make(chats, len(chatHistories))
	for chatId, messages := range chatHistories {
		log.Printf("trying to save messages: %v\n", messages)
		err := h.addMessages(chatId, messages)
		if err != nil {
			log.Printf("failed to save messages: %v\n", err)
			unsaved[chatId] = chatHistories[chatId]
		}
	}
	if len(unsaved) > 0 {
		<-time.After(time.Duration(time.Millisecond * firebase_write_delay))
		log.Println("retrying")
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

func mapDomainChatMessagesToStorageChats(domainChatMessages chat.ChatMessages) chats {
	chats := make(chats, len(domainChatMessages))
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
