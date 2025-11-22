package server_http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"data_base_project/service_logic"
)

type httpServer struct {
	services *service_logic.ServiceModule
	logger   *log.Logger
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (s *httpServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func (s *httpServer) handleGetClient(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad id")
		return
	}
	profile, err := s.services.ClientService.GetClientProfile(id, 0, 100)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (s *httpServer) handleGetRepetitor(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad id")
		return
	}
	profile, err := s.services.RepetitorService.GetRepetitorProfile(id, 0, 100)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (s *httpServer) handleGetContract(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad id")
		return
	}
	contract, err := s.services.ContractService.GetContract(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, contract)
}

func (s *httpServer) handleGetLessons(w http.ResponseWriter, r *http.Request) {
	contractIDStr := r.URL.Query().Get("contractId")
	if contractIDStr == "" {
		writeError(w, http.StatusBadRequest, "missing contractId")
		return
	}
	contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad contractId")
		return
	}
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if limit <= 0 {
		limit = 100
	}
	lessons, err := s.services.LessonService.GetLessons(contractID, offset, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, lessons)
}

func (s *httpServer) handleGetClientReviews(w http.ResponseWriter, r *http.Request) {
	clientIDStr := r.URL.Query().Get("clientId")
	if clientIDStr == "" {
		writeError(w, http.StatusBadRequest, "missing clientId")
		return
	}
	clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad clientId")
		return
	}
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if limit <= 0 {
		limit = 100
	}
	reviews, err := s.services.ReviewService.GetReviewsByClientID(clientID, offset, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, reviews)
}

func (s *httpServer) handleGetRepetitorReviews(w http.ResponseWriter, r *http.Request) {
	repetitorIDStr := r.URL.Query().Get("repetitorId")
	if repetitorIDStr == "" {
		writeError(w, http.StatusBadRequest, "missing repetitorId")
		return
	}
	repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad repetitorId")
		return
	}
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if limit <= 0 {
		limit = 100
	}
	reviews, err := s.services.ReviewService.GetReviewsByRepetitorID(repetitorID, offset, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, reviews)
}

func StartHTTPServer(services *service_logic.ServiceModule, logger *log.Logger) *http.Server {
	s := &httpServer{services: services, logger: logger}
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/client", s.handleGetClient)
	mux.HandleFunc("/api/repetitor", s.handleGetRepetitor)
	mux.HandleFunc("/api/contract", s.handleGetContract)
	mux.HandleFunc("/api/lessons", s.handleGetLessons)
	mux.HandleFunc("/api/reviews/client", s.handleGetClientReviews)
	mux.HandleFunc("/api/reviews/repetitor", s.handleGetRepetitorReviews)

	addr := ":8081"
	if p := os.Getenv("CONSOLE_SERVER_PORT"); p != "" {
		if p[0] == ':' {
			addr = p
		} else {
			addr = ":" + p
		}
	}
	server := &http.Server{Addr: addr, Handler: mux}
	go func() {
		logger.Printf("console http server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Printf("http server error: %v", err)
		}
	}()
	return server
}
