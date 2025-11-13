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

func SetupRepetitorRouterV2(repetitorService service_logic.IRepetitorService, contractService service_logic.IContractService, transactionService service_logic.ITransactionService, resumeService service_logic.IResumeService) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(EXACT_REPETITOR_V2, RepetitorGetHandlerV2(repetitorService)).Methods("GET")
	router.HandleFunc(EXACT_REPETITOR_V2, RepetitorAssignContractHandlerV2(contractService)).Methods("PATCH")
	return router
}

func RepetitorGetHandlerV2(repetitorService service_logic.IRepetitorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repetitorID := mux.Vars(r)["repetitorId"]
		repetitorIDInt, err := strconv.Atoi(repetitorID)
		if err != nil {
			http.Error(w, "Invalid repetitor ID", http.StatusBadRequest)
			return
		}
		repetitor, err := repetitorService.GetRepetitorProfile(int64(repetitorIDInt), 0, 0)
		if err != nil {
			http.Error(w, "Repetitor not found", http.StatusNotFound)
			return
		}
		serverRepetitor := types.MapperRepetitorProfileServiceToServerV2(repetitor)
		json.NewEncoder(w).Encode(serverRepetitor)
		w.WriteHeader(http.StatusOK)
	}
}

func RepetitorAssignContractHandlerV2(contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repetitorIDStr := mux.Vars(r)["repetitorId"]
		repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid repetitor ID", http.StatusBadRequest)
			return
		}
		contractIDStr := r.URL.Query().Get("contract_id")
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil || contractID <= 0 {
			http.Error(w, "Invalid contract_id", http.StatusBadRequest)
			return
		}
		// Ensure contract exists first
		if _, err := contractService.GetContract(contractID); err != nil {
			http.Error(w, "Contract not found", http.StatusNotFound)
			return
		}
		if err := contractService.UpdateContractRepetitorID(contractID, repetitorID); err != nil {
			http.Error(w, "Error updating contract repetitor", http.StatusBadRequest)
			return
		}
		contract, err := contractService.GetContract(contractID)
		if err != nil {
			http.Error(w, "Contract not found", http.StatusNotFound)
			return
		}
		resp := types.MapperContractServiceToServerV2(contract)
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusOK)
	}
}

func SetupRepetitorRouter(
	repetitorService service_logic.IRepetitorService,
	contractService service_logic.IContractService,
	transactionService service_logic.ITransactionService,
	resumeService service_logic.IResumeService,
	logger *log.Logger,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(REPETITOR_GET_PROFILE, RepetitorGetProfileHandler(repetitorService, logger))
	router.HandleFunc(REPETITOR_GET_CONTRACTS, RepetitorGetContractsHandler(contractService, logger))
	router.HandleFunc(REPETITOR_GET_AVAILABLE_CONTRACTS, RepetitorGetAvailableContractsHandler(contractService, logger))
	router.HandleFunc(REPETITOR_ACCEPT_CONTRACT, RepetitorAcceptContractHandler(repetitorService, contractService, logger))
	router.HandleFunc(REPETITOR_MAKE_REVIEW, RepetitorMakeReviewHandler(contractService, logger))
	router.HandleFunc(REPETITOR_PAY_FOR_CONTRACT, RepetitorPayForContractHandler(repetitorService, contractService, transactionService, logger))
	router.HandleFunc(REPETITOR_CANCEL_CONTRACT, RepetitorCancelContractHandler(repetitorService, contractService, logger))
	router.HandleFunc(REPETITOR_COMPLETE_CONTRACT, RepetitorCompleteContractHandler(repetitorService, contractService, logger))
	router.HandleFunc(REPETITOR_CHANGE_RESUME, RepetitorChangeResumeHandler(repetitorService, resumeService, logger))
	return router
}

func RepetitorGetProfileHandler(repetitorService service_logic.IRepetitorService, logger *log.Logger) http.HandlerFunc {
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
		reviewsOffsetStr := r.URL.Query().Get("reviews_offset")
		reviewsOffset, err := strconv.ParseInt(reviewsOffsetStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting reviewsOffset to int: %v", err)
			http.Error(w, "Invalid reviewsOffset", http.StatusBadRequest)
			return
		}
		logger.Printf("Reviews Offset: %v", reviewsOffset)
		reviewsLimitStr := r.URL.Query().Get("reviews_limit")
		reviewsLimit, err := strconv.ParseInt(reviewsLimitStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting reviewsLimit to int: %v", err)
			http.Error(w, "Invalid reviewsLimit", http.StatusBadRequest)
			return
		}
		logger.Printf("Reviews Limit: %v", reviewsLimit)
		repetitor, err := repetitorService.GetRepetitorProfile(int64(repetitorID), int64(reviewsOffset), int64(reviewsLimit))
		if err != nil {
			logger.Printf("Error getting repetitor: %v", err)
			http.Error(w, "Error getting repetitor", http.StatusBadRequest)
			return
		}
		logger.Printf("Repetitor retrieved: %v", repetitor)
		serverRepetitor := types.MapperRepetitorProfileServiceToServer(repetitor)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverRepetitor)
	}
}

