package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func SetupChatRouterV2(chatService service_logic.IChatService) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(CHATS_V2, ChatGetChatsHandlerV2(chatService)).Methods("GET")
	router.HandleFunc(CHATS_V2, ChatCreateChatHandlerV2(chatService)).Methods("POST")
	router.HandleFunc(EXACT_CHAT_V2, ChatGetChatHandlerV2(chatService)).Methods("GET")
	router.HandleFunc(EXACT_CHAT_V2, ChatUpdateChatHandlerV2(chatService)).Methods("PATCH")
	router.HandleFunc(EXACT_CHAT_V2, ChatDeleteChatHandlerV2(chatService)).Methods("DELETE")
	router.HandleFunc(EXACT_CHAT_V2, ChatClearChatHandlerV2(chatService)).Methods("PUT")
	router.HandleFunc(EXACT_CHAT_MESSAGES_V2, ChatGetMessagesHandlerV2(chatService)).Methods("GET")
	router.HandleFunc(EXACT_CHAT_MESSAGES_V2, ChatSendMessageHandlerV2(chatService)).Methods("POST")
	return router
}

func ChatGetChatsHandlerV2(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		clientID := queryParams.Get("userId")
		clientIDInt, err := strconv.Atoi(clientID)
		if err != nil {
			http.Error(w, "Invalid client ID", http.StatusBadRequest)
			return
		}
		chatsOffset := queryParams.Get("offset")
		chatsOffsetInt, err := strconv.Atoi(chatsOffset)
		if err != nil {
			http.Error(w, "Invalid chats offset", http.StatusBadRequest)
			return
		}
		chatsLimit := queryParams.Get("limit")
		chatsLimitInt, err := strconv.Atoi(chatsLimit)
		if err != nil {
			http.Error(w, "Invalid chats limit", http.StatusBadRequest)
			return
		}
		chats, err := chatService.GetChatListByClientID(int64(clientIDInt), int64(chatsOffsetInt), int64(chatsLimitInt))
		if err != nil {
			http.Error(w, "Error getting chats", http.StatusBadRequest)
			return
		}
		serverChats := make([]types.ServerChat, 0)
		for _, chat := range chats {
			serverChats = append(serverChats, *types.MapperChatServiceToServer(&chat))
		}
		json.NewEncoder(w).Encode(serverChats)
		w.WriteHeader(http.StatusOK)
	}
}

func ChatCreateChatHandlerV2(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		var chat types.ServerChatV2
		if err := json.Unmarshal(body, &chat); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		if chat.Type != "client_moderator" && chat.Type != "repetitor_moderator" && chat.Type != "client_repetitor" {
			http.Error(w, "Invalid chat type", http.StatusBadRequest)
			return
		}
		var chatID int64
		if chat.Type == "client_moderator" {
			chatID, err = chatService.CreateCMChat(chat.ClientID, chat.ModeratorID)
			if err != nil {
				http.Error(w, "Error creating chat", http.StatusBadRequest)
				return
			}
		}
		if chat.Type == "repetitor_moderator" {
			chatID, err = chatService.CreateRMChat(chat.RepetitorID, chat.ModeratorID)
			if err != nil {
				http.Error(w, "Error creating chat", http.StatusBadRequest)
				return
			}
		}
		if chat.Type == "client_repetitor" {
			chatID, err = chatService.CreateCRChat(chat.ClientID, chat.RepetitorID)
			if err != nil {
				http.Error(w, "Error creating chat", http.StatusBadRequest)
				return
			}
		}
		createdChat, err := chatService.GetChat(int64(chatID))
		if err != nil {
			http.Error(w, "Error getting chat", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.MapperChatServiceToServerV2(createdChat))
	}
}

func ChatUpdateChatHandlerV2(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		var req types.ServerChatUpdateV2
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		if req.Status != "active" && req.Status != "closed" {
			http.Error(w, "Invalid chat status", http.StatusBadRequest)
			return
		}
		chatID := mux.Vars(r)["chatId"]
		chatIDInt, err := strconv.Atoi(chatID)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		err = chatService.UpdateChat(int64(chatIDInt), req.Status)
		if err != nil {
			http.Error(w, ERR_MSG_CHAT_NOT_FOUND, http.StatusNotFound)
			return
		}
		updatedChat, err := chatService.GetChat(int64(chatIDInt))
		if err != nil {
			http.Error(w, ERR_MSG_CHAT_NOT_FOUND, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(types.MapperChatServiceToServerV2(updatedChat))
	}
}

func ChatGetChatHandlerV2(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID := mux.Vars(r)["chatId"]
		chatIDInt, err := strconv.Atoi(chatID)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		chat, err := chatService.GetChat(int64(chatIDInt))
		if err != nil {
			http.Error(w, ERR_MSG_CHAT_NOT_FOUND, http.StatusNotFound)
			return
		}
		serverChat := types.MapperChatServiceToServerV2(chat)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverChat)
	}
}

