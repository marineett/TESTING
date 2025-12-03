package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func SetupGuestRouter(
	repetitorService service_logic.IRepetitorService,
	logger *log.Logger,
) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc(GUEST_GET_REPETITORS, GuestGetRepetitorsHandler(repetitorService, logger))

	return router
}

func GuestGetRepetitorsHandler(repetitorService service_logic.IRepetitorService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorsOffsetStr := r.URL.Query().Get("repetitors_offset")
		repetitorsOffset, err := strconv.ParseInt(repetitorsOffsetStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting offset to int: %v", err)
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		logger.Printf("Repetitors offset: %v", repetitorsOffset)
		repetitorsLimitStr := r.URL.Query().Get("repetitors_limit")
		repetitorsLimit, err := strconv.ParseInt(repetitorsLimitStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting limit to int: %v", err)
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		logger.Printf("Repetitors limit: %v", repetitorsLimit)
		repetitors, err := repetitorService.GetRepetitors(repetitorsOffset, repetitorsLimit)
		if err != nil {
			logger.Printf("Error getting repetitors: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Printf("Repetitors retrieved: %v", repetitors)
		serverRepetitors := make([]types.ServerRepetitorView, len(repetitors))
		for i, repetitor := range repetitors {
			serverRepetitors[i] = *types.MapperRepetitorViewServiceToServer(repetitor)
		}
		json.NewEncoder(w).Encode(serverRepetitors)
		w.WriteHeader(http.StatusOK)
	}
}
