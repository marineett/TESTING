package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func SetupModeratorRouter(
	transactionService service_logic.ITransactionService,
	contractService service_logic.IContractService,
	moderatorService service_logic.IModeratorService,
	logger *log.Logger,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(MODERATOR_GET_PROFILE, ModeratorGetProfileHandler(moderatorService, logger))
	router.HandleFunc(MODERATOR_GET_TRANSACTION_TO_APPROVE, ModeratorGetTransactionsToApproveHandler(transactionService, logger))
	router.HandleFunc(MODERATOR_APPROVE_TRANSACTION, ModeratorApproveTransactionHandler(transactionService, logger))
	router.HandleFunc(MODERATOR_GET_CONTRACTS, ModeratorGetContractsHandler(contractService, logger))
	router.HandleFunc(MODERATOR_BAN_CONTRACT, ModeratorBanContractHandler(contractService, logger))
	return router
}

func ModeratorGetProfileHandler(moderatorService service_logic.IModeratorService, logger *log.Logger) http.HandlerFunc {
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
			logger.Printf("Invalid moderator ID: %v", err)
			http.Error(w, "Invalid moderator ID", http.StatusBadRequest)
			return
		}
		logger.Printf("Moderator ID: %v", moderatorID)
		moderator, err := moderatorService.GetModeratorProfile(moderatorID)
		if err != nil {
			logger.Printf("Error getting moderator: %v", err)
			http.Error(w, "Error getting moderator", http.StatusInternalServerError)
			return
		}
		serverModerator := types.MapperModeratorProfileServiceToServer(moderator)
		logger.Printf("Moderator retrieved: %v", serverModerator)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverModerator)
	}
}

func ModeratorGetTransactionsToApproveHandler(transactionService service_logic.ITransactionService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		transaction, err := transactionService.GetPendingContractPaymentTransaction()
		if err != nil {
			logger.Printf("Error getting transaction: %v", err)
			http.Error(w, "Error getting transaction", http.StatusInternalServerError)
			return
		}
		serverTransaction := types.ServerPendingContractPaymentTransaction(*transaction)
		logger.Printf("Transaction retrieved: %v", serverTransaction)
		json.NewEncoder(w).Encode(serverTransaction)
		w.WriteHeader(http.StatusOK)
	}
}

func ModeratorApproveTransactionHandler(transactionService service_logic.ITransactionService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		transactionIDStr := r.URL.Query().Get("id")
		if transactionIDStr == "" {
			logger.Printf("Transaction ID is required")
			http.Error(w, "Transaction ID is required", http.StatusBadRequest)
			return
		}
		logger.Printf("Transaction ID: %v", transactionIDStr)
		transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid transaction ID: %v", err)
			http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
			return
		}
		logger.Printf("Transaction ID: %v", transactionID)
		err = transactionService.ApproveTransaction(transactionID)
		if err != nil {
			logger.Printf("Error approving transaction: %v", err)
			http.Error(w, "Error approving transaction", http.StatusInternalServerError)
			return
		}
		logger.Printf("Transaction approved")
		w.WriteHeader(http.StatusOK)
	}
}

func ModeratorGetContractsHandler(contractService service_logic.IContractService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fromStr := r.URL.Query().Get("from")
		if fromStr == "" {
			logger.Printf("From is required")
			http.Error(w, "From is required", http.StatusBadRequest)
			return
		}
		logger.Printf("From: %v", fromStr)
		from, err := strconv.ParseInt(fromStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid from: %v", err)
			http.Error(w, "Invalid from", http.StatusBadRequest)
			return
		}
		logger.Printf("From: %v", from)
		sizeStr := r.URL.Query().Get("size")
		if sizeStr == "" {
			logger.Printf("Size is required")
			http.Error(w, "Size is required", http.StatusBadRequest)
			return
		}
		logger.Printf("Size: %v", sizeStr)
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid size: %v", err)
			http.Error(w, "Invalid size", http.StatusBadRequest)
			return
		}
		logger.Printf("Size: %v", size)
		contracts, err := contractService.GetAllContracts(from, size)
		if err != nil {
			logger.Printf("Error getting contracts: %v", err)
			http.Error(w, "Error getting contracts", http.StatusInternalServerError)
			return
		}
		serverContracts := make([]types.ServerContract, len(contracts))
		for i, contract := range contracts {
			serverContracts[i] = *types.MapperContractServiceToServer(&contract)
		}
		logger.Printf("Contracts retrieved: %v", serverContracts)
		json.NewEncoder(w).Encode(serverContracts)
		w.WriteHeader(http.StatusOK)
	}
}

func ModeratorBanContractHandler(contractService service_logic.IContractService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		contractIDStr := r.URL.Query().Get("id")
		if contractIDStr == "" {
			logger.Printf("Contract ID is required")
			http.Error(w, "Contract ID is required", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract ID: %v", contractIDStr)
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid contract ID: %v", err)
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract ID: %v", contractID)
		contract, err := contractService.GetContract(contractID)
		if err != nil {
			logger.Printf("Error getting contract: %v", err)
			http.Error(w, "Error getting contract", http.StatusInternalServerError)
			return
		}
		logger.Printf("Contract retrieved: %v", contract)
		if contract.Status != types.ContractStatusActive && contract.Status != types.ContractStatusPending {
			logger.Printf("Contract is not in valid status")
			http.Error(w, "Contract is not in valid status", http.StatusBadRequest)
			return
		}
		err = contractService.UpdateContractStatus(contractID, types.ContractStatusBanned)
		if err != nil {
			logger.Printf("Error banning contract: %v", err)
			http.Error(w, "Error banning contract", http.StatusInternalServerError)
			return
		}
		logger.Printf("Contract banned")
		w.WriteHeader(http.StatusOK)
	}
}
