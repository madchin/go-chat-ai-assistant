package chat

import (
	"errors"
	"strings"
	"testing"
)

var testSize = []struct {
	size     int
	expected int
}{
	{0, 0},
	{1, 1},
	{2, 2},
	{maxConversationSize, maxConversationSize},
	{maxConversationSize + 1, maxConversationSize},
}

func TestSize(t *testing.T) {
	for idx, tc := range testSize {
		chat := NewChat("")
		seedConversation(chat.conversation, seedMessages(tc.size))
		if chat.Size() != tc.expected {
			t.Fatalf("Test case: %d failed. Expected size: %d, Actual size: %d", idx+1, tc.expected, chat.Size())
		}
	}
}

var testSendMessage = []struct {
	content  string
	expected error
}{
	{strings.Repeat("a", minCharLength), ErrMessageShort},
	{strings.Repeat("a", minCharLength-1), ErrMessageShort},
	{strings.Repeat("a", minCharLength+1), nil},
	{"", ErrMessageShort},
	{strings.Repeat("a", maxCharLength-1), nil},
	{strings.Repeat("a", maxCharLength), nil},
	{strings.Repeat("a", maxCharLength+1), ErrMessageLong},
}

func TestSendMessage(t *testing.T) {
	for idx, tc := range testSendMessage {
		chatId := "random"
		chat := Chat{chatId, "", &conversation{}}
		err := chat.SendMessage(NewMessage(tc.content))
		if !errors.Is(err, tc.expected) {
			t.Fatalf("Test case: %d failed.\nExpected error is not wrapped. expected: %v\n actual: %v", idx+1, tc.expected, err)
		}
		if err != nil && !errors.Is(err, ErrChat{chatId}) {
			t.Fatalf("Test case: %d failed.\nExpected error is not with chat context. expected: %v\n actual: %v", idx+1, ErrChat{chatId}, err)
		}
	}
}

var testRemoveMessage = []struct {
	name                    string
	olderThan               int64
	conversationSize        int
	expectedSizeAfterRemove int
	expectedErr             error
}{
	{"remove when conversation is empty and first message is outdated", -1, 0, 0, ErrEmptyConversation},
	{"remove when conversation is empty and first message is up to date", messageMaxTime, 0, 0, ErrEmptyConversation},
	{"remove when conversation is full and first message is outdated", -1, maxConversationSize, maxConversationSize - 1, nil},
	{"remove when conversation is full and first message is up to date", messageMaxTime, maxConversationSize, maxConversationSize, ErrMessageUpToDate},
	{"remove when conversation have capacity and first message is up to date", messageMaxTime, 2, 2, ErrMessageUpToDate},
}

func TestRemoveMessage(t *testing.T) {
	for _, tc := range testRemoveMessage {
		chatId := "random"
		chat := Chat{chatId, "", &conversation{}}
		seedConversation(chat.conversation, seedMessages(tc.conversationSize))
		_, err := chat.RemoveMessage(tc.olderThan)
		if chat.Size() != tc.expectedSizeAfterRemove {
			t.Fatalf("Test case: %s failed.\nConversation size is wrong. Expected: %d, Actual: %d", tc.name, tc.conversationSize, tc.expectedSizeAfterRemove)
		}
		if !errors.Is(err, tc.expectedErr) {
			t.Fatalf("Test case: %s failed.\nExpected error is not wrapped. expected: %v\n actual: %v", tc.name, tc.expectedErr, err)
		}
		if err != nil && !errors.Is(err, ErrChat{chatId}) {
			t.Fatalf("Test case: %s failed.\nExpected error is not with chat context. expected: %v\n actual: %v", tc.name, ErrChat{chatId}, err)
		}
	}
}

var testConversation = []struct {
	name        string
	size        int
	expectedErr error
}{
	{"Retrieve messages when conversation is empty", 0, ErrEmptyConversation},
	{"Retrieve messages when conversation is full", maxConversationSize, nil},
	{"Retrieve messages when conversation have 2 messages", 2, nil},
}

func TestConversation(t *testing.T) {
	for _, tc := range testConversation {
		chat := NewChat("")
		seededMsgs := seedMessages(tc.size)
		seedConversation(chat.conversation, seededMsgs)
		msgs, err := chat.Conversation()
		if !errors.Is(err, tc.expectedErr) {
			t.Fatalf("Test case: %s failed.\nExpected error is not wrapped. expected: %v\n actual: %v", tc.name, tc.expectedErr, err)
		}
		if err != nil && !errors.Is(err, ErrChat{chat.id}) {
			t.Fatalf("Test case: %s failed.\nExpected error is not with chat context. expected: %v\n actual: %v", tc.name, ErrChat{chat.id}, err)
		}
		for idx, msg := range msgs {
			if msg != seededMsgs[idx] {
				t.Fatalf("Test case: %s failed.\n Retrieved message is different than should be.\n expected: %v, actual: %v", tc.name, seededMsgs[idx], msg)
			}
		}
	}
}
