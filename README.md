## Chat AI Assistant
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
3. Requirements:
    List of requirements using BDD pattern:
    1. Feature: Consumer can chat with predefined "trained" model
    2. Feature: Client can retrieve history of all chats.
    3. Feature: 
