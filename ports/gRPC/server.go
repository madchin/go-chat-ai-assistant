package grpc_server

import (
	"bytes"
	context "context"
	"fmt"
	"net"
	"time"

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
	domainMsg, err := chat.NewMessage(chat.Customer, content)
	if err != nil {
		return &MessageResponse{}, err
	}
	response, err := g.service.SendMessage(chatId, domainMsg)
	if err != nil {
		return &MessageResponse{}, err
	}
	return &MessageResponse{Author: response.Author().Role(), Content: response.Content()}, nil
}

func (g *GrpcServer) SendMessageStream(msg *MessageRequest, stream Chat_SendMessageStreamServer) error {
	chatId := msg.GetChatId()
	content := msg.GetContent()
	domainMsg, err := chat.NewMessage(chat.Customer, content)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(make([]byte, 0, 128))
	go func() {
		for {
			time.Sleep(time.Millisecond * 50)
			stream.Send(&MessageResponse{Author: chat.Assistant.Role(), Content: buffer.String()})
		}
	}()
	err = g.service.SendMessageStream(buffer, chatId, domainMsg)
	if err != nil {
		return err
	}
	return nil
}
