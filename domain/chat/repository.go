package chat

type Repository interface {
	SendMessage(role participant, msg Message, onSendMessage func(Message) error)
	RemoveMessage(onRemoveMessage func(Message) (Message, error)) (Message, error)
	RetrieveAllConversations() ([]Chat, error)
}

type HistoryRepository interface {
	SaveHistory(chat Chat) error
	RetrieveAllChatsHistory() ([]Chat, error)
}
