package chat

import (
	"strings"
	"testing"
	"time"
)

func TestMessageIsUpToDateAfterCreation(t *testing.T) {
	msg, _ := NewCustomerMessage("example")
	isOutdated := msg.isOutdated(time.Now().UnixMilli(), MaxMessageTime)

	if isOutdated {
		t.Fatalf("message should be up to date, now its outdated. msg: %v", msg)
	}
}

var testCasesOutdate = []struct {
	currentTime int64
	messageTime int64
	outdateTime int64
	expected    bool
}{
	{100, 50, 49, true},
	{100, 50, 50, true},
	{100, 50, 51, false},
	{100, 50, 60, false},
}

func TestMessageIsOutdated(t *testing.T) {
	for idx, tc := range testCasesOutdate {
		msg := Message{Customer, "", tc.messageTime}
		actual := msg.isOutdated(tc.currentTime, tc.outdateTime)
		if actual != tc.expected {
			t.Fatalf(`
			test failed at test case %d
			message should be outdated, but is up to date
			messageTime: %d
			currentTime: %d
			outdateTime: %d
			expected to be out dated: %v
			`, idx+1, tc.messageTime, tc.currentTime, tc.outdateTime, tc.expected)
		}
	}
}

var testCasesContentValidation = []struct {
	content  string
	expected error
}{
	{strings.Repeat("a", minCharLength), ErrMessageShort},
	{strings.Repeat("a", minCharLength+1), nil},
	{"a", ErrMessageShort},
	{"", ErrMessageShort},
	{strings.Repeat("a", maxCharLength), nil},
	{strings.Repeat("a", maxCharLength+1), ErrMessageLong},
}

func TestMessageContentValidationError(t *testing.T) {
	for idx, tc := range testCasesContentValidation {
		_, err := NewCustomerMessage(tc.content)
		if err != tc.expected {
			t.Fatalf("At test case %d expected err: %v, actual err: %v", idx+1, tc.expected, err)
		}
	}
}
