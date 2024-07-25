package grpc_server

import (
	context "context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
	"github.com/madchin/go-chat-ai-assistant/ports"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	service ports.ChatService
}

func RegisterGrpcServer(service ports.ChatService, host string, port int) {
	chatServer := &GrpcServer{service}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err.Error())
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	RegisterChatServer(grpcServer, chatServer)
	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			panic(err.Error())
		}
	}()
}

func (g *GrpcServer) mustEmbedUnimplementedChatServer() {}

func (g *GrpcServer) CreateChat(ctx context.Context, chat *ChatRequest) (*ChatResponse, error) {
	chatId := uuid.NewString()
	err := g.service.CreateChat(chatId, "")
	if err != nil {
		return &ChatResponse{}, err
	}
	return &ChatResponse{ChatId: chatId}, nil
}

func (g *GrpcServer) SendMessage(ctx context.Context, msg *MessageRequest) (*MessageResponse, error) {
	chatId := msg.GetChatId()
	content := msg.GetContent()
	customerMsg, err := chat.NewCustomerMessage(content)
	if err != nil {
		return &MessageResponse{}, err
	}
	response, err := g.service.SendMessage(chatId, customerMsg)
	if err != nil {
		return &MessageResponse{}, err
	}
	return &MessageResponse{Author: response.Author().Role(), Content: response.Content()}, nil
}

func (g *GrpcServer) SendMessageStream(msg *MessageRequest, stream Chat_SendMessageStreamServer) error {
	chatId := msg.GetChatId()
	content := msg.GetContent()
	customerMsg, err := chat.NewCustomerMessage(content)
	if err != nil {
		return err
	}
	assistantResponseCh := make(chan string)
	errCh := make(chan error)
	go streamResponsePart(assistantResponseCh, errCh, stream)
	err = g.service.SendMessageStream(assistantResponseCh, chatId, customerMsg)
	if err != nil {
		return err
	}
	if err = <-errCh; err != nil {
		return err
	}
	return nil
}

func streamResponsePart(partialResponseCh <-chan string, errCh chan<- error, stream Chat_SendMessageStreamServer) {
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
