package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func SetupChatRouter(
	chatService service_logic.IChatService,
	logger *log.Logger,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(CHAT_GET_CLIENT_CHATS, ChatGetClientChatsHandler(chatService, logger))
	router.HandleFunc(CHAT_GET_REPETITOR_CHATS, ChatGetRepetitorChatsHandler(chatService, logger))
	router.HandleFunc(CHAT_GET_MODERATOR_CHATS, ChatGetModeratorChatsHandler(chatService, logger))
	router.HandleFunc(CHAT_START_CM_CHAT, ChatStartCMHandler(chatService, logger))
	router.HandleFunc(CHAT_START_RM_CHAT, ChatStartRMHandler(chatService, logger))
	router.HandleFunc(CHAT_START_CR_CHAT, ChatStartCRHandler(chatService, logger))
	router.HandleFunc(CHAT_GET_CHAT, ChatGetChatHandler(chatService, logger))
	router.HandleFunc(CHAT_SEND_MESSAGE, ChatSendMessageHandler(chatService, logger))
	router.HandleFunc(CHAT_GET_MESSAGES, ChatGetChatMessagesHandler(chatService, logger))
	router.HandleFunc(CHAT_DELETE_CHAT, ChatDeleteChatHandler(chatService, logger))
	router.HandleFunc(CHAT_CLEAR_MESSAGES, ChatClearMessagesHandler(chatService, logger))
	return router
}

func ChatGetClientChatsHandler(chatService service_logic.IChatService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clientIDStr := r.URL.Query().Get("id")
		clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting clientID to int: %v", err)
			http.Error(w, "Invalid clientID", http.StatusBadRequest)
			return
		}
		logger.Printf("Client ID: %v", clientID)
		chatsOffsetStr := r.URL.Query().Get("chats_offset")
		chatsOffset, err := strconv.ParseInt(chatsOffsetStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting chatsOffset to int: %v", err)
			http.Error(w, "Invalid chatsOffset", http.StatusBadRequest)
			return
		}
		logger.Printf("Chats offset: %v", chatsOffset)
		chatsLimitStr := r.URL.Query().Get("chats_limit")
		chatsLimit, err := strconv.ParseInt(chatsLimitStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting chatsLimit to int: %v", err)
			http.Error(w, "Invalid chatsLimit", http.StatusBadRequest)
			return
		}
		logger.Printf("Chats limit: %v", chatsLimit)
		chats, err := chatService.GetChatListByClientID(clientID, chatsOffset, chatsLimit)
		if err != nil {
			logger.Printf("Error getting chats: %v", err)
			http.Error(w, "Error getting chats", http.StatusInternalServerError)
			return
		}
		serverChats := make([]types.ServerChat, 0)
		for _, chat := range chats {
			serverChats = append(serverChats, *types.MapperChatServiceToServer(&chat))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverChats)
		logger.Printf("Chats retrieved: %v", serverChats)
	}
}

func ChatGetRepetitorChatsHandler(chatService service_logic.IChatService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("id")
		repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		logger.Printf("Repetitor ID: %v", repetitorID)
		chatsOffsetStr := r.URL.Query().Get("chats_offset")
		chatsOffset, err := strconv.ParseInt(chatsOffsetStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting chatsOffset to int: %v", err)
			http.Error(w, "Invalid chatsOffset", http.StatusBadRequest)
			return
		}
		logger.Printf("Chats offset: %v", chatsOffset)
		chatsLimitStr := r.URL.Query().Get("chats_limit")
		chatsLimit, err := strconv.ParseInt(chatsLimitStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting chatsLimit to int: %v", err)
			http.Error(w, "Invalid chatsLimit", http.StatusBadRequest)
			return
		}
		logger.Printf("Chats limit: %v", chatsLimit)
		chats, err := chatService.GetChatListByRepetitorID(repetitorID, chatsOffset, chatsLimit)
		if err != nil {
			logger.Printf("Error getting chats: %v", err)
			http.Error(w, "Error getting chats", http.StatusInternalServerError)
			return
		}
		serverChats := make([]types.ServerChat, 0)
		for _, chat := range chats {
			serverChats = append(serverChats, *types.MapperChatServiceToServer(&chat))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverChats)
		logger.Printf("Chats retrieved: %v", serverChats)
	}
}

