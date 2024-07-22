package chat

import (
	"errors"
	"testing"
)

var testSize = []struct {
	size     int
	expected int
}{
	{0, 0},
	{1, 1},
	{2, 2},
	{MaxConversationSize, MaxConversationSize},
	{MaxConversationSize + 1, MaxConversationSize},
}

func TestSize(t *testing.T) {
	for idx, tc := range testSize {
		chat := NewChat("")
		msgs, err := seedMessages(tc.size)
		if err != nil {
			t.Fatalf("Test case: %d failed.\nCreating messages failed: %v", idx+1, err)
		}
		seedConversation(chat.conversation, msgs)
		if chat.Size() != tc.expected {
			t.Fatalf("Test case: %d failed. Expected size: %d, Actual size: %d", idx+1, tc.expected, chat.Size())
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
	{"remove when conversation is empty and first message is up to date", MaxMessageTime, 0, 0, ErrEmptyConversation},
	{"remove when conversation is full and first message is outdated", -1, MaxConversationSize, MaxConversationSize - 1, nil},
	{"remove when conversation is full and first message is up to date", MaxMessageTime, MaxConversationSize, MaxConversationSize, ErrMessageUpToDate},
	{"remove when conversation have capacity and first message is up to date", MaxMessageTime, 2, 2, ErrMessageUpToDate},
}

func TestRemoveMessage(t *testing.T) {
	for _, tc := range testRemoveMessage {
		chat := Chat{"", &conversation{}}
		seededMsgs, err := seedMessages(tc.conversationSize)
		if err != nil {
			t.Fatalf("Test case: %s failed.\nCreating messages failed: %v", tc.name, err)
		}
		seedConversation(chat.conversation, seededMsgs)
		_, err = chat.RemoveMessage(tc.olderThan)
		if chat.Size() != tc.expectedSizeAfterRemove {
			t.Fatalf("Test case: %s failed.\nConversation size is wrong. Expected: %d, Actual: %d", tc.name, tc.conversationSize, tc.expectedSizeAfterRemove)
		}
		if !errors.Is(err, tc.expectedErr) {
			t.Fatalf("Test case: %s failed.\nExpected error is not wrapped. expected: %v\n actual: %v", tc.name, tc.expectedErr, err)
		}
	}
}

var testConversation = []struct {
	name        string
	size        int
	expectedErr error
}{
	{"Retrieve messages when conversation is empty", 0, ErrEmptyConversation},
	{"Retrieve messages when conversation is full", MaxConversationSize, nil},
	{"Retrieve messages when conversation have 2 messages", 2, nil},
}

func TestConversation(t *testing.T) {
	for _, tc := range testConversation {
		chat := NewChat("")
		seededMsgs, err := seedMessages(tc.size)
		if err != nil {
			t.Fatalf("Test case: %s failed.\nCreating messages failed: %v", tc.name, err)
		}
		seedConversation(chat.conversation, seededMsgs)
		msgs, err := chat.Conversation()
		if !errors.Is(err, tc.expectedErr) {
			t.Fatalf("Test case: %s failed.\nExpected error is not wrapped. expected: %v\n actual: %v", tc.name, tc.expectedErr, err)
		}
		for idx, msg := range msgs {
			if msg != seededMsgs[idx] {
				t.Fatalf("Test case: %s failed.\n Retrieved message is different than should be.\n expected: %v, actual: %v", tc.name, seededMsgs[idx], msg)
			}
		}
	}
}
