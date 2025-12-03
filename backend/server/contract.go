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

func SetupContractRouterV2(
	contractService service_logic.IContractService,
	lessonService service_logic.ILessonService,
	reviewService service_logic.IReviewService,
	transactionService service_logic.ITransactionService,
) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(CONTRACTS_V2, ContractsListHandlerV2(contractService)).Methods("GET")
	r.HandleFunc(CONTRACTS_V2, ContractCreateHandlerV2(contractService)).Methods("POST")
	r.HandleFunc(EXACT_CONTRACT_V2, ContractGetHandlerV2(contractService)).Methods("GET")
	r.HandleFunc(EXACT_CONTRACT_V2, ContractStatusPatchHandlerV2(contractService)).Methods("PATCH")
	r.HandleFunc(CONTRACT_LESSONS_V2, ContractLessonsListHandlerV2(lessonService)).Methods("GET")
	r.HandleFunc(CONTRACT_LESSONS_V2, ContractLessonCreateHandlerV2(lessonService)).Methods("POST")
	r.HandleFunc(CONTRACT_REVIEWS_V2, ContractReviewsListHandlerV2(reviewService, contractService)).Methods("GET")
	r.HandleFunc(CONTRACT_REVIEWS_V2, ContractReviewCreateHandlerV2(reviewService, contractService)).Methods("POST")
	r.HandleFunc(CONTRACT_TRANSACTIONS_V2, ContractTransactionsListHandlerV2(transactionService, contractService)).Methods("GET")
	r.HandleFunc(CONTRACT_TRANSACTIONS_V2, ContractTransactionCreateHandlerV2(transactionService)).Methods("POST")
	r.HandleFunc(TRANSACTION_APPROVAL_V2, TransactionApproveHandlerV2(transactionService)).Methods("PATCH")
	return r
}

func ContractsListHandlerV2(contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Getting contracts")
		offsetStr := r.URL.Query().Get("offset")
		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		sizeStr := r.URL.Query().Get("limit")
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid size", http.StatusBadRequest)
			return
		}
		clientIDStr := r.URL.Query().Get("client_id")
		clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid client ID", http.StatusBadRequest)
			return
		}
		repetitorIDStr := r.URL.Query().Get("repetitor_id")
		repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid repetitor ID", http.StatusBadRequest)
			return
		}
		fmt.Println("Getting contracts for client ID: ", clientID, " and repetitor ID: ", repetitorID)
		contracts, err := contractService.GetContracts(clientID, repetitorID, offset, size)
		if err != nil {
			http.Error(w, "Error getting contracts", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(contracts)
		w.WriteHeader(http.StatusOK)
	}
}

func ContractCreateHandlerV2(contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ServerContractCreateV2
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		contractID, err := contractService.CreateContract(*types.MapperContractCreateV2ServerToServiceInit(&req))
		if err != nil {
			fmt.Println("Error creating contract: ", err)
			http.Error(w, "Error creating contract", http.StatusBadRequest)
			return
		}
		contract, err := contractService.GetContract(contractID)
		if err != nil {
			http.Error(w, "Error getting contract", http.StatusBadRequest)
			return
		}
		serverContract := types.MapperContractServiceToServerV2(contract)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(serverContract)
	}
}

func ContractGetHandlerV2(contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contractIDStr := mux.Vars(r)["contractId"]
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}
		contract, err := contractService.GetContract(contractID)
		if err != nil {
			http.Error(w, "Contract not found", http.StatusNotFound)
			return
		}
		serverContract := types.MapperContractServiceToServerV2(contract)
		json.NewEncoder(w).Encode(serverContract)
		w.WriteHeader(http.StatusOK)
	}
}

func ContractStatusPatchHandlerV2(contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ServerContractStatusPatchV2
		contractID := mux.Vars(r)["contractId"]
		contractIDInt, err := strconv.ParseInt(contractID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}
		// Existence check
		if _, err := contractService.GetContract(contractIDInt); err != nil {
			http.Error(w, "Contract not found", http.StatusNotFound)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		statusEnum, err := types.ParseContractStatus(req.Status)
		if err != nil {
			http.Error(w, "Invalid contract status", http.StatusBadRequest)
			return
		}
		err = contractService.UpdateContractStatus(contractIDInt, statusEnum)
		if err != nil {
			http.Error(w, "Error updating contract status", http.StatusBadRequest)
			return
		}
		contract, err := contractService.GetContract(contractIDInt)
		if err != nil {
			http.Error(w, "Contract not found", http.StatusNotFound)
			return
		}
		serverContract := types.MapperContractServiceToServerV2(contract)
		json.NewEncoder(w).Encode(serverContract)
		w.WriteHeader(http.StatusOK)
	}
}

