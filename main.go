package main

import (
	"github.com/madchin/go-chat-ai-assistant/adapters/ai_model/gemini"
	inmemory_storage "github.com/madchin/go-chat-ai-assistant/adapters/repository/in_memory"
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
	http_server "github.com/madchin/go-chat-ai-assistant/ports/http"
	service "github.com/madchin/go-chat-ai-assistant/services"
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
	model := gemini.NewModel("gemini-1.5-flash-001", "randomProjectId", "us-central1")
	application := service.NewApplication(storage, &historyRepositoryMock{}, storage, model)
	http_server.NewHttpServer(&application.ChatService)

	// err := storage.CreateChat("1", "")
	// if err != nil {
	// 	fmt.Printf("err is %v", err)
	// }
	// storage.CreateChat("2", "")
	// storage.CreateChat("3", "")
	// storage.CreateChat("4", "")

	// err = storage.SendMessage("1", chat.NewMessage(chat.Customer, "1111111"))
	// if err != nil {
	// 	fmt.Printf("err is %v", err)
	// }
	// storage.SendMessage("2", chat.NewMessage(chat.Assistant, "1111111"))
	// storage.SendMessage("2", chat.NewMessage(chat.Assistant, "2222222"))
	// storage.SendMessage("2", chat.NewMessage(chat.Assistant, "3333333"))
	// storage.SendMessage("3", chat.NewMessage(chat.Customer, "111121111"))
	// storage.SendMessage("3", chat.NewMessage(chat.Assistant, "2222222"))
	// storage.SendMessage("4", chat.NewMessage(chat.Customer, "1111111"))
	// for {
	// 	usersMsgs, _ := storage.RetrieveAllConversations()
	// 	fmt.Printf("actual chat count %d\n\n", len(usersMsgs))
	// 	for cId, msg := range usersMsgs {
	// 		fmt.Printf("Actual msg count %d in chat %s \n\n", len(msg), cId)
	// 	}
	// 	fmt.Println("\n\nNEXT:\n\n")
	// 	time.Sleep(time.Second * 1)
	// 	if len(usersMsgs) == 0 {
	// 		break
	// 	}
	// }
	// storage.SendMessage("3", chat.NewMessage(chat.Customer, "111121111"))
	// storage.SendMessage("3", chat.NewMessage(chat.Assistant, "2222222"))
	// storage.SendMessage("4", chat.NewMessage(chat.Customer, "1111111"))

	// for {
	// 	usersMsgs, _ := storage.RetrieveAllConversations()
	// 	fmt.Printf("actual chat count %d\n\n", len(usersMsgs))
	// 	for cId, msg := range usersMsgs {
	// 		fmt.Printf("Actual msg count %d in chat %s \n\n", len(msg), cId)
	// 	}
	// 	fmt.Println("\nNEXT:\n")
	// 	time.Sleep(time.Second * 3)
	// }

}
