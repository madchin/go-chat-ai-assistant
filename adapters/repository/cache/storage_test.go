package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

func TestAddingChat(t *testing.T) {
	chatId := "chatId"
	chatContext := "chatContext"
	storage := newStorageWithChat(chatId, chatContext, t)

	if chat := storage.chat(chatId); chat == nil {
		t.Fatalf("chat has not been added to the storage")
	}
}

func TestAddingExistingChat(t *testing.T) {
	chatId := "chatId"
	chatContext := "chatContext"
	storage := newStorageWithChat(chatId, chatContext, t)
	err := storage.CreateChat(chatId, chatContext)
	if err != ErrChatAlreadyExists {
		t.Fatalf("error should be returned when creating chat for same user. \n Expected: %v \n Actual: %v", ErrChatAlreadyExists, err)
	}
}

var testCasesAddingChatWithContext = []struct {
	name    string
	context string
}{
	{"context provided", "randomContext"},
	{"context is empty", ""},
}

func TestAddingChatWithContext(t *testing.T) {
	chatId := "chatId"
	for _, tc := range testCasesAddingChatWithContext {
		storage := newStorageWithChat(chatId, tc.context, t)
		chat := storage.chat(chatId)
		if chat.Context() != tc.context {
			t.Fatalf("Test case with name: %s failed.\n Context is wrong. \n Expected: %s\n Actual: %s", tc.name, tc.context, chat.Context())
		}
	}
}

func TestSendMessageToExistingChat(t *testing.T) {
	chatId := "chatId"
	storage := newStorageWithChat(chatId, "", t)
	err := storage.SendMessage(chatId, chat.NewValidMessage())
	if err != nil {
		t.Fatalf("error occured during sending message. err: %v", err)
	}
	if storage.chats[chatId].Size() != 1 {
		t.Fatalf("message has not been sent. actual conversation size: %d, expected: %d", storage.chats[chatId].Size(), 1)
	}
}

func TestSendMessageToNotExistingChat(t *testing.T) {
	chatId := "chatId"
	storage := New()
	msg := chat.NewValidMessage()
	err := storage.SendMessage(chatId, msg)
	if err != errChatNotExists {
		t.Fatalf("unexpected error occured\n Expected %v \n Actual: %v", errChatNotExists, err)
	}
}

func TestRemovingMessageWhenStorageNotHaveChats(t *testing.T) {
	storage := New()
	msgs := storage.RemoveOutdatedMessages()
	if len(msgs) > 0 {
		t.Fatalf("%d count of messages have been removed when storage had no chats", len(msgs))
	}
}

func TestRemovingMessageWhenChatNotHaveMessages(t *testing.T) {
	chatId := "chatId"
	storage := newStorageWithChat(chatId, "", t)
	msgs := storage.RemoveOutdatedMessages()
	if len(msgs[chatId]) > 0 {
		t.Fatalf("%d count of messages have been removed from chat which should not have any messages", len(msgs))
	}
}

func TestRemovingOutdatedMessages(t *testing.T) {
	chatId := "chatId"
	storage := newStorageWithChat(chatId, "", t)
	seedChatWithMessages(storage, chatId, outdatedMessages(chat.MaxConversationSize))
	msgs := storage.RemoveOutdatedMessages()
	if len(msgs[chatId]) != chat.MaxConversationSize {
		t.Fatalf("wrong count of messages has been removed. \n Expected: %d \n Actual: %d", chat.MaxConversationSize, len(msgs[chatId]))
	}
}

func TestRemovingUpToDateMessages(t *testing.T) {
	chatId := "chatId"
	storage := newStorageWithChat(chatId, "", t)
	seedChatWithMessages(storage, chatId, upToDateMessages(chat.MaxConversationSize))
	msgs := storage.RemoveOutdatedMessages()
	if len(msgs[chatId]) != 0 {
		t.Fatalf("any of the message should not be deleted. Deletion count: \n Expected: %d \n Actual: %d", 0, len(msgs[chatId]))
	}
}

var testCasesRetrievingAllConversations = []struct {
	name      string
	msgs      []chat.Message
	chatCount int
}{
	{"should not retrieve any conversation when there is no chats", nil, 0},
	{"should not retrieve any conversation when there are chats but empty", nil, ChatsCapacity / 2},
	{"should retrieve all up to date messages from all chats", upToDateMessages(chat.MaxConversationSize), ChatsCapacity / 2},
	{"should retrieve all outdated messages from all chats", outdatedMessages(chat.MaxConversationSize), ChatsCapacity / 2},
	{
		"should retrieve all outdated and up to date messages from all chats",
		[]chat.Message{
			chat.NewValidMessageWithTimestamp(time.Now().UnixMilli()),
			chat.NewValidMessageWithTimestamp(time.Now().UnixMilli() - chat.MaxMessageTime*2),
		}, ChatsCapacity / 2,
	},
}

func TestRetrievingAllConversations(t *testing.T) {
	for _, tc := range testCasesRetrievingAllConversations {
		storage := New()
		for i := 0; i < tc.chatCount; i++ {
			chatId := fmt.Sprintf("chat %d", i)
			storage.CreateChat(chatId, "")
			seedChatWithMessages(storage, chatId, tc.msgs)
		}
		expectedMsgRetrieveCount := tc.chatCount * len(tc.msgs)
		actualMsgRetrieveCount := allConversationsMessageCount(storage)
		if expectedMsgRetrieveCount != actualMsgRetrieveCount {
			t.Fatalf("Test case with name %s failed. \n Wrong count of messages retrieved. \n Expected: %d \n Actual: %d", tc.name, expectedMsgRetrieveCount, actualMsgRetrieveCount)
		}
	}
}

func newStorageWithChat(chatId, context string, t *testing.T) *storage {
	storage := New()
	err := storage.CreateChat(chatId, context)
	if err != nil {
		t.Fatalf("unexpected error occured during adding new chat, %v", err)
	}
	return storage
}

func seedChatWithMessages(storage *storage, chatId string, msgs []chat.Message) {
	for _, msg := range msgs {
		storage.SendMessage(chatId, msg)
	}
}

func upToDateMessages(count int) []chat.Message {
	msgs := make([]chat.Message, count)
	for i := 0; i < count; i++ {
		msg := chat.NewValidMessageWithTimestamp(time.Now().UnixMilli())
		msgs[i] = msg
	}
	return msgs
}
func outdatedMessages(count int) []chat.Message {
	msgs := make([]chat.Message, count)
	outdatedTimestamp := time.Now().UnixMilli() - chat.MaxMessageTime*2
	for i := 0; i < count; i++ {
		msg := chat.NewValidMessageWithTimestamp(outdatedTimestamp)
		msgs[i] = msg
	}
	return msgs
}

func allConversationsMessageCount(storage *storage) (count int) {
	retrievedMsgs, _ := storage.RetrieveAllConversations()
	for _, msgs := range retrievedMsgs {
		count += len(msgs)
	}
	return
}