func ChatGetModeratorChatsHandler(chatService service_logic.IChatService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		moderatorIDStr := r.URL.Query().Get("id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting moderatorID to int: %v", err)
			http.Error(w, "Invalid moderatorID", http.StatusBadRequest)
			return
		}
		logger.Printf("Moderator ID: %v", moderatorID)
		chatsOffsetStr := r.URL.Query().Get("chats_offset")
		chatsOffset, err := strconv.ParseInt(chatsOffsetStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting chatsOffset to int: %v", err)
			http.Error(w, "Invalid chatsOffset", http.StatusBadRequest)
			return
		}
		logger.Printf("Chats offset: %v", chatsOffset)
		chatsLimitStr := r.URL.Query().Get("chats_limit")
		chatsLimit, err := strconv.ParseInt(chatsLimitStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting chatsLimit to int: %v", err)
			http.Error(w, "Invalid chatsLimit", http.StatusBadRequest)
			return
		}
		logger.Printf("Chats limit: %v", chatsLimit)
		chats, err := chatService.GetChatListByModeratorID(moderatorID, chatsOffset, chatsLimit)
		if err != nil {
			logger.Printf("Error getting chats: %v", err)
			http.Error(w, "Error getting chats", http.StatusInternalServerError)
			return
		}
		serverChats := make([]types.ServerChat, 0)
		for _, chat := range chats {
			serverChats = append(serverChats, *types.MapperChatServiceToServer(&chat))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverChats)
		logger.Printf("Chats retrieved: %v", serverChats)
	}
}

func ChatStartCMHandler(chatService service_logic.IChatService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clientIDStr := r.URL.Query().Get("c_id")
		clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting clientID to int: %v", err)
			http.Error(w, "Invalid clientID", http.StatusBadRequest)
			return
		}
		logger.Printf("Client ID: %v", clientID)
		moderatorIDStr := r.URL.Query().Get("m_id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting moderatorID to int: %v", err)
			http.Error(w, "Invalid moderatorID", http.StatusBadRequest)
			return
		}
		logger.Printf("Moderator ID: %v", moderatorID)
		_, err = chatService.CreateCMChat(clientID, moderatorID)
		if err != nil {
			logger.Printf("Error creating chat: %v", err)
			http.Error(w, "Error creating chat", http.StatusInternalServerError)
			return
		}
		logger.Printf("Chat created")
		chatID, err := chatService.GetChatIdByCIDAndMID(clientID, moderatorID)
		if err != nil {
			logger.Printf("Error getting chat ID: %v", err)
			http.Error(w, "Error getting chat ID", http.StatusInternalServerError)
			return
		}
		logger.Printf("Chat ID: %v", chatID)
		json.NewEncoder(w).Encode(chatID)
	}
}

func ChatStartRMHandler(chatService service_logic.IChatService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("r_id")
		repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		logger.Printf("Repetitor ID: %v", repetitorID)
		moderatorIDStr := r.URL.Query().Get("m_id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting moderatorID to int: %v", err)
			http.Error(w, "Invalid moderatorID", http.StatusBadRequest)
			return
		}
		logger.Printf("Moderator ID: %v", moderatorID)
		_, err = chatService.CreateRMChat(repetitorID, moderatorID)
		if err != nil {
			logger.Printf("Error creating chat: %v", err)
			http.Error(w, "Error creating chat", http.StatusInternalServerError)
			return
		}
		logger.Printf("Chat created")
		chatID, err := chatService.GetChatIdByMIDAndRID(moderatorID, repetitorID)
		if err != nil {
			logger.Printf("Error getting chat ID: %v", err)
			http.Error(w, "Error getting chat ID", http.StatusInternalServerError)
			return
		}
		logger.Printf("Chat ID: %v", chatID)
		json.NewEncoder(w).Encode(chatID)
	}
}

func ChatStartCRHandler(chatService service_logic.IChatService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clientIDStr := r.URL.Query().Get("c_id")
		clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting clientID to int: %v", err)
			http.Error(w, "Invalid clientID", http.StatusBadRequest)
			return
		}
		logger.Printf("Client ID: %v", clientID)
		repetitorIDStr := r.URL.Query().Get("r_id")
		repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		logger.Printf("Repetitor ID: %v", repetitorID)
		_, err = chatService.CreateCRChat(clientID, repetitorID)
		if err != nil {
			logger.Printf("Error creating chat: %v", err)
			http.Error(w, "Error creating chat", http.StatusInternalServerError)
			return
		}
		logger.Printf("Chat created")
		chatID, err := chatService.GetChatIdByCIDAndRID(clientID, repetitorID)
		if err != nil {
			logger.Printf("Error getting chat ID: %v", err)
			http.Error(w, "Error getting chat ID", http.StatusInternalServerError)
			return
		}
		logger.Printf("Chat ID: %v", chatID)
		json.NewEncoder(w).Encode(chatID)
	}
}

