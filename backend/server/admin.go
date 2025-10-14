package server

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func SetupAdminRouter(
	adminService service_logic.IAdminService,
	departmentService service_logic.IDepartmentService,
	moderatorService service_logic.IModeratorService,
	logger *log.Logger,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(ADMIN_GET_PROFILE, AdminGetProfileHandler(adminService, logger))
	router.HandleFunc(ADMIN_CREATE_DEPARTMENT, AdminCreateDepartmentHandler(adminService, departmentService, logger))
	router.HandleFunc(ADMIN_GET_DEPARTMENTS, AdminGetDepartmentsHandler(departmentService, moderatorService, logger))
	router.HandleFunc(ADMIN_GET_MODERATORS, AdminGetModeratorsHandler(moderatorService, logger))
	router.HandleFunc(ADMIN_HIRE_MODERATOR, AdminHireModeratorHandler(departmentService, logger))
	router.HandleFunc(ADMIN_FIRE_MODERATOR, AdminFireModeratorHandler(departmentService, logger))
	router.HandleFunc(ADMIN_CHANGE_MODERATOR_SALARY, AdminChangeModeratorSalaryHandler(moderatorService, logger))
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
			http.Error(w, "Failed to get admin profile", http.StatusInternalServerError)
			return
		}
		serverAdminProfile := types.MapperAdminProfileServiceToServer(adminProfile)
		logger.Printf("Got admin profile: %v", serverAdminProfile)
		json.NewEncoder(w).Encode(serverAdminProfile)
	}
}

func AdminCreateDepartmentHandler(
	adminService service_logic.IAdminService,
	departmentService service_logic.IDepartmentService,
	logger *log.Logger,
) http.HandlerFunc {
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
		departmentName := r.URL.Query().Get("name")
		if departmentName == "" {
			logger.Printf("Error:Department name is required")
			http.Error(w, "Department name is required", http.StatusBadRequest)
			return
		}
		logger.Printf("Department name: %v", departmentName)
		department := types.ServiceDepartmentInitData{
			Name:   departmentName,
			HeadID: userID,
		}
		err = departmentService.CreateDepartment(department)
		if err != nil {
			logger.Printf("Failed to create department: %v", err)
			http.Error(w, "Failed to create department", http.StatusInternalServerError)
			return
		}
		logger.Printf("Department created successfully")
		w.WriteHeader(http.StatusCreated)
	}
}

func AdminGetDepartmentsHandler(departmentService service_logic.IDepartmentService, moderatorService service_logic.IModeratorService, logger *log.Logger) http.HandlerFunc {
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
		departments, err := departmentService.GetDepartmentsByHeadID(userID)
		if err != nil {
			logger.Printf("Failed to get departments: %v", err)
			http.Error(w, "Failed to get departments", http.StatusInternalServerError)
			return
		}
		logger.Printf("Got departments: %v", departments)
		completeDepartments := make([]types.ServerDepartment, len(departments))
		for i, department := range departments {
			moderatorsIDs, err := departmentService.GetDepartmentUsersIDs(department.ID)
			if err != nil {
				logger.Printf("Failed to get complete department info: %v", err)
				http.Error(w, "Failed to get complete department info", http.StatusInternalServerError)
				return
			}
			moderators := make([]types.ServerModeratorProfileWithID, len(moderatorsIDs))
			for j, moderatorID := range moderatorsIDs {
				moderator, err := moderatorService.GetModeratorProfileWithId(moderatorID)
				if err != nil {
					logger.Printf("Failed to get complete department info: %v", err)
					http.Error(w, "Failed to get complete department info", http.StatusInternalServerError)
					return
				}
				moderators[j] = *types.MapperModeratorProfileWithIDServiceToServer(moderator)
			}
			completeDepartments[i] = types.ServerDepartment{
				Name:       department.Name,
				HeadID:     department.HeadID,
				Moderators: moderators,
			}
		}
		logger.Printf("Got complete departments: %v", completeDepartments)
		json.NewEncoder(w).Encode(completeDepartments)
	}
}

func AdminGetDepartmentHandler(departmentService service_logic.IDepartmentService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		departmentIDStr := r.URL.Query().Get("id")
		departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid department ID: %v", err)
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		logger.Printf("Department ID: %v", departmentID)
		department, err := departmentService.GetDepartment(departmentID)
		if err != nil {
			logger.Printf("Failed to get department: %v", err)
			http.Error(w, "Failed to get department", http.StatusInternalServerError)
			return
		}
		serverDepartment := types.MapperDepartmentServiceToServer(&department)
		logger.Printf("Got department: %v", serverDepartment)
		json.NewEncoder(w).Encode(serverDepartment)
	}
}

