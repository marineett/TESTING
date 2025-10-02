package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

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
			http.Error(w, "Error getting contract", http.StatusInternalServerError)
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
			http.Error(w, "Error getting review", http.StatusInternalServerError)
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
			http.Error(w, "Error creating lesson", http.StatusInternalServerError)
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
			http.Error(w, "Error getting lessons", http.StatusInternalServerError)
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
