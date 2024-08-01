package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"github.com/madchin/go-chat-ai-assistant/adapters/ai_model/gemini"
	"github.com/madchin/go-chat-ai-assistant/adapters/repository/cache"
	"github.com/madchin/go-chat-ai-assistant/adapters/repository/history"
	grpc_server "github.com/madchin/go-chat-ai-assistant/ports/gRPC"
	http_server "github.com/madchin/go-chat-ai-assistant/ports/http"

	service "github.com/madchin/go-chat-ai-assistant/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	history := history.New(&firebase.Config{
		ProjectID: "",
	})
	storage := cache.New()
	model := gemini.NewModel("gemini-1.5-flash-001", "", "us-central1")
	application := service.NewApplication(storage, history, model.Connections, storage, model)
	http_server.RegisterHttpServer(&application.ChatService)
	grpc_server.RegisterGrpcServer(&application.ChatService, &application.HistoryService, "127.0.0.1", 8081)
	creds, err := credentials.NewClientTLSFromFile("cert/serv.cert", "chat.grpc-server.com")
	if err != nil {
		panic(err)
	}
	conn, err := grpc.NewClient("127.0.0.1:8081", grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}
	grpcClient := grpc_server.NewChatClient(conn)
	resp, err := grpcClient.RetrieveHistory(context.Background(), &grpc_server.HistoryRetrieveRequest{})

	if err != nil {
		panic(err)
	}
	log.Printf("in main %v", resp)
	// resp, err := grpcClient.CreateChat(context.Background(), &grpc_server.ChatRequest{})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(resp.GetChatId())
	// resp2, err := grpcClient.CreateChat(context.Background(), &grpc_server.ChatRequest{})
	// if err != nil {
	// 	panic(err)
	// }
	// iter := 0
	// for {
	// 	log.Printf("we start")
	// 	respp, err := grpcClient.SendMessage(context.Background(), &grpc_server.MessageRequest{ChatId: resp.GetChatId(), Content: fmt.Sprintf("Hello Yabadubi %d", iter)})
	// 	iter++
	// 	if err != nil {
	// 		log.Printf("error during sending message%v", err)
	// 	} else {
	// 		log.Printf("full message is %s %s", respp.GetAuthor(), respp.GetContent())
	// 	}

	// 	<-time.After(time.Millisecond * time.Duration(10000))
	// 	respp2, err := grpcClient.SendMessage(context.Background(), &grpc_server.MessageRequest{ChatId: resp2.GetChatId(), Content: fmt.Sprintf("Hello AAHAHAH %d", iter)})
	// 	iter++
	// 	if err != nil {
	// 		log.Printf("error during sending message%v", err)
	// 	} else {
	// 		log.Printf("full message is %s %s", respp2.GetAuthor(), respp2.GetContent())
	// 	}
	// 	log.Printf("we end")

	// 	<-time.After(time.Millisecond * time.Duration(10000))
	// }
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