func ChatGetMessagesHandlerV2(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID := mux.Vars(r)["chatId"]
		chatIDInt, err := strconv.Atoi(chatID)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		offset := r.URL.Query().Get("offset")
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		limit := r.URL.Query().Get("limit")
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		messages, err := chatService.GetMessages(int64(chatIDInt), int64(offsetInt), int64(limitInt))
		if err != nil {
			http.Error(w, "Error getting messages", http.StatusBadRequest)
			return
		}
		serverMessages := make([]types.ServerMessageV2, 0)
		for _, message := range messages {
			serverMessages = append(serverMessages, *types.MapperMessageServiceToServerV2(&message))
		}
		json.NewEncoder(w).Encode(serverMessages)
		w.WriteHeader(http.StatusOK)
	}
}

func ChatSendMessageHandlerV2(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ServerMessageCreateV2
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		chatID := mux.Vars(r)["chatId"]
		chatIDInt, err := strconv.Atoi(chatID)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}

		message, err := chatService.SendMessage(int64(chatIDInt), req.SenderID, req.Content)
		if err != nil {
			http.Error(w, ERR_MSG_CHAT_NOT_FOUND, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	}
}

func ChatDeleteChatHandlerV2(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID := mux.Vars(r)["chatId"]
		chatIDInt, err := strconv.Atoi(chatID)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		err = chatService.DeleteChat(int64(chatIDInt))
		if err != nil {
			http.Error(w, ERR_MSG_CHAT_NOT_FOUND, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func ChatClearChatHandlerV2(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID := mux.Vars(r)["chatId"]
		chatIDInt, err := strconv.Atoi(chatID)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		err = chatService.ClearChat(int64(chatIDInt))
		if err != nil {
			http.Error(w, "Error clearing chat", http.StatusBadRequest)
			return
		}
		updatedChat, err := chatService.GetChat(int64(chatIDInt))
		if err != nil {
			http.Error(w, ERR_MSG_CHAT_NOT_FOUND, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(types.MapperChatServiceToServerV2(updatedChat))
	}
}

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
			http.Error(w, "Error getting chats", http.StatusBadRequest)
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
			http.Error(w, "Error getting chats", http.StatusBadRequest)
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
			http.Error(w, "Error getting chats", http.StatusBadRequest)
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
			http.Error(w, "Error creating chat", http.StatusBadRequest)
			return
		}
		logger.Printf("Chat created")
		chatID, err := chatService.GetChatIdByCIDAndMID(clientID, moderatorID)
		if err != nil {
			logger.Printf("Error getting chat ID: %v", err)
			http.Error(w, "Error getting chat ID", http.StatusBadRequest)
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
			http.Error(w, "Error creating chat", http.StatusBadRequest)
			return
		}
		logger.Printf("Chat created")
		chatID, err := chatService.GetChatIdByMIDAndRID(moderatorID, repetitorID)
		if err != nil {
			logger.Printf("Error getting chat ID: %v", err)
			http.Error(w, "Error getting chat ID", http.StatusBadRequest)
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
			http.Error(w, "Error creating chat", http.StatusBadRequest)
			return
		}
		logger.Printf("Chat created")
		chatID, err := chatService.GetChatIdByCIDAndRID(clientID, repetitorID)
		if err != nil {
			logger.Printf("Error getting chat ID: %v", err)
			http.Error(w, "Error getting chat ID", http.StatusBadRequest)
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
			http.Error(w, "Error getting chat", http.StatusBadRequest)
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
		_, err = chatService.SendMessage(chatID, senderID, message)
		if err != nil {
			logger.Printf("Error sending message: %v", err)
			http.Error(w, "Error sending message", http.StatusBadRequest)
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
			http.Error(w, "Error getting messages", http.StatusBadRequest)
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
			http.Error(w, "Error deleting chat", http.StatusBadRequest)
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
			http.Error(w, "Error clearing messages", http.StatusBadRequest)
			return
		}
		logger.Printf("Messages cleared")
		w.WriteHeader(http.StatusOK)
	}
}
