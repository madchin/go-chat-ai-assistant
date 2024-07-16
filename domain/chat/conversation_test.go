package chat

import "testing"

func TestConversationEnqueueOneMessage(t *testing.T) {
	cnvrst := &conversation{}
	msg := Message{"first", 1}
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

func TestConversationEnqueueTwoMessages(t *testing.T) {
	cnvrst := &conversation{}
	firstMsg := Message{"first", 1}
	err := cnvrst.enqueue(firstMsg)
	if err != nil {
		t.Fatalf("Expected to enqueue correctly, err occured: %v", err)
	}
	secondMsg := Message{"second", 2}
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

func TestConversationEnqueueThreeMessages(t *testing.T) {
	cnvrst := &conversation{}
	firstMsg := Message{"first", 1}
	err := cnvrst.enqueue(firstMsg)
	if err != nil {
		t.Fatalf("Expected to enqueue correctly, err occured: %v", err)
	}
	secondMsg := Message{"second", 2}
	err = cnvrst.enqueue(secondMsg)
	if err != nil {
		t.Fatalf("Expected to enqueue correctly, err occured: %v", err)
	}
	thirdMsg := Message{"third", 3}
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
	{maxConversationSize, nil},
	{maxConversationSize + 1, ErrConversationMaxSizeExceeded},
}

func TestConversationEnqueueExceedMaxConversationSize(t *testing.T) {
	for idx, tc := range testCasesExceedMaxConversationSize {
		cnvrst := &conversation{}
		msg := Message{"msg", 1}
		var err error
		for i := 1; i <= tc.size; i++ {
			err = cnvrst.enqueue(msg)
		}
		if err != tc.expected {
			t.Fatalf("At test case: %d\nExpected error to be %v, actual: %v", idx+1, err, tc.expected)
		}
	}
}
