package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func SetupAdminRouterV2(adminService service_logic.IAdminService, moderatorService service_logic.IModeratorService, departmentService service_logic.IDepartmentService) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(EXACT_ADMIN_V2, AdminGetProfileHandlerV2(adminService)).Methods("GET")
	return r
}

func AdminGetProfileHandlerV2(adminService service_logic.IAdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminID := mux.Vars(r)["adminId"]
		adminIDInt, err := strconv.Atoi(adminID)
		if err != nil {
			http.Error(w, "Invalid admin ID", http.StatusBadRequest)
			return
		}
		adminProfile, err := adminService.GetAdminProfile(int64(adminIDInt))
		if err != nil {
			http.Error(w, "Error getting admin profile", http.StatusBadRequest)
			return
		}
		serverAdminProfile := *types.MapperAdminProfileServiceToServer(adminProfile)
		json.NewEncoder(w).Encode(serverAdminProfile)
		w.WriteHeader(http.StatusOK)
	}
}

func SetupAdminRouter(
	adminService service_logic.IAdminService,
	departmentService service_logic.IDepartmentService,
	moderatorService service_logic.IModeratorService,
	logger *log.Logger,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(ADMIN_GET_PROFILE, AdminGetProfileHandler(adminService, logger))
	return router
}

func AdminGetProfileHandler(adminService service_logic.IAdminService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		userIDStr := r.URL.Query().Get("id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid user ID: %v", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		logger.Printf("User ID: %v", userID)
		adminProfile, err := adminService.GetAdminProfile(userID)
		if err != nil {
			logger.Printf("Failed to get admin profile: %v", err)
			http.Error(w, "Failed to get admin profile", http.StatusBadRequest)
			return
		}
		serverAdminProfile := types.MapperAdminProfileServiceToServer(adminProfile)
		logger.Printf("Got admin profile: %v", serverAdminProfile)
		json.NewEncoder(w).Encode(serverAdminProfile)
	}
}

// moved V1 moderator handlers to moderator.go
