package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

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