func RepetitorGetContractsHandler(contractService service_logic.IContractService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("repetitor_id")
		repetitorID, err := strconv.Atoi(repetitorIDStr)
		if err != nil {
			logger.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		logger.Printf("Repetitor ID: %v", repetitorID)
		offsetStr := r.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			logger.Printf("Error converting offset to int: %v", err)
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
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
		contracts, err := contractService.GetRepetitorContractList(int64(repetitorID), int64(offset), int64(limit), types.ContractStatus(status))
		if err != nil {
			logger.Printf("Error getting contracts: %v", err)
			http.Error(w, "Error getting contracts", http.StatusBadRequest)
			return
		}
		logger.Printf("Contracts retrieved: %v", contracts)
		serverContracts := make([]types.ServerContract, len(contracts))
		for i, contract := range contracts {
			serverContracts[i] = *types.MapperContractServiceToServer(&contract)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverContracts)
	}
}

func RepetitorGetAvailableContractsHandler(
	contractService service_logic.IContractService,
	logger *log.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
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
		contracts, err := contractService.GetContractList(int64(offset), int64(limit), types.ContractStatus(status))
		if err != nil {
			logger.Printf("Error getting contracts: %v", err)
			http.Error(w, "Error getting contracts", http.StatusBadRequest)
			return
		}
		logger.Printf("Contracts retrieved: %v", contracts)
		serverContracts := make([]types.ServerContract, len(contracts))
		for i, contract := range contracts {
			serverContracts[i] = *types.MapperContractServiceToServer(&contract)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(serverContracts)
	}
}

