package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/madchin/go-chat-ai-assistant/domain/chat"
)

const (
	chatEndpoint = "/chat"
)

type HttpServer struct {
	ms     MessageService
	router *httprouter.Router
}

type MessageService interface {
	CreateChat(id, context string) error
	SendMessage(chatId string, msg chat.Message) (chat.Message, error)
}

func NewHttpServer(messageService MessageService) *HttpServer {
	server := &HttpServer{messageService, httprouter.New()}
	server.router.POST(chatEndpoint, server.chatHandler)
	server.router.PanicHandler = crashHandler
	http.ListenAndServe(":8080", registerMiddlewares(server.router))
	return server
}

func (h HttpServer) chatHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var customerMessage IncomingMessage
	err := json.NewDecoder(r.Body).Decode(&customerMessage)
	if err != nil {
		badRequest(w, "client", "JSON parse failed")
		return
	}
	chatId := chatIdFromCookie(w)
	setChatInIdCookie(w, chatId)

	//FIXME context
	err = h.ms.CreateChat(chatId, "")
	if err != nil {
		badRequest(w, "client", err.Error())
		return
	}

	msg, err := h.ms.SendMessage(chatId, customerMessage.toDomainMessage())
	if err != nil {
		badRequest(w, "client", err.Error())
		return
	}

	assistantMessage := OutcomingMessage{Author: msg.Author().Role(), Content: msg.Content()}
	err = json.NewEncoder(w).Encode(&assistantMessage)
	if err != nil {
		internal(w)
	}

}

func badRequest(w http.ResponseWriter, code, description string) {
	httpErr := HttpError{Code: code, Description: description}
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(&httpErr)
}

func internal(w http.ResponseWriter) {
	httpErr := HttpError{Code: "server", Description: "Oops! Something went wrong"}
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(&httpErr)
}

func setChatInIdCookie(w http.ResponseWriter, chatId string) {
	cookie := fmt.Sprintf("chatId=%s; Max-Age=240; Secure; Path=/chat; HttpOnly; SameSite=Strict", chatId)
	w.Header().Add("Set-Cookie", cookie)
}

func chatIdFromCookie(w http.ResponseWriter) string {
	chatId := w.Header().Get("Cookie")
	re := regexp.MustCompile(`chatId=(\w+)`)
	match := re.FindStringSubmatch(chatId)
	if len(match) > 1 {
		chatId = match[1]
	} else {
		chatId = uuid.NewString()
	}
	return chatId
}

func crashHandler(w http.ResponseWriter, _ *http.Request, _ any) {
	internal(w)
}
