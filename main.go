package main

import (
	"fmt"
	"time"

	inmemory_storage "github.com/madchin/go-chat-ai-assistant/adapters/repository/in_memory"
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
	"github.com/madchin/go-chat-ai-assistant/services"
)

type historyRepositoryMock struct{}

func (h *historyRepositoryMock) SaveHistory(um chat.UserMessages) error {
	return nil
}

func (h *historyRepositoryMock) RetrieveAllChatsHistory() (chat.UserMessages, error) {
	return chat.UserMessages{}, nil
}

func main() {
	storage := inmemory_storage.New()
	services.Application(storage, &historyRepositoryMock{})
	err := storage.AddNewChat("1", "")
	if err != nil {
		fmt.Printf("err is %v", err)
	}
	storage.AddNewChat("2", "")
	storage.AddNewChat("3", "")
	storage.AddNewChat("4", "")

	err = storage.SendMessage("1", chat.NewMessage(chat.Customer, "1111111"))
	if err != nil {
		fmt.Printf("err is %v", err)
	}
	storage.SendMessage("2", chat.NewMessage(chat.Assistant, "1111111"))
	storage.SendMessage("2", chat.NewMessage(chat.Assistant, "2222222"))
	storage.SendMessage("2", chat.NewMessage(chat.Assistant, "3333333"))
	storage.SendMessage("3", chat.NewMessage(chat.Customer, "111121111"))
	storage.SendMessage("3", chat.NewMessage(chat.Assistant, "2222222"))
	storage.SendMessage("4", chat.NewMessage(chat.Customer, "1111111"))
	for {
		usersMsgs, _ := storage.RetrieveAllConversations()
		fmt.Printf("actual chat count %d\n\n", len(usersMsgs))
		for cId, msg := range usersMsgs {
			fmt.Printf("Actual msg count %d in chat %s \n\n", len(msg), cId)
		}
		fmt.Println("\n\nNEXT:\n\n")
		time.Sleep(time.Second * 1)
		if len(usersMsgs) == 0 {
			break
		}
	}
	storage.SendMessage("3", chat.NewMessage(chat.Customer, "111121111"))
	storage.SendMessage("3", chat.NewMessage(chat.Assistant, "2222222"))
	storage.SendMessage("4", chat.NewMessage(chat.Customer, "1111111"))
	for {
		usersMsgs, _ := storage.RetrieveAllConversations()
		fmt.Printf("actual chat count %d\n\n", len(usersMsgs))
		for cId, msg := range usersMsgs {
			fmt.Printf("Actual msg count %d in chat %s \n\n", len(msg), cId)
		}
		fmt.Println("\n\nNEXT:\n\n")
		time.Sleep(time.Second * 1)
		if len(usersMsgs) == 0 {
			break
		}
	}
}
