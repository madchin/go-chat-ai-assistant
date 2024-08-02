package http_server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	service "github.com/madchin/go-chat-ai-assistant/services"
)

type route struct {
	r string
}

var (
	chatSend   = route{"/chat/send"}
	chatCreate = route{"/chat/create"}
)

type HttpServer struct {
	chatService *service.ChatService
	router      *httprouter.Router
}

func (h *HttpServer) registerRoutes() {
	h.router.POST(chatSend.r, h.sendMessageHandler)
	h.router.POST(chatCreate.r, h.createChatHandler)
}

func Register(chatService *service.ChatService, host string) {
	server := &HttpServer{chatService, httprouter.New()}
	server.router.PanicHandler = crashHandler
	server.registerRoutes()
	go func() {
		err := http.ListenAndServeTLS(host, "./cert/serv.cert", "./cert/priv.key", registerMiddlewares(server.router))
		if err != nil {
			panic(err.Error())
		}
	}()
}

func (h *HttpServer) createChatHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	chatId := uuid.NewString()
	err := h.chatService.CreateChat(chatId, "")
	if err != nil {
		badRequest(w, serverCodeError.c, err.Error())
		return
	}
	setChatIdCookie(w, chatId)
	w.WriteHeader(http.StatusCreated)
}

func (h *HttpServer) sendMessageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var customerMessage IncomingMessage
	err := json.NewDecoder(r.Body).Decode(&customerMessage)
	if err != nil {
		badRequest(w, clientCodeError.c, "JSON parse failed")
		return
	}
	chatIdCookie, err := r.Cookie("chatId")
	if err != nil || !isValidUUID(chatIdCookie.Value) {
		badRequest(w, clientCodeError.c, err.Error())
		return
	}
	customerDomainMsg, err := customerMessage.toDomainMessage()
	if err != nil {
		badRequest(w, clientCodeError.c, err.Error())
		return
	}
	msg, err := h.chatService.SendMessage(chatIdCookie.Value, customerDomainMsg)
	if err != nil {
		badRequest(w, clientCodeError.c, err.Error())
		return
	}

	assistantMessage := mapDomainMessageToOutcomingMessage(msg)
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
