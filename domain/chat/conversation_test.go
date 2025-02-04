package chat

import (
	"fmt"
	"strings"
	"testing"
)

func TestEnqueueOneMessage(t *testing.T) {
	cnvrst := &conversation{}
	msg := Message{Customer, "first", 1}
	err := cnvrst.enqueue(msg)
	if err != nil {
		t.Fatalf("Expected to enqueue correctly, err occured: %v", err)
	}
	if cnvrst.first.message != msg {
		t.Fatalf("First message should equal %v, but is %v", msg, cnvrst.first.message)
	}
	if cnvrst.first.message != cnvrst.last.message {
		t.Fatalf("Conversation first and last message should be equal, but arent")
	}
	if cnvrst.first.message.content != msg.content {
		t.Fatalf("Enqueued message content is different than first message content")
	}
	if cnvrst.last.message.content != msg.content {
		t.Fatalf("Enqueued message content is different than last message content")
	}
	if cnvrst.first.next != nil {
		t.Fatalf("First message next should be nil after first enqueue")
	}
	if cnvrst.last.next != nil {
		t.Fatalf("Last message next should be nil after first enqueue")
	}
	if cnvrst.size != 1 {
		t.Fatalf("Conversation size should be equal 1, but is %d", cnvrst.size)
	}
}

func TestEnqueueTwoMessages(t *testing.T) {
	cnvrst := &conversation{}
	firstMsg := Message{Customer, "first", 1}
	err := cnvrst.enqueue(firstMsg)
	if err != nil {
		t.Fatalf("Expected to enqueue correctly, err occured: %v", err)
	}
	secondMsg := Message{Customer, "second", 2}
	err = cnvrst.enqueue(secondMsg)
	if err != nil {
		t.Fatalf("Expected to enqueue correctly, err occured: %v", err)
	}
	if cnvrst.first.message == cnvrst.last.message {
		t.Fatalf("First message should be different than last message")
	}
	if cnvrst.first.next != cnvrst.last {
		t.Fatalf("First next is not pointing to last node")
	}
	if cnvrst.last.next != nil {
		t.Fatalf("Last next should be nil, but is pointing to memory address somewhere")
	}
	if cnvrst.size != 2 {
		t.Fatalf("Conversation queue size should equal 2, but is %d", cnvrst.size)
	}
	if cnvrst.first.message != firstMsg {
		t.Fatalf("Conversation First message should equal %v, but is: %v", firstMsg, cnvrst.first.message)
	}
	if cnvrst.last.message != secondMsg {
		t.Fatalf("Conversation Last message should equal %v, but is: %v", secondMsg, cnvrst.last.message)
	}
}

func TestEnqueueThreeMessages(t *testing.T) {
	cnvrst := &conversation{}
	firstMsg := Message{Customer, "first", 1}
	err := cnvrst.enqueue(firstMsg)
	if err != nil {
		t.Fatalf("Expected to enqueue correctly, err occured: %v", err)
	}
	secondMsg := Message{Customer, "second", 2}
	err = cnvrst.enqueue(secondMsg)
	if err != nil {
		t.Fatalf("Expected to enqueue correctly, err occured: %v", err)
	}
	thirdMsg := Message{Customer, "third", 3}
	err = cnvrst.enqueue(thirdMsg)
	if err != nil {
		t.Fatalf("Expected to enqueue correctly, err occured: %v", err)
	}
	if cnvrst.first.message == cnvrst.last.message {
		t.Fatalf("First message should be different than last message")
	}
	if cnvrst.first.next == cnvrst.last {
		t.Fatalf("First next is pointing to last node, but should not")
	}
	if cnvrst.first.next.next != cnvrst.last {
		t.Fatalf("Node chain is broken, two nodes further from first node there should be last node")
	}
	if cnvrst.last.next != nil {
		t.Fatalf("Last next should be nil, but is pointing to memory address somewhere")
	}
	if cnvrst.size != 3 {
		t.Fatalf("Conversation queue size should equal 2, but is %d", cnvrst.size)
	}
	if cnvrst.first.message != firstMsg {
		t.Fatalf("Conversation First message should equal %v, but is: %v", firstMsg, cnvrst.first.message)
	}
	if cnvrst.first.next.message != secondMsg {
		t.Fatalf("Conversation Second message should equal %v, but is: %v", secondMsg, cnvrst.last.message)
	}
	if cnvrst.last.message != thirdMsg {
		t.Fatalf("Conversation Third message should equal %v, but is: %v", secondMsg, cnvrst.last.message)
	}
}

var testCasesExceedMaxConversationSize = []struct {
	size     int
	expected error
}{
	{MaxConversationSize, nil},
	{MaxConversationSize + 1, ErrConversationMaxSizeExceeded},
}

