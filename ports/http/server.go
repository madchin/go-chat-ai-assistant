package http_server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	inmemory_storage "github.com/madchin/go-chat-ai-assistant/adapters/repository/in_memory"
	"github.com/madchin/go-chat-ai-assistant/ports"
)

type HttpServer struct {
	chatService ports.ChatService
	router      *httprouter.Router
}

func RegisterHttpServer(chatService ports.ChatService) {
	server := &HttpServer{chatService, httprouter.New()}
	server.router.POST("/chat/send", server.sendMessageHandler)
	server.router.POST("/chat/create", server.createChatHandler)
	server.router.PanicHandler = crashHandler
	go func() {
		err := http.ListenAndServe(":8080", registerMiddlewares(server.router))
		if err != nil {
			panic(err.Error())
		}
	}()
}

func (h *HttpServer) createChatHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	chatId := uuid.NewString()
	err := h.chatService.CreateChat(chatId, "")
	if err != nil && err != inmemory_storage.ErrChatAlreadyExists {
		badRequest(w, "server", err.Error())
		return
	}
	setChatIdCookie(w, chatId)
	w.WriteHeader(http.StatusCreated)
}

func (h *HttpServer) sendMessageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var customerMessage IncomingMessage
	err := json.NewDecoder(r.Body).Decode(&customerMessage)
	if err != nil {
		badRequest(w, "client", "JSON parse failed")
		return
	}
	chatIdCookie, err := r.Cookie("chatId")
	if err != nil || !isValidUUID(chatIdCookie.Value) {
		badRequest(w, "client", err.Error())
		return
	}
	customerDomainMsg, err := customerMessage.toDomainMessage()
	if err != nil {
		badRequest(w, "client", err.Error())
		return
	}
	msg, err := h.chatService.SendMessage(chatIdCookie.Value, customerDomainMsg)
	if err != nil {
		badRequest(w, "client", err.Error())
		return
	}

	assistantMessage := OutcomingMessage{Author: msg.Author().Role(), Content: msg.Content()}
	err = json.NewEncoder(w).Encode(&assistantMessage)
	if err != nil {
		badRequest(w, "internal", err.Error())
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

func setChatIdCookie(w http.ResponseWriter, chatId string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "chatId",
		Value:    chatId,
		Path:     "/chat",
		MaxAge:   240,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	return cookie
}

func isValidUUID(uid string) bool {
	_, err := uuid.Parse(uid)
	return err == nil
}

func crashHandler(w http.ResponseWriter, _ *http.Request, _ any) {
	internal(w)
}
