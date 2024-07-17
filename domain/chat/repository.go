package chat

type OnRemoveMessage func(Message) error
type UserMessages map[string][]Message

type Repository interface {
	SendMessage(chatId string, msg Message) error
	RemoveMessage(chatId string, onRemoveMessage OnRemoveMessage) (Message, error)
	RetrieveAllConversations() (UserMessages, error)
}

type HistoryRepository interface {
	SaveHistory(id string) error
	RetrieveAllChatsHistory() (UserMessages, error)
}
