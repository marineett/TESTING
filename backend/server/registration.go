package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func SetupRegistrationRouterV2(
	authService service_logic.IAuthService,
	moderatorService service_logic.IModeratorService,
	clientService service_logic.IClientService,
	adminService service_logic.IAdminService,
	repetitorService service_logic.IRepetitorService,
) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(REGISTRATION_API_V2, RegistrationHandlerV2(clientService, moderatorService, adminService, repetitorService, authService)).Methods("POST")
	return router
}

func RegistrationHandlerV2(
	clientService service_logic.IClientService,
	moderatorService service_logic.IModeratorService,
	adminService service_logic.IAdminService,
	repetitorService service_logic.IRepetitorService,
	authService service_logic.IAuthService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		var registrationData types.ServerRegistrationDataV2
		if err := json.Unmarshal(body, &registrationData); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		switch registrationData.Role {
		case "client":
			err = clientService.CreateClient(*types.MapperRegistrationV2ToServiceInitClient(&registrationData))
		case "moderator":
			err = moderatorService.CreateModerator(*types.MapperRegistrationV2ToServiceInitModerator(&registrationData))
		case "admin":
			err = adminService.CreateAdmin(*types.MapperRegistrationV2ToServiceInitAdmin(&registrationData, 0))
		case "repetitor":
			err = repetitorService.CreateRepetitor(*types.MapperRegistrationV2ToServiceInitRepetitor(&registrationData))
		default:
			http.Error(w, "Invalid user type", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		verdict, err := authService.Authorize(types.ServiceAuthData{Login: registrationData.Login, Password: registrationData.Password})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			http.Error(w, "JWT secret is not configured", http.StatusInternalServerError)
			return
		}
		token, err := createJWT(verdict.UserID, verdict.UserType.String(), 24*time.Hour, secret)
		if err != nil {
			http.Error(w, "Failed to issue token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(types.ServerAuthResponseV2{Token: token, Role: verdict.UserType.String(), UserID: verdict.UserID})
	}
}

func SetupRegistrationRouter(
	authService service_logic.IAuthService,
	moderatorService service_logic.IModeratorService,
	clientService service_logic.IClientService,
	adminService service_logic.IAdminService,
	repetitorService service_logic.IRepetitorService,
	logger *log.Logger,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(REGISTRATION_CLIENT, RegistrationClientHandler(clientService, authService, logger))
	router.HandleFunc(REGISTRATION_MODERATOR, RegistrationModeratorHandler(moderatorService, authService, logger))
	router.HandleFunc(REGISTRATION_ADMIN, RegistrationAdminHandler(adminService, authService, logger))
	router.HandleFunc(REGISTRATION_REPETITOR, RegistrationRepetitorHandler(repetitorService, authService, logger))
	return router
}

func RegistrationModeratorHandler(
	moderatorService service_logic.IModeratorService,
	authService service_logic.IAuthService,
	logger *log.Logger,
) http.HandlerFunc {
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
		logger.Printf("Request body: %s", string(body))

		var initData types.ServerInitModeratorData
		if err := json.Unmarshal(body, &initData); err != nil {
			logger.Printf("Error unmarshaling request body: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		logger.Printf("Parsed init data: %+v", initData)

		inSystem, err := authService.CheckLogin(initData.Login)
		if err != nil {
			logger.Printf("Error checking login: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if inSystem {
			logger.Printf("User already exists: %s", initData.Login)
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		serviceInitData := types.MapperInitModeratorServerToService(&initData)
		err = moderatorService.CreateModerator(*serviceInitData)
		if err != nil {
			logger.Printf("Error creating moderator: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Printf("Moderator created successfully: %s", initData.Login)
		w.WriteHeader(http.StatusCreated)
	}
}

func RegistrationClientHandler(
	clientService service_logic.IClientService,
	authService service_logic.IAuthService,
	logger *log.Logger,
) http.HandlerFunc {
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
		logger.Printf("Request body: %s", string(body))

		var initData types.ServerInitClientData
		if err := json.Unmarshal(body, &initData); err != nil {
			logger.Printf("Error unmarshaling request body: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		logger.Printf("Parsed init data: %+v", initData)

		inSystem, err := authService.CheckLogin(initData.Login)
		if err != nil {
			logger.Printf("Error checking login: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if inSystem {
			logger.Printf("User already exists: %s", initData.Login)
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		serviceInitData := types.MapperInitClientServerToService(&initData)
		err = clientService.CreateClient(*serviceInitData)
		if err != nil {
			log.Printf("Error creating client: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Printf("Client created successfully: %s", initData.Login)
		w.WriteHeader(http.StatusCreated)
	}
}

func RegistrationAdminHandler(
	adminService service_logic.IAdminService,
	authService service_logic.IAuthService,
	logger *log.Logger,
) http.HandlerFunc {
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
		logger.Printf("Request body: %s", string(body))

		var initData types.ServerInitAdminData
		if err := json.Unmarshal(body, &initData); err != nil {
			logger.Printf("Error unmarshaling request body: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		logger.Printf("Parsed init data: %+v", initData)

		inSystem, err := authService.CheckLogin(initData.Login)
		if err != nil {
			logger.Printf("Error checking login: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if inSystem {
			logger.Printf("User already exists: %s", initData.Login)
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		serviceInitData := types.MapperInitAdminServerToService(&initData)
		err = adminService.CreateAdmin(*serviceInitData)
		if err != nil {
			logger.Printf("Error creating admin: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Printf("Admin created successfully: %s", initData.Login)
		w.WriteHeader(http.StatusCreated)
	}
}

func RegistrationRepetitorHandler(
	repetitorService service_logic.IRepetitorService,
	authService service_logic.IAuthService,
	logger *log.Logger,
) http.HandlerFunc {
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
		logger.Printf("Request body: %s", string(body))

		var initData types.ServerInitRepetitorData
		if err := json.Unmarshal(body, &initData); err != nil {
			logger.Printf("Error unmarshaling request body: %v", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		logger.Printf("Parsed init data: %+v", initData)

		inSystem, err := authService.CheckLogin(initData.Login)
		if err != nil {
			logger.Printf("Error checking login: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if inSystem {
			logger.Printf("User already exists: %s", initData.Login)
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		serviceInitData := types.MapperInitRepetitorServerToService(&initData)
		err = repetitorService.CreateRepetitor(*serviceInitData)
		if err != nil {
			logger.Printf("Error creating repetitor: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Printf("Repetitor created successfully: %s", initData.Login)
		w.WriteHeader(http.StatusCreated)
	}
}