func ContractLessonsListHandlerV2(lessonService service_logic.ILessonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contractIDStr := mux.Vars(r)["contractId"]
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}
		offsetStr := r.URL.Query().Get("offset")
		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		sizeStr := r.URL.Query().Get("size")
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid size", http.StatusBadRequest)
			return
		}
		lessons, err := lessonService.GetLessons(contractID, offset, size)
		if err != nil {
			http.Error(w, "Error getting lessons", http.StatusBadRequest)
			return
		}
		serverLessons := make([]types.ServerLessonV2, len(lessons))
		for i, lesson := range lessons {
			serverLessons[i] = *types.MapperLessonServiceToServerV2(&lesson)
		}
		json.NewEncoder(w).Encode(serverLessons)
		w.WriteHeader(http.StatusOK)
	}
}

func ContractLessonCreateHandlerV2(lessonService service_logic.ILessonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contractIDStr := mux.Vars(r)["contractId"]
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}
		var req types.ServerLessonCreateV2
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		lessonID, err := lessonService.CreateLesson(*types.MapperLessonCreateV2ServerToService(contractID, &req))
		if err != nil {
			fmt.Println("Error creating lesson: ", err)
			http.Error(w, "Error creating lesson", http.StatusBadRequest)
			return
		}
		lesson, err := lessonService.GetLesson(lessonID)
		if err != nil {
			fmt.Println("Error getting lesson: ", err)
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}
		serverLesson := types.MapperLessonServiceToServerV2(lesson)
		json.NewEncoder(w).Encode(*serverLesson)
		w.WriteHeader(http.StatusCreated)
	}
}

func LessonGetHandlerV2(lessonService service_logic.ILessonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lessonIDStr := mux.Vars(r)["lessonId"]
		lessonID, err := strconv.ParseInt(lessonIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
			return
		}
		lesson, err := lessonService.GetLesson(lessonID)
		if err != nil {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}
		serverLesson := types.MapperLessonServiceToServerV2(lesson)
		json.NewEncoder(w).Encode(serverLesson)
		w.WriteHeader(http.StatusOK)
	}
}

func LessonPatchHandlerV2(lessonService service_logic.ILessonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lessonIDStr := mux.Vars(r)["lessonId"]
		lessonID, err := strconv.ParseInt(lessonIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
			return
		}
		// Existence check
		if _, err := lessonService.GetLesson(lessonID); err != nil {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}
		var req types.ServerLessonPatchV2
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		if err := lessonService.UpdateLesson(lessonID, req.DurationMin, req.Format); err != nil {
			http.Error(w, "Error updating lesson", http.StatusBadRequest)
			return
		}
		lesson, err := lessonService.GetLesson(lessonID)
		if err != nil {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(types.MapperLessonServiceToServerV2(lesson))
		w.WriteHeader(http.StatusOK)
	}
}

