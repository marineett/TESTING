package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func SetupAuthorizeRouter(authService service_logic.IAuthService, logger *log.Logger) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(AUTH_AUTHORIZE, AuthorizeHandler(authService, logger))

	return router
}

func AuthorizeHandler(authService service_logic.IAuthService, logger *log.Logger) http.HandlerFunc {
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
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var authData types.ServerAuthData
		if err := json.Unmarshal(body, &authData); err != nil {
			logger.Printf("Error unmarshaling request body: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		verdict, err := authService.Authorize(*types.MapperAuthServerToService(&authData))
		if err != nil {
			logger.Printf("Error authorizing: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		serverVerdict := types.MapperVerdictServiceToServer(&verdict)
		json.NewEncoder(w).Encode(serverVerdict)
		logger.Printf("Authorized: %v", verdict)
	}
}