func AdminGetModeratorsHandler(moderatorService service_logic.IModeratorService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		moderators, err := moderatorService.GetModerators()
		if err != nil {
			logger.Printf("Failed to get moderators: %v", err)
			http.Error(w, "Failed to get moderators", http.StatusInternalServerError)
			return
		}
		serverModerators := make([]types.ServerModeratorProfileWithID, len(moderators))
		for i, moderator := range moderators {
			serverModerators[i] = *types.MapperModeratorProfileWithIDServiceToServer(moderator)
		}
		logger.Printf("Got moderators: %v", serverModerators)
		json.NewEncoder(w).Encode(serverModerators)
	}
}

func AdminHireModeratorHandler(
	departmentService service_logic.IDepartmentService,
	logger *log.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		adminIDStr := r.URL.Query().Get("id")
		adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid admin ID: %v", err)
			http.Error(w, "Invalid admin ID", http.StatusBadRequest)
			return
		}
		logger.Printf("Admin ID: %v", adminID)
		departmentIDStr := r.URL.Query().Get("d_id")
		departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid department ID: %v", err)
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		logger.Printf("Department ID: %v", departmentID)
		moderatorIDStr := r.URL.Query().Get("m_id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid moderator ID: %v", err)
			http.Error(w, "Invalid moderator ID", http.StatusBadRequest)
			return
		}
		logger.Printf("Moderator ID: %v", moderatorID)
		department, err := departmentService.GetDepartment(departmentID)
		if err != nil {
			logger.Printf("Failed to get departments: %v", err)
			http.Error(w, "Failed to get departments", http.StatusInternalServerError)
			return
		}
		logger.Printf("Got department: %v", department)
		departments, err := departmentService.GetUserDepartmentsIDs(moderatorID)
		if err != nil {
			logger.Printf("Failed to get departments: %v", err)
			http.Error(w, "Failed to get departments", http.StatusInternalServerError)
			return
		}
		logger.Printf("Got departments: %v", departments)
		for _, departmentID := range departments {
			if departmentID == department.ID {
				logger.Printf("Moderator is already in this department")
				http.Error(w, "Moderator is already in this department", http.StatusBadRequest)
				return
			}
		}
		if department.HeadID != adminID {
			logger.Printf("You are not the head of this department")
			http.Error(w, "You are not the head of this department", http.StatusBadRequest)
			return
		}
		err = departmentService.AssignModeratorToDepartment(moderatorID, departmentID)
		if err != nil {
			logger.Printf("Failed to hire moderator: %v", err)
			http.Error(w, "Failed to hire moderator", http.StatusInternalServerError)
			return
		}
		logger.Printf("Hired moderator successfully")
		w.WriteHeader(http.StatusOK)
	}
}

func AdminFireModeratorHandler(departmentService service_logic.IDepartmentService, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		adminIDStr := r.URL.Query().Get("id")
		adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid admin ID: %v", err)
			http.Error(w, "Invalid admin ID", http.StatusBadRequest)
			return
		}
		logger.Printf("Admin ID: %v", adminID)
		departmentIDStr := r.URL.Query().Get("d_id")
		departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid department ID: %v", err)
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		logger.Printf("Department ID: %v", departmentID)
		moderatorIDStr := r.URL.Query().Get("m_id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid moderator ID: %v", err)
			http.Error(w, "Invalid moderator ID", http.StatusBadRequest)
			return
		}
		logger.Printf("Moderator ID: %v", moderatorID)
		departments, err := departmentService.GetDepartment(departmentID)
		if err != nil {
			logger.Printf("Failed to get departments: %v", err)
			http.Error(w, "Failed to get departments", http.StatusInternalServerError)
			return
		}
		logger.Printf("Got department: %v", departments)
		if departments.HeadID != adminID {
			logger.Printf("You are not the head of this department")
			http.Error(w, "You are not the head of this department", http.StatusBadRequest)
			return
		}
		err = departmentService.FireModeratorFromDepartment(moderatorID, departmentID)
		if err != nil {
			logger.Printf("Failed to fire moderator: %v", err)
			http.Error(w, "Failed to fire moderator", http.StatusInternalServerError)
			return
		}
		logger.Printf("Fired moderator successfully")
		w.WriteHeader(http.StatusOK)
	}
}

func AdminChangeModeratorSalaryHandler(
	moderatorService service_logic.IModeratorService,
	logger *log.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		if r.Method != "GET" {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		newSalaryStr := r.URL.Query().Get("salary")
		newSalary, err := strconv.ParseInt(newSalaryStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid salary: %v", err)
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		logger.Printf("New salary: %v", newSalary)
		moderatorIDStr := r.URL.Query().Get("m_id")
		moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
		if err != nil {
			logger.Printf("Invalid moderator ID: %v", err)
			http.Error(w, "Invalid moderator ID", http.StatusBadRequest)
			return
		}
		logger.Printf("Moderator ID: %v", moderatorID)
		err = moderatorService.UpdateModeratorSalary(moderatorID, newSalary)
		if err != nil {
			logger.Printf("Failed to change moderator salary: %v", err)
			http.Error(w, "Failed to change moderator salary", http.StatusInternalServerError)
			return
		}
		logger.Printf("Changed moderator salary successfully")
		w.WriteHeader(http.StatusOK)
	}
}
