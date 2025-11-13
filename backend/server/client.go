package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gorilla/mux"
)

func SetupClientRouterV2(clientService service_logic.IClientService, contractService service_logic.IContractService) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(EXACT_CLIENT_V2, ClientGetHandlerV2(clientService)).Methods("GET")
	return router
}

func ClientGetHandlerV2(clientService service_logic.IClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := mux.Vars(r)["clientId"]
		fmt.Printf("client ID: %s\n", clientID)
		clientIDInt, err := strconv.Atoi(clientID)
		if err != nil {
			http.Error(w, "Invalid client ID", http.StatusBadRequest)
			return
		}
		client, err := clientService.GetClientProfile(int64(clientIDInt), 0, 0)
		if err != nil {
			http.Error(w, "Client not found", http.StatusNotFound)
			return
		}
		serverClient := types.MapperClientProfileServiceToServerV2(client)
		json.NewEncoder(w).Encode(serverClient)
		w.WriteHeader(http.StatusOK)
	}
}

func SetupClientRouter(
	clientService service_logic.IClientService,
	contractService service_logic.IContractService,
	logger *log.Logger,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(CLIENT_GET_PROFILE, ClientGetProfileHandler(clientService, logger))
	router.HandleFunc(CLIENT_CREATE_CONTRACT, ClientCreateContractHandler(contractService, logger))
	router.HandleFunc(CLIENT_GET_CONTRACTS, ClientGetContractsHandler(contractService, logger))
	router.HandleFunc(CLIENT_MAKE_REVIEW, ClientMakeReviewHandler(contractService, logger))
	return router
}

func ClientGetProfileHandler(clientService service_logic.IClientService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clientIDStr := r.URL.Query().Get("id")
		clientID, err := strconv.Atoi(clientIDStr)
		if err != nil {
			logger.Printf("Error converting clientID to int: %v", err)
			http.Error(w, "Invalid clientID", http.StatusBadRequest)
			return
		}
		logger.Printf("Client ID: %v", clientID)
		reviewsOffsetStr := r.URL.Query().Get("reviews_offset")
		reviewsOffset, err := strconv.Atoi(reviewsOffsetStr)
		if err != nil {
			logger.Printf("Error converting reviewsOffset to int: %v", err)
			http.Error(w, "Invalid reviewsOffset", http.StatusBadRequest)
			return
		}
		logger.Printf("Reviews offset: %v", reviewsOffset)
		reviewsLimitStr := r.URL.Query().Get("reviews_limit")
		reviewsLimit, err := strconv.Atoi(reviewsLimitStr)
		if err != nil {
			logger.Printf("Error converting reviewsLimit to int: %v", err)
			http.Error(w, "Invalid reviewsLimit", http.StatusBadRequest)
			return
		}
		logger.Printf("Reviews limit: %v", reviewsLimit)
		client, err := clientService.GetClientProfile(int64(clientID), int64(reviewsOffset), int64(reviewsLimit))
		if err != nil {
			http.Error(w, "Error getting client", http.StatusBadRequest)
			return
		}
		serverClient := types.MapperClientProfileServiceToServer(client)
		logger.Printf("Client retrieved: %v", serverClient)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverClient)
	}
}

func ClientCreateContractHandler(contractService service_logic.IContractService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		logger.Printf("Request body: %v", string(body))
		contractInitInfo := types.ServerContractInitData{}
		if err := json.Unmarshal(body, &contractInitInfo); err != nil {
			logger.Printf("Error unmarshalling request body: %v", err)
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract init info: %v", contractInitInfo)
		serviceContractInitInfo := types.ServiceContractInitData(contractInitInfo)
		contractID, err := contractService.CreateContract(serviceContractInitInfo)
		if err != nil {
			logger.Printf("Error creating contract: %v", err)
			http.Error(w, "Error creating contract", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract created with ID: %d", contractID)
		json.NewEncoder(w).Encode(contractID)
		w.WriteHeader(http.StatusCreated)
	}
}

func ClientGetContractsHandler(contractService service_logic.IContractService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clientIDStr := r.URL.Query().Get("client_id")
		clientID, err := strconv.Atoi(clientIDStr)
		if err != nil {
			logger.Printf("Error converting clientID to int: %v", err)
			http.Error(w, "Invalid clientID", http.StatusBadRequest)
			return
		}
		logger.Printf("Client ID: %v", clientID)
		offsetStr := r.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			logger.Printf("Error converting offset to int: %v", err)
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		logger.Printf("Offset: %v", offset)
		limitStr := r.URL.Query().Get("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			logger.Printf("Error converting limit to int: %v", err)
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		logger.Printf("Limit: %v", limit)
		statusStr := r.URL.Query().Get("status")
		status, err := strconv.Atoi(statusStr)
		if err != nil {
			logger.Printf("Error converting status to int: %v", err)
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}
		logger.Printf("Status: %v", status)
		contracts, err := contractService.GetClientContractList(int64(clientID), int64(offset), int64(limit), types.ContractStatus(status))
		if err != nil {
			logger.Printf("Error getting contracts: %v", err)
			http.Error(w, "Error getting contracts", http.StatusBadRequest)
			return
		}
		serverContracts := make([]types.ServerContract, len(contracts))
		for i, contract := range contracts {
			serverContracts[i] = *types.MapperContractServiceToServer(&contract)
		}
		logger.Printf("Contracts retrieved: %v", serverContracts)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverContracts)
	}
}

func ClientMakeReviewHandler(contractService service_logic.IContractService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		contractIDStr := r.URL.Query().Get("contract_id")
		contractID, err := strconv.Atoi(contractIDStr)
		if err != nil {
			logger.Printf("Error converting contractID to int: %v", err)
			http.Error(w, "Invalid contractID", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract ID: %v", contractID)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		review := types.ServerReview{}
		if err := json.Unmarshal(body, &review); err != nil {
			logger.Printf("Error unmarshalling request body: %v", err)
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		review.CreatedAt = time.Now()
		logger.Printf("Review: %v", review)
		serviceReview := types.MapperReviewServerToService(&review)
		_, err = contractService.CreateContractReviewClient(int64(contractID), *serviceReview)
		if err != nil {
			logger.Printf("Error making review: %v", err)
			http.Error(w, "Error making review", http.StatusBadRequest)
			return
		}
		logger.Printf("Review made: %v", review)
		json.NewEncoder(w).Encode(review)
		w.WriteHeader(http.StatusOK)
	}
}
