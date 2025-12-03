package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func SetupMessageRouterV2(chatService service_logic.IChatService) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(EXACT_MESSAGE_V2, UpdateMessageContentHandlerV2(chatService)).Methods("PATCH")
	r.HandleFunc(EXACT_MESSAGE_V2, DeleteMessageHandlerV2(chatService)).Methods("DELETE")
	return r
}

func UpdateMessageContentHandlerV2(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messageIDStr := mux.Vars(r)["messageId"]
		messageID, err := strconv.Atoi(messageIDStr)
		if err != nil {
			http.Error(w, "Invalid message ID", http.StatusBadRequest)
			return
		}
		var req types.ServerMessageContentPatchV2
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, ERR_MSG_INVALID_BODY, http.StatusBadRequest)
			return
		}
		err = chatService.UpdateMessageContent(int64(messageID), req.Content)
		if err != nil {
			http.Error(w, ERR_MSG_MESSAGE_NOT_FOUND, http.StatusNotFound)
			return
		}
		updated, err := chatService.GetMessage(int64(messageID))
		if err != nil {
			http.Error(w, ERR_MSG_MESSAGE_NOT_FOUND, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(types.MapperMessageServiceToServerV2(updated))
	}
}

func DeleteMessageHandlerV2(chatService service_logic.IChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messageIDStr := mux.Vars(r)["messageId"]
		messageID, err := strconv.Atoi(messageIDStr)
		if err != nil {
			http.Error(w, "Invalid message ID", http.StatusBadRequest)
			return
		}
		if err := chatService.DeleteMessage(int64(messageID)); err != nil {
			http.Error(w, ERR_MSG_MESSAGE_NOT_FOUND, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
