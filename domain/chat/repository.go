package chat

type ChatMessages map[string][]Message

type StorageRepository interface {
	CreateChat(id, context string) error
	SendMessage(chatId string, msg Message) error
}

type HistoryRepository interface {
	HistoryRetriever
	HistorySaver
}

type AssistantService interface {
	Ask(content, chatId string) (Message, error)
	AskStream(response chan<- string, content, chatId string) (Message, error)
}

type HistoryRetriever interface {
	RetrieveAllChatsHistory() (ChatMessages, error)
	RetrieveAllChatsHistoryStream(partialResponseCh chan<- ChatMessages) error
}

type HistorySaver interface {
	SaveHistory(ChatMessages) error
}