func LessonDeleteHandlerV2(lessonService service_logic.ILessonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lessonIDStr := mux.Vars(r)["lessonId"]
		lessonID, err := strconv.ParseInt(lessonIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
			return
		}
		// Existence check
		if _, err := lessonService.GetLesson(lessonID); err != nil {
			http.Error(w, "Lesson not found", http.StatusNotFound)
			return
		}
		if err := lessonService.DeleteLesson(lessonID); err != nil {
			http.Error(w, "Error deleting lesson", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func ContractReviewsListHandlerV2(reviewService service_logic.IReviewService, contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contractIDStr := mux.Vars(r)["contractId"]
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}
		contract, err := contractService.GetContract(contractID)
		if err != nil {
			http.Error(w, "Contract not found", http.StatusNotFound)
			return
		}
		reviews := make([]types.ServiceReview, 0)
		if contract.ReviewClientID != 0 {
			review, err := reviewService.GetReview(contract.ReviewClientID)
			if err != nil {
				http.Error(w, "Error getting review", http.StatusBadRequest)
				return
			}
			reviews = append(reviews, *review)
		}
		if contract.ReviewRepetitorID != 0 {
			review, err := reviewService.GetReview(contract.ReviewRepetitorID)
			if err != nil {
				http.Error(w, "Error getting review", http.StatusBadRequest)
				return
			}
			reviews = append(reviews, *review)
		}
		serverReviews := make([]types.ServerReviewV2, len(reviews))
		for i, review := range reviews {
			serverReviews[i] = *types.MapperReviewServiceToServerV2(&review)
		}
		json.NewEncoder(w).Encode(serverReviews)
		w.WriteHeader(http.StatusOK)
	}
}

func ContractReviewCreateHandlerV2(reviewService service_logic.IReviewService, contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ServerReviewCreateV2
		contractIDStr := mux.Vars(r)["contractId"]
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		contract, err := contractService.GetContract(contractID)
		if err != nil {
			http.Error(w, "Contract not found", http.StatusNotFound)
			return
		}
		review := types.MapperReviewCreateV2ServerToService(&req, contractID, contract.ClientID, contract.RepetitorID)
		reviewID := int64(0)
		if req.SenderID == contract.ClientID {
			reviewID, err = contractService.CreateContractReviewClient(contractID, *review)
			if err != nil {
				http.Error(w, "Error creating review", http.StatusBadRequest)
				return
			}
		} else if req.SenderID == contract.RepetitorID {
			reviewID, err = contractService.CreateContractReviewRepetitor(contractID, *review)
			if err != nil {
				http.Error(w, "Error creating review", http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Invalid sender ID", http.StatusBadRequest)
			return
		}
		reviewV2 := types.ServerReviewV2{
			ID:         reviewID,
			ContractID: contractID,
			FromUserID: req.SenderID,
			ToUserID:   contract.ClientID,
			Score:      review.Rating,
			Text:       review.Comment,
			CreatedAt:  review.CreatedAt,
		}
		json.NewEncoder(w).Encode(reviewV2)
		w.WriteHeader(http.StatusCreated)
	}
}

func ContractTransactionsListHandlerV2(transactionService service_logic.ITransactionService, contractService service_logic.IContractService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contractIDStr := mux.Vars(r)["contractId"]
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}
		_, err = contractService.GetContract(contractID)
		if err != nil {
			http.Error(w, "Contract not found", http.StatusNotFound)
			return
		}
		offsetStr := r.URL.Query().Get("offset")
		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		sizeStr := r.URL.Query().Get("size")
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid size", http.StatusBadRequest)
			return
		}
		transactions, err := transactionService.GetContractTransactionsList(contractID, offset, size)
		if err != nil {
			http.Error(w, "Error getting transactions", http.StatusBadRequest)
			return
		}
		serverTransactions := make([]types.ServerTransactionV2, len(transactions))
		for i, transaction := range transactions {
			serverTransactions[i] = *types.MapperTransactionServiceToServerV2(&transaction)
		}
		json.NewEncoder(w).Encode(serverTransactions)
		w.WriteHeader(http.StatusOK)
	}
}

func ContractTransactionCreateHandlerV2(transactionService service_logic.ITransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contractIDStr := mux.Vars(r)["contractId"]
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid contract ID", http.StatusBadRequest)
			return
		}
		var req types.ServerTransactionCreateV2
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		// Note: creation will fail if contract doesn't exist; map to 404
		transactionID, err := transactionService.CreateContractPaymentTransaction(req.Amount, 0, contractID)
		if err != nil {
			fmt.Println("Error creating transaction: ", err)
			http.Error(w, "Error creating transaction", http.StatusBadRequest)
			return
		}
		transaction, err := transactionService.GetTransaction(transactionID)
		if err != nil {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}
		serverTx := types.MapperTransactionServiceToServerV2(transaction)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(*serverTx)
	}
}

func TransactionApproveHandlerV2(transactionService service_logic.ITransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactionIDStr := mux.Vars(r)["transactionId"]
		transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
			return
		}
		// Existence check
		if _, err := transactionService.GetTransaction(transactionID); err != nil {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}
		if err := transactionService.ApproveTransaction(transactionID); err != nil {
			http.Error(w, "Error approving transaction", http.StatusBadRequest)
			return
		}
		tx, err := transactionService.GetTransaction(transactionID)
		if err != nil {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(types.MapperTransactionServiceToServerV2(tx))
		w.WriteHeader(http.StatusOK)
	}
}

