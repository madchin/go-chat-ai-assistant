package chat

type OnRemoveMessage func(Message) error
type UserMessages map[string][]Message

type Repository interface {
	CreateChat(id, context string) error
	SendMessage(chatId string, msg Message) error
	RemoveOutdatedMessages() UserMessages
	RetrieveAllConversations() (UserMessages, error)
}

type HistoryRepository interface {
	SaveHistory(UserMessages) error
	RetrieveAllChatsHistory() (UserMessages, error)
}
