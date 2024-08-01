package chat

type ChatMessages map[string][]Message

type Repository interface {
	CreateChat(id, context string) error
	SendMessage(chatId string, msg Message) error
	RetrieveAllConversations() (ChatMessages, error)
}

type HistoryRepository interface {
	SaveHistory(ChatMessages) error
	RetrieveAllChatsHistory(responseCh chan<- ChatMessages) error
}
