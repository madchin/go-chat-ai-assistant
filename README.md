## Chat AI Assistant

## System Architecture

![Architecture](/.docs/architecture.png)

## Release 0.3.0

### Added tests for port/grpc and port/http


---


## Release 0.2.0

### Integration with gemini ai model

### Possibility to communicate with ai model via gRPC framework and HTTP protocol

### All messages sent from customer (client-side) to assistant(ai-response) are cached in-memory. Consider it as "chat"

### Messages within "chat" become outdated after certain time, and batch-saved to history repository in storage_cleaner service which can be referred as worker (runs continously in separate goroutine on specified interval). If save fails for any of the message, it is retried certain times (with respect to firebase delay for each document which is 1 second)

### Each connection to gemini model (its based on RPC) is also cached and reused for specific chat, and removed periodically like messages within "chat"

### Each chat bases on queue data structure, and for each chat message limit is established

### Each message sent is validated only for client-side message (Customer)


---


## First planning 

1. Introduction
    This README provides all information about the Chat AI Assistant module, its requirements, future features, basic use cases.
    1. Terms:

    * Consumer -> Entity which communicate with server via client

    * Client -> Entity which manages consumers and request resources from server 
    
    * Server -> Entity which handles responses to the client, and forwarding requests to the AI model

    * Token -> Unit used in AI model for counting message price


2. Proof of concept
    Related chapter is about proofing basic concept which is client chat with AI model with server as mediator.
    What AI model to use?
    * Chat GPT?
    * Host model locally - resource intensive
    * Mistral AI - requires sign in with phone for free trial
    * models via HuggingFace - llama-3 chatbot, no integration to golang, availble for js / py
    * google cloud ai - free credits for start

    For POC we gonna use google cloud ai, not multimodal, only text based.
    [gcloudcli](https://cloud.google.com/vertex-ai/generative-ai/docs/start/quickstarts/quickstart-multimodal#local-shell)
    Initialize gcloud cli and install necessary deps, then using base code we gonna create simple example of request to the ai model.
    **DONE** 
3. AI model - Possibilities / Considerations
    1. Input / output price model (we pay 2x for one message)
    2. Context caching - cache model (75% cost reduction in comparison to input data)
    - cache previous messages to preserve context?
    - TTL of cached messages? e.g. browser session + timeout based on last message? what if cache is flushed and customer again will start interacting?
    3. Function calling - model respond with predefined function names and its params based on func description, based on question customer made, then we can make external api call for e.g. and provide response to model and he gives us response based on our data
    4. Grounding - ability to ground response to predefined data (google search / own data)
4. Data storage - Chat history
    1. Sql / NoSQL?
    2. Redis + Firestore?
5. Vertex AI / Google AI?
    1. comparison:
    https://cloud.google.com/vertex-ai/generative-ai/docs/migrate/migrate-google-ai



4. Requirements:
    List of requirements using BDD pattern:
    1. Feature: Consumer can chat with AI model with provided context
    - Stories
        - As a client I want be able to specify domain so that bot can reply with domain related answers
        ```
            Given: As a client
            When: I will want to specify bot context
            Then: It will be preserved for every customer chat

            Given: As a client
            When: I will specify bot context
            Then: It will be saved in bot cache again when cache expiry and user send message to save resources
        ```
        - As a customer I want be able to communicate so that i can get information about specified subject
        ```
            Given: As a customer
            When: I will want to ask chat assistant with specified question about predefined domain
            Then: I will receive correct answer about asked subject from ai model

            Given: As a customer
            When: I will ask question not related to predefined domain
            Then: I will receive answer about wrong query and information to correct message to related one to our domain from ai model

            Given: As a customer
            When: I will send message which token count is less than 2
            Then: I will receive answer about too short query and information to correct message from server

            Given: As a customer
            When: I will send message which token count is greater than rate limit 
            (1 mln TPM / 4 mln TPM) -- Google AI -- Vertex AI - 4mln TPM (1.5 flash)
            Then: I will receive answer about too long query and information to correct message from server

            Given: As a customer
            When: I will extend request rate limit 
            (15 RPM / 1000 RPM) Google AI -- Vertex AI -- 200 RPM (1.5 flash) 
            Then: I will receive message about slowing down from server
        ```
        - As a customer I wont be able to send more than specified amount of messages in specified timeframe
        ```            
            Given: As a customer
            When: I will send 5 messages in 15 seconds timeframe
            Then: On next message, I will receive answer from server about slowing down and will be delayed to send next message
            until (1st message timestamp - now > 15)
        ```
    2. Feature: Client can retrieve history of all chats.
    - Stories
        - As a client I want have ability to retrieve history of all clients chats
        ```
            Given: As a customer
            When: My last message timestamp - now > 120 seconds
            Then: Chat will be flushed from in-memory storage and saved in persistent place

            Given: As a client
            When: I will send request to retrieve history of all clients chats
            Then: Chats from in-memory storage will be saved to persistent storage and marked as saved, then chats from database will be retrieved
        ```
5. Project:
- Adapters
    - http
        - Firebase
    - grpc
        - GeminiAdapter {AiModel}
    - Repository
        - Firebase
            - SendMessage(chatId string, role Participant, message Message) error
            - RetrieveConversation(chatId string) ([]Message, error)
            - SaveChatHistory(chat Chat) error
            - RetrieveAllChats() ([]Chat, error)
        - In-memory chat storage
            - SendMessage(chatId string, role Participant, message Message) error
            - RetrieveConversation(chatId string) ([]Message, error)
- Application 
    - Interface definition - AiModel 
    - firebase ChatHistoryRepository
    - redis ChatRepository
    - aiModel AiModel
    - ChatService
        - (a Application) SendMessage(chatId string, role Participant, message Message) error
        - (a Application) RetrieveConversation(chatId string) ([]Message, error)
    - ChatHistoryService
        - (a Application) Retrieve() ([]Chat, error)
        - (a Application) Save(chat Chat) error
- Core
    - Chat
        - Chat_id
        - Context
        - MessageQueue
        - Participant
            - Customer
            - Assistant
        - Message
            - Content
            - Timestamp
            - (m Message) Send(role Participant, onSendMessage(msg Message))
        - MessageValidator
    - Repository
        - Chat
            - SendMessage(chatId string, role Participant, message Message) error
            - RetrieveConversation(chatId string) ([]Message, error)
        - Chat History
            - RetrieveAllChats() ([]Chat, error)
            - SaveChatHistory(chat Chat) error
- Ports
    - grpc
        - SendMessage(customerId string, content string)
        - RetrieveConversation(customerId string) []Message (ServerStream)
        - RetrieveChatsHistory() []Chat (ServerStream)
    - http
        - SendMessage(customerId string, content string)
        - RetrieveConversation(customerId string) []Message
        - RetrieveChatsHistory() []Chat
    - ws
        - SendMessage(customerId string, content string)
        - RetrieveConversation(customerId string) []Message
        - RetrieveChatsHistory() []Chat