func TestEnqueueExceedMaxConversationSize(t *testing.T) {
	for idx, tc := range testCasesExceedMaxConversationSize {
		cnvrst := &conversation{}
		msg := Message{Customer, "msg", 1}
		var err error
		for i := 1; i <= tc.size; i++ {
			err = cnvrst.enqueue(msg)
		}
		if err != tc.expected {
			t.Fatalf("At test case: %d\nExpected error to be %v, actual: %v", idx+1, err, tc.expected)
		}
	}
}

func TestDequeueWhenConversationIsEmpty(t *testing.T) {
	cnvrst := &conversation{}
	_, err := cnvrst.dequeue()
	if err != ErrEmptyConversation {
		t.Fatalf("Expected error: %v, actual error: %v", ErrEmptyConversation, err)
	}
}

func TestDequeueLastElementInConversation(t *testing.T) {
	cnvrst := &conversation{}
	enqueueMsg, err := NewCustomerMessage("example")
	if err != nil {
		t.Fatalf("Unexpected error occured during creating message, err: %v", err)
	}
	cnvrst.enqueue(enqueueMsg)
	dequeueMsg, err := cnvrst.dequeue()
	if err != nil {
		t.Fatalf("Unexpected error occured during dequeueing last message in conversation, err: %v", err)
	}
	if enqueueMsg != dequeueMsg {
		t.Fatalf("Enqueued msg %v != dequeued msg %v", enqueueMsg, dequeueMsg)
	}
	if cnvrst.first != nil {
		t.Fatalf("First node in conversation should be nil")
	}
	if cnvrst.last != nil {
		t.Fatalf("Last node in conversation should be nil")
	}
	if cnvrst.size != 0 {
		t.Fatalf("Conversation size should be decremented after dequeue, expected size: %d actual size: %d", 1, cnvrst.size)
	}
}

func TestDequeueTwoElementsAfterEachOther(t *testing.T) {
	cnvrst := &conversation{}
	firstMsg, err := NewCustomerMessage(strings.Repeat("f", MinCharLength+1))
	if err != nil {
		t.Fatalf("Unexpected error occured during creating first message, err: %v", err)
	}
	secondMsg, err := NewCustomerMessage(strings.Repeat("s", MinCharLength+1))
	if err != nil {
		t.Fatalf("Unexpected error occured during creating second message, err: %v", err)
	}
	cnvrst.enqueue(firstMsg)
	cnvrst.enqueue(secondMsg)
	firstDequeued, _ := cnvrst.dequeue()
	secondDequeued, err := cnvrst.dequeue()
	if err != nil {
		t.Fatalf("Unexpected error occured during dequeueing second message, err: %v", err)
	}
	if firstDequeued == secondDequeued {
		t.Fatalf("first dequeued: %v == second dequeued: %v ", firstDequeued, secondDequeued)
	}
	if secondDequeued != secondMsg {
		t.Fatalf("Second enqueued msg: %v != second dequeued msg: %v", secondMsg, secondDequeued)
	}
	if cnvrst.size != 0 {
		t.Fatalf("Conversation size should equal to: %d, actual: %d", 0, cnvrst.size)
	}
	if cnvrst.first != nil {
		t.Fatalf("First node in conversation should be nil")
	}
	if cnvrst.last != nil {
		t.Fatalf("Last node in conversation should be nil")
	}
}

var testCasesPeek = []struct {
	msgCount int
	expected error
}{
	{0, ErrEmptyConversation},
	{1, nil},
	{2, nil},
}

func TestPeek(t *testing.T) {
	for idx, tc := range testCasesPeek {
		cnvrst := &conversation{}
		msgs, err := seedMessages(tc.msgCount)
		if err != nil {
			t.Fatalf("Test case: %d failed. Seeding messages err %v", idx+1, err)
		}
		seedConversation(cnvrst, msgs)
		peeked, err := cnvrst.peek()
		if err != tc.expected {
			t.Fatalf("Test case %d failed. Expected err to equal: %v, actual: %v", idx+1, tc.expected, err)
		}
		if len(msgs) > 0 && peeked != msgs[0] {
			t.Fatalf("Peeked msg is wrong. Expected msg: %v, actual msg: %v", msgs[0], peeked)
		}
	}
}

func seedConversation(cnvrst *conversation, msgs []Message) {
	for _, msg := range msgs {
		cnvrst.enqueue(msg)
	}
}

func seedMessages(count int) ([]Message, error) {
	msgs := make([]Message, count)
	for i := 0; i < count; i++ {
		content := strings.Repeat(fmt.Sprintf("%d", i), MinCharLength+1)
		msg, err := NewCustomerMessage(content)
		if err != nil {
			return nil, err
		}
		msgs[i] = msg
	}
	return msgs, nil
}