func SetupContractRouter(
	contractService service_logic.IContractService,
	reviewService service_logic.IReviewService,
	lessonService service_logic.ILessonService,
	logger *log.Logger,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(CONTRACT_GET, ContractGetContractHandler(contractService, logger))
	router.HandleFunc(CONTRACT_GET_REVIEW, ContractGetReviewHandler(reviewService, logger))
	router.HandleFunc(ADD_LESSON, ContractAddLessonHandler(lessonService, logger))
	router.HandleFunc(GET_LESSONS, ContractGetLessonsHandler(lessonService, logger))
	return router
}

func ContractGetContractHandler(contractService service_logic.IContractService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
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
		contract, err := contractService.GetContract(int64(contractID))
		if err != nil {
			logger.Printf("Error getting contract: %v", err)
			http.Error(w, "Error getting contract", http.StatusBadRequest)
			return
		}
		serverContract := types.MapperContractServiceToServer(contract)
		logger.Printf("Contract retrieved: %v", serverContract)
		json.NewEncoder(w).Encode(serverContract)
		w.WriteHeader(http.StatusOK)
	}
}

func ContractGetReviewHandler(reviewService service_logic.IReviewService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		reviewIDStr := r.URL.Query().Get("review_id")
		reviewID, err := strconv.Atoi(reviewIDStr)
		if err != nil {
			logger.Printf("Error converting reviewID to int: %v", err)
			http.Error(w, "Invalid reviewID", http.StatusBadRequest)
			return
		}
		logger.Printf("Review ID: %v", reviewID)
		review, err := reviewService.GetReview(int64(reviewID))
		if err != nil {
			logger.Printf("Error getting review: %v", err)
			http.Error(w, "Error getting review", http.StatusBadRequest)
			return
		}
		serverReview := types.MapperReviewServiceToServer(review)
		logger.Printf("Review retrieved: %v", serverReview)
		json.NewEncoder(w).Encode(serverReview)
	}
}

func ContractAddLessonHandler(lessonService service_logic.ILessonService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "POST" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		lesson := types.ServerLesson{}
		err := json.NewDecoder(r.Body).Decode(&lesson)
		if err != nil {
			logger.Printf("Error decoding lesson: %v", err)
			http.Error(w, "Error decoding lesson", http.StatusBadRequest)
			return
		}
		lesson.CreatedAt = time.Now()
		logger.Printf("Lesson: %v", lesson)
		serviceLesson := types.MapperLessonServerToService(&lesson)
		lessonID, err := lessonService.CreateLesson(*serviceLesson)
		if err != nil {
			logger.Printf("Error creating lesson: %v", err)
			http.Error(w, "Error creating lesson", http.StatusBadRequest)
			return
		}
		logger.Printf("Lesson created with ID: %v", lessonID)
		json.NewEncoder(w).Encode(lessonID)
	}
}

func ContractGetLessonsHandler(lessonService service_logic.ILessonService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		contractIDStr := r.URL.Query().Get("contract_id")
		contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting contractID to int: %v", err)
			http.Error(w, "Invalid contractID", http.StatusBadRequest)
			return
		}
		logger.Printf("Contract ID: %v", contractID)
		offsetStr := r.URL.Query().Get("lessons_offset")
		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting from to int: %v", err)
			http.Error(w, "Invalid from", http.StatusBadRequest)
			return
		}
		logger.Printf("Offset: %v", offset)
		sizeStr := r.URL.Query().Get("lessons_size")
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			logger.Printf("Error converting size to int: %v", err)
			http.Error(w, "Invalid size", http.StatusBadRequest)
			return
		}
		logger.Printf("Size: %v", size)
		lessons, err := lessonService.GetLessons(contractID, offset, size)
		if err != nil {
			logger.Printf("Error getting lessons: %v", err)
			http.Error(w, "Error getting lessons", http.StatusBadRequest)
			return
		}
		serverLessons := make([]types.ServerLesson, len(lessons))
		for i, lesson := range lessons {
			serverLessons[i] = *types.MapperLessonServiceToServer(&lesson)
		}
		logger.Printf("Lessons retrieved: %v", serverLessons)
		json.NewEncoder(w).Encode(serverLessons)
		w.WriteHeader(http.StatusOK)
	}
}