func ChatGetChatHandler(chatService service_logic.IChatService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		chatIDStr := r.URL.Query().Get("id")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting chatID to int: %v", err)
			http.Error(w, "Invalid chatID", http.StatusBadRequest)
			return
		}
		logger.Printf("Chat ID: %v", chatID)
		chat, err := chatService.GetChat(chatID)
		if err != nil {
			logger.Printf("Error getting chat: %v", err)
			http.Error(w, "Error getting chat", http.StatusInternalServerError)
			return
		}
		serverChat := types.MapperChatServiceToServer(chat)
		json.NewEncoder(w).Encode(serverChat)
		logger.Printf("Chat retrieved: %v", serverChat)
	}
}

func ChatSendMessageHandler(chatService service_logic.IChatService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "PATCH" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		message := ""
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			logger.Printf("Error decoding message: %v", err)
			http.Error(w, "Error decoding message", http.StatusBadRequest)
			return
		}
		logger.Printf("Message: %v", message)
		senderIDStr := r.URL.Query().Get("sender_id")
		senderID, err := strconv.ParseInt(senderIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting senderID to int: %v", err)
			http.Error(w, "Invalid senderID", http.StatusBadRequest)
			return
		}
		logger.Printf("Sender ID: %v", senderID)
		chatIDStr := r.URL.Query().Get("chat_id")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting chatID to int: %v", err)
			http.Error(w, "Invalid chatID", http.StatusBadRequest)
			return
		}
		logger.Printf("Chat ID: %v", chatID)
		err = chatService.SendMessage(chatID, senderID, message)
		if err != nil {
			logger.Printf("Error sending message: %v", err)
			http.Error(w, "Error sending message", http.StatusInternalServerError)
			return
		}
		logger.Printf("Message sent")
		w.WriteHeader(http.StatusOK)
	}
}

func ChatGetChatMessagesHandler(chatService service_logic.IChatService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		chatIDStr := r.URL.Query().Get("id")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting chatID to int: %v", err)
			http.Error(w, "Invalid chatID", http.StatusBadRequest)
			return
		}
		logger.Printf("Chat ID: %v", chatID)
		messagesOffsetStr := r.URL.Query().Get("messages_offset")
		messagesOffset, err := strconv.ParseInt(messagesOffsetStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting messagesOffset to int: %v", err)
			http.Error(w, "Invalid messagesOffset", http.StatusBadRequest)
			return
		}
		logger.Printf("Messages offset: %v", messagesOffset)
		messagesLimitStr := r.URL.Query().Get("messages_limit")
		messagesLimit, err := strconv.ParseInt(messagesLimitStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting messagesLimit to int: %v", err)
			http.Error(w, "Invalid messagesLimit", http.StatusBadRequest)
			return
		}
		logger.Printf("Messages limit: %v", messagesLimit)
		messages, err := chatService.GetMessages(chatID, messagesOffset, messagesLimit)
		if err != nil {
			logger.Printf("Error getting messages: %v", err)
			http.Error(w, "Error getting messages", http.StatusInternalServerError)
			return
		}
		serverMessages := make([]types.ServerMessage, 0)
		for _, message := range messages {
			serverMessages = append(serverMessages, *types.MapperMessageServiceToServer(&message))
		}
		logger.Printf("Messages retrieved: %v", serverMessages)
		json.NewEncoder(w).Encode(serverMessages)
	}
}

func ChatDeleteChatHandler(chatService service_logic.IChatService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "DELETE" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		chatIDStr := r.URL.Query().Get("id")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting chatID to int: %v", err)
			http.Error(w, "Invalid chatID", http.StatusBadRequest)
			return
		}
		logger.Printf("Chat ID: %v", chatID)
		err = chatService.DeleteChat(chatID)
		if err != nil {
			logger.Printf("Error deleting chat: %v", err)
			http.Error(w, "Error deleting chat", http.StatusInternalServerError)
			return
		}
		logger.Printf("Chat deleted")
		w.WriteHeader(http.StatusOK)
	}
}

func ChatClearMessagesHandler(chatService service_logic.IChatService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "PUT" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		chatIDStr := r.URL.Query().Get("id")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting chatID to int: %v", err)
			http.Error(w, "Invalid chatID", http.StatusBadRequest)
			return
		}
		logger.Printf("Chat ID: %v", chatID)
		err = chatService.DeleteChat(chatID)
		if err != nil {
			logger.Printf("Error clearing messages: %v", err)
			http.Error(w, "Error clearing messages", http.StatusInternalServerError)
			return
		}
		logger.Printf("Messages cleared")
		w.WriteHeader(http.StatusOK)
	}
}