func RepetitorAcceptContractHandler(
	repetitorService service_logic.IRepetitorService,
	contractService service_logic.IContractService,
	logger *log.Logger,
) http.HandlerFunc {
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
		repetitorIDStr := r.URL.Query().Get("repetitor_id")
		repetitorID, err := strconv.Atoi(repetitorIDStr)
		if err != nil {
			logger.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		logger.Printf("Repetitor ID: %v", repetitorID)
		_, err = repetitorService.GetRepetitorData(int64(repetitorID))
		if err != nil {
			logger.Printf("Error getting repetitor: %v", err)
			http.Error(w, "Error getting repetitor", http.StatusBadRequest)
			return
		}
		logger.Printf("Repetitor ID: %v", repetitorID)
		err = contractService.UpdateContractRepetitorID(int64(contractID), int64(repetitorID))
		if err != nil {
			logger.Printf("Error updating contract repetitor: %v", err)
			http.Error(w, "Error updating contract repetitor", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract %d updated with repetitor %d", contractID, repetitorID)
		w.WriteHeader(http.StatusOK)
	}
}

func RepetitorMakeReviewHandler(
	contractService service_logic.IContractService,
	logger *log.Logger,
) http.HandlerFunc {
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
		logger.Printf("Request body: %s", string(body))
		review := types.ServerReview{}
		if err := json.Unmarshal(body, &review); err != nil {
			logger.Printf("Error unmarshalling request body: %v", err)
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		logger.Printf("Review: %v", review)
		serviceReview := types.MapperReviewServerToService(&review)
		_, err = contractService.CreateContractReviewRepetitor(int64(contractID), *serviceReview)
		if err != nil {
			logger.Printf("Error creating review: %v", err)
			http.Error(w, "Error creating review", http.StatusBadRequest)
			return
		}
		logger.Printf("Review created: %v", review)
		w.WriteHeader(http.StatusOK)
	}
}

func RepetitorPayForContractHandler(
	repetitorService service_logic.IRepetitorService,
	contractService service_logic.IContractService,
	transactionService service_logic.ITransactionService,
	logger *log.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		userIDStr := r.URL.Query().Get("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			logger.Printf("Error converting userID to int: %v", err)
			http.Error(w, "Invalid userID", http.StatusBadRequest)
			return
		}
		logger.Printf("User ID: %v", userID)
		contractIDStr := r.URL.Query().Get("contract_id")
		contractID, err := strconv.Atoi(contractIDStr)
		if err != nil {
			logger.Printf("Error converting contractID to int: %v", err)
			http.Error(w, "Invalid contractID", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract ID: %v", contractID)
		amountStr := r.URL.Query().Get("amount")
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			logger.Printf("Error converting amount to int: %v", err)
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}
		logger.Printf("Amount: %v", amount)
		_, err = contractService.GetContract(int64(contractID))
		if err != nil {
			logger.Printf("Error getting contract: %v", err)
			http.Error(w, "Error getting contract", http.StatusBadRequest)
			return
		}
		transactionID, err := transactionService.CreateContractPaymentTransaction(int64(amount), int64(userID), int64(contractID))
		if err != nil {
			logger.Printf("Error creating transaction: %v", err)
			http.Error(w, "Error creating transaction", http.StatusBadRequest)
			return
		}
		logger.Printf("Transaction created: %v", transactionID)
		err = contractService.UpdateContractPaymentStatus(int64(contractID), types.PaymentStatusPaid)
		if err != nil {
			logger.Printf("Error updating contract status: %v", err)
			http.Error(w, "Error updating contract status", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract status updated")
		w.WriteHeader(http.StatusOK)
	}
}

func RepetitorCancelContractHandler(
	repetitorService service_logic.IRepetitorService,
	contractService service_logic.IContractService,
	logger *log.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("id")
		repetitorID, err := strconv.Atoi(repetitorIDStr)
		if err != nil {
			logger.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		logger.Printf("Repetitor ID: %v", repetitorID)
		contractIDStr := r.URL.Query().Get("c_id")
		contractID, err := strconv.Atoi(contractIDStr)
		if err != nil {
			logger.Printf("Error converting contractID to int: %v", err)
			http.Error(w, "Invalid contractID", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract ID: %v", contractID)
		contract, err := contractService.GetContract(int64(contractID))
		if err != nil {
			logger.Printf("Error getting contract: %v", err)
			http.Error(w, "Error getting contract", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract: %v", contract)
		if contract.RepetitorID != int64(repetitorID) {
			logger.Printf("Contract doesn't belong to repetitor")
			http.Error(w, "Contract doesn't belong to repetitor", http.StatusNotFound)
			return
		}
		if contract.Status != types.ContractStatusActive {
			logger.Printf("Contract is not active")
			http.Error(w, "Contract is not active", http.StatusBadRequest)
			return
		}
		err = contractService.UpdateContractStatus(int64(contractID), types.ContractStatusCancelled)
		if err != nil {
			logger.Printf("Error updating contract status: %v", err)
			http.Error(w, "Error updating contract status", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract status updated")
		w.WriteHeader(http.StatusOK)
	}
}

func RepetitorCompleteContractHandler(
	repetitorService service_logic.IRepetitorService,
	contractService service_logic.IContractService,
	logger *log.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		repetitorIDStr := r.URL.Query().Get("id")
		repetitorID, err := strconv.Atoi(repetitorIDStr)
		if err != nil {
			logger.Printf("Error converting repetitorID to int: %v", err)
			http.Error(w, "Invalid repetitorID", http.StatusBadRequest)
			return
		}
		logger.Printf("Repetitor ID: %v", repetitorID)
		contractIDStr := r.URL.Query().Get("c_id")
		contractID, err := strconv.Atoi(contractIDStr)
		if err != nil {
			logger.Printf("Error converting contractID to int: %v", err)
			http.Error(w, "Invalid contractID", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract ID: %v", contractID)
		contract, err := contractService.GetContract(int64(contractID))
		if err != nil {
			logger.Printf("Error getting contract: %v", err)
			http.Error(w, "Error getting contract", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract: %v", contract)
		if contract.RepetitorID != int64(repetitorID) {
			logger.Printf("Contract doesn't belong to repetitor")
			http.Error(w, "Contract doesn't belong to repetitor", http.StatusNotFound)
			return
		}
		err = contractService.UpdateContractStatus(int64(contractID), types.ContractStatusCompleted)
		if err != nil {
			logger.Printf("Error updating contract status: %v", err)
			http.Error(w, "Error updating contract status", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract status updated")
		w.WriteHeader(http.StatusOK)
	}
}

func RepetitorChangeResumeHandler(
	repetitorService service_logic.IRepetitorService,
	resumeService service_logic.IResumeService,
	logger *log.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		logger.Printf("Request body: %s", string(body))
		repetitor, err := repetitorService.GetRepetitorData(repetitorID)
		if err != nil {
			logger.Printf("Error getting repetitor: %v", err)
			http.Error(w, "Error getting repetitor", http.StatusBadRequest)
			return
		}
		logger.Printf("Request body: %s", string(body))
		resume := types.ServerResume{}
		if err := json.Unmarshal(body, &resume); err != nil {
			logger.Printf("Error unmarshalling request body: %v", err)
			http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
			return
		}
		logger.Printf("Resume: %v", resume)
		err = resumeService.UpdateResumeTitle(repetitor.ResumeID, resume.Title)
		if err != nil {
			logger.Printf("Error updating repetitor resume: %v", err)
			http.Error(w, "Error updating repetitor resume", http.StatusBadRequest)
			return
		}
		err = resumeService.UpdateResumeDescription(repetitor.ResumeID, resume.Description)
		if err != nil {
			logger.Printf("Error updating repetitor resume: %v", err)
			http.Error(w, "Error updating repetitor resume", http.StatusBadRequest)
			return
		}
		logger.Printf("Resume title updated")
		err = resumeService.UpdateResumePrices(repetitor.ResumeID, resume.Prices)
		if err != nil {
			logger.Printf("Error updating repetitor resume: %v", err)
			http.Error(w, "Error updating repetitor resume", http.StatusBadRequest)
			return
		}
		logger.Printf("Resume prices updated")
		logger.Printf("Resume updated")
		w.WriteHeader(http.StatusOK)
	}
}
