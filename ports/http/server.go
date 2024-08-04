package http_server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/madchin/go-chat-ai-assistant/ports"
)

type route struct {
	r string
}

var (
	chatSend    = route{"/chat/send"}
	chatCreate  = route{"/chat/create"}
	chatHistory = route{"/chat/history"}
)

type HttpServer struct {
	chatService    ports.ChatService
	historyService ports.HistoryService
	router         *httprouter.Router
}

func (h *HttpServer) registerRoutes() {
	h.router.POST(chatSend.r, h.sendMessageHandler)
	h.router.POST(chatCreate.r, h.createChatHandler)
	h.router.GET(chatHistory.r, h.retrieveHistoryHandler)
}

func Register(chatService ports.ChatService, historyService ports.HistoryService, host string) {
	server := &HttpServer{chatService, historyService, httprouter.New()}
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
	cookie := chatIdCookie(chatId)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)
}

func (h *HttpServer) sendMessageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Header.Get("Content-Type") != "application/json" {
		badRequest(w, clientCodeError.c, wrongContentTypeError.d)
		return
	}
	var customerMessage CustomerMessage
	err := json.NewDecoder(r.Body).Decode(&customerMessage)
	if err != nil {
		badRequest(w, clientCodeError.c, bodyParseError.d)
		return
	}
	chatIdCookie, err := r.Cookie("chatId")
	if err != nil {
		badRequest(w, clientCodeError.c, err.Error())
		return
	}
	if !isValidUUID(chatIdCookie.Value) {
		badRequest(w, clientCodeError.c, wrongCookieValue.d)
		return
	}
	msg, err := h.chatService.SendMessage(chatIdCookie.Value, customerMessage.Content)
	if err != nil {
		badRequest(w, clientCodeError.c, err.Error())
		return
	}

	assistantMessage := mapDomainMessageToHttpAssistantMessage(msg)
	err = json.NewEncoder(w).Encode(&assistantMessage)
	if err != nil {
		internal(w)
	}

}

func (h *HttpServer) retrieveHistoryHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	history, err := h.historyService.RetrieveAllChatsHistory()
	if err != nil {
		badRequest(w, serverCodeError.c, err.Error())
		return
	}
	httpHistory := mapDomainChatMessagesToHttpChatsHistory(history)
	err = json.NewEncoder(w).Encode(httpHistory)
	if err != nil {
		badRequest(w, serverCodeError.c, bodyParseError.d)
	}
}

func badRequest(w http.ResponseWriter, code, description string) {
	httpErr := HttpError{Code: code, Description: description}
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(&httpErr)
}

func internal(w http.ResponseWriter) {
	httpErr := HttpError{Code: serverCodeError.c, Description: genericDescriptionError.d}
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(&httpErr)
}

func chatIdCookie(chatId string) *http.Cookie {
	return &http.Cookie{
		Name:     "chatId",
		Value:    chatId,
		Path:     "/chat",
		MaxAge:   240,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func isValidUUID(uid string) bool {
	_, err := uuid.Parse(uid)
	return err == nil
}

func crashHandler(w http.ResponseWriter, _ *http.Request, _ any) {
	internal(w)
}
