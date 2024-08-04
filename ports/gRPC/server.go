package grpc_server

import (
	context "context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
	service "github.com/madchin/go-chat-ai-assistant/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcServer struct {
	chatService            *service.ChatService
	historyRetrieveService *service.HistoryRetrieveService
}

func Register(chatService *service.ChatService, historyRetrieveService *service.HistoryRetrieveService, host string, port int) {
	chatServer := &GrpcServer{chatService, historyRetrieveService}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err.Error())
	}
	creds, err := credentials.NewServerTLSFromFile("./cert/serv.cert", "./cert/priv.key")
	if err != nil {
		panic(err.Error())
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	RegisterChatServer(grpcServer, chatServer)
	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			panic(err.Error())
		}
	}()

}

func (g *GrpcServer) mustEmbedUnimplementedChatServer() {}

func (g *GrpcServer) RetrieveHistory(_ *HistoryRetrieveRequest, stream Chat_RetrieveHistoryServer) error {
	errCh := make(chan error)
	partialResponseCh := make(chan chat.ChatMessages)
	go streamRetrieveHistoryPartialContent(partialResponseCh, errCh, stream)
	err := g.historyRetrieveService.RetrieveAllChatsHistory(partialResponseCh)
	if err != nil {
		return err
	}
	if err = <-errCh; err != nil {
		return err
	}
	return nil
}

func (g *GrpcServer) CreateChat(ctx context.Context, chat *ChatRequest) (*ChatResponse, error) {
	chatId := uuid.NewString()
	err := g.chatService.CreateChat(chatId, "")
	if err != nil {
		return &ChatResponse{}, err
	}

	return &ChatResponse{ChatId: chatId}, nil
}

func (g *GrpcServer) SendMessage(ctx context.Context, msg *MessageRequest) (*MessageResponse, error) {
	chatId := msg.GetChatId()
	content := msg.GetContent()
	response, err := g.chatService.SendMessage(chatId, content)
	if err != nil {
		return &MessageResponse{}, err
	}
	return &MessageResponse{Author: response.Author().Role(), Content: response.Content()}, nil
}

func (g *GrpcServer) SendMessageStream(msg *MessageRequest, stream Chat_SendMessageStreamServer) error {
	chatId := msg.GetChatId()
	content := msg.GetContent()
	assistantResponseCh := make(chan string)
	errCh := make(chan error)
	go streamSendMessagePartialContent(assistantResponseCh, errCh, stream)
	err := g.chatService.SendMessageStream(assistantResponseCh, chatId, content)
	if err != nil {
		return err
	}
	if err = <-errCh; err != nil {
		return err
	}

	return nil
}

func streamSendMessagePartialContent(partialResponseCh <-chan string, errCh chan<- error, stream Chat_SendMessageStreamServer) {
	for {
		responsePart, waitForNextPart := <-partialResponseCh
		if !waitForNextPart {
			errCh <- nil
			return
		}
		if err := stream.Send(&MessageResponse{Author: chat.Assistant.Role(), Content: responsePart}); err != nil {
			errCh <- err
			close(errCh)
			return
		}
	}
}

func streamRetrieveHistoryPartialContent(partialResponseCh <-chan chat.ChatMessages, errCh chan<- error, stream Chat_RetrieveHistoryServer) {
	for {
		responsePart, waitForNextPart := <-partialResponseCh
		if !waitForNextPart {
			errCh <- nil
			return
		}
		if err := stream.Send(mapDomainChatMessagesToHistoryResponse(responsePart)); err != nil {
			errCh <- err
			close(errCh)
			return
		}
	}

}

func mapDomainChatMessagesToHistoryResponse(domainChatMessages chat.ChatMessages) *HistoryRetrieveResponse {
	chatHistory := make([]*HistoryRetrieveResponseChat, 0, len(domainChatMessages))
	for chatId, domainMessages := range domainChatMessages {
		messages := make([]*HistoryRetrieveResponseChatMsg, len(domainMessages))
		for i := 0; i < len(domainMessages); i++ {
			author, content, timestamp := domainMessages[i].Author().Role(), domainMessages[i].Content(), domainMessages[i].Timestamp()
			message := &HistoryRetrieveResponseChatMsg{Author: author, Content: content, Timestamp: timestamp}
			messages[i] = message
		}
		chatMessages := &HistoryRetrieveResponseChat{ChatId: chatId, Messages: messages}
		chatHistory = append(chatHistory, chatMessages)
	}
	return &HistoryRetrieveResponse{Chats: chatHistory}
}
