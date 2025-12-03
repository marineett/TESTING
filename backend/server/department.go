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

	"github.com/gorilla/mux"
)

func SetupDepartmentRouterV2(departmentService service_logic.IDepartmentService, moderatorService service_logic.IModeratorService) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(DEPARTMENTS_V2, AdminCreateDepartmentHandlerV2(departmentService)).Methods("POST")
	r.HandleFunc(ADMIN_DEPARTMENTS_V2, AdminListDepartmentsHandlerV2(departmentService)).Methods("GET")
	r.HandleFunc(EXACT_DEPARTMENT_V2, DepartmentReplaceHandlerV2(departmentService)).Methods("PUT")
	r.HandleFunc(EXACT_DEPARTMENT_V2, DepartmentDeleteHandlerV2(departmentService)).Methods("DELETE")
	r.HandleFunc(DEPARTMENT_MODERATOR_V2, DepartmentAssignModeratorHandlerV2(departmentService)).Methods("PUT")
	r.HandleFunc(DEPARTMENT_MODERATOR_V2, DepartmentRemoveModeratorHandlerV2(departmentService)).Methods("DELETE")
	return r
}

func AdminCreateDepartmentHandlerV2(departmentService service_logic.IDepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ServerDepartmentCreateV2
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		departmentId, err := departmentService.CreateDepartment(types.ServiceDepartmentInitData{
			Name:   req.Name,
			HeadID: req.HeadID,
		})
		if err != nil {
			http.Error(w, "Error creating department", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(types.ServerDepartmentV2{
			ID:         departmentId,
			Name:       req.Name,
			HeadID:     req.HeadID,
			Moderators: []types.ServerModeratorProfileWithIDV2{},
		})
	}
}

func AdminListDepartmentsHandlerV2(departmentService service_logic.IDepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminId := mux.Vars(r)["adminId"]
		adminIdInt, err := strconv.ParseInt(adminId, 10, 64)
		if err != nil {
			http.Error(w, "Invalid admin ID", http.StatusBadRequest)
			return
		}
		departments, err := departmentService.GetDepartmentsByHeadIdWithModerators(adminIdInt)
		if err != nil {
			fmt.Print(err.Error())
			http.Error(w, "Error getting departments", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(departments)
	}
}

func DepartmentReplaceHandlerV2(departmentService service_logic.IDepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ServerDepartmentNameUpdateV2
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		departmentId := mux.Vars(r)["departmentId"]
		departmentIdInt, err := strconv.ParseInt(departmentId, 10, 64)
		if err != nil {
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		// Existence check
		if _, err := departmentService.GetDepartment(departmentIdInt); err != nil {
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}
		err = departmentService.UpdateDepartmentName(departmentIdInt, req.Name)
		if err != nil {
			http.Error(w, "Error updating department name", http.StatusBadRequest)
			return
		}
		updatedDepartment, err := departmentService.GetDepartment(departmentIdInt)
		if err != nil {
			http.Error(w, "Error getting department", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(updatedDepartment)
	}
}

func DepartmentDeleteHandlerV2(departmentService service_logic.IDepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		departmentId := mux.Vars(r)["departmentId"]
		departmentIdInt, err := strconv.ParseInt(departmentId, 10, 64)
		if err != nil {
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		// Existence check
		if _, err := departmentService.GetDepartment(departmentIdInt); err != nil {
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}
		err = departmentService.DeleteDepartment(departmentIdInt)
		if err != nil {
			http.Error(w, "Error deleting department", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func DepartmentAssignModeratorHandlerV2(departmentService service_logic.IDepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		moderatorId := mux.Vars(r)["moderatorId"]
		moderatorIdInt, err := strconv.ParseInt(moderatorId, 10, 64)
		if err != nil {
			http.Error(w, "Invalid moderator ID", http.StatusBadRequest)
			return
		}
		departmentId := mux.Vars(r)["departmentId"]
		departmentIdInt, err := strconv.ParseInt(departmentId, 10, 64)
		if err != nil {
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		// Existence check
		if _, err := departmentService.GetDepartment(departmentIdInt); err != nil {
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}
		err = departmentService.AssignModeratorToDepartment(moderatorIdInt, departmentIdInt)
		if err != nil {
			http.Error(w, "Error assigning moderator to department", http.StatusBadRequest)
			return
		}
		updatedDepartment, err := departmentService.GetDepartment(departmentIdInt)
		if err != nil {
			http.Error(w, "Error getting department", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedDepartment)
	}
}

func DepartmentRemoveModeratorHandlerV2(departmentService service_logic.IDepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		moderatorId := mux.Vars(r)["moderatorId"]
		moderatorIdInt, err := strconv.ParseInt(moderatorId, 10, 64)
		if err != nil {
			http.Error(w, "Invalid moderator ID", http.StatusBadRequest)
			return
		}
		departmentId := mux.Vars(r)["departmentId"]
		departmentIdInt, err := strconv.ParseInt(departmentId, 10, 64)
		if err != nil {
			http.Error(w, "Invalid department ID", http.StatusBadRequest)
			return
		}
		// Existence check
		if _, err := departmentService.GetDepartment(departmentIdInt); err != nil {
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}
		err = departmentService.FireModeratorFromDepartment(moderatorIdInt, departmentIdInt)
		if err != nil {
			http.Error(w, "Error removing moderator from department", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func SetupDepartmentRouter(
	adminService service_logic.IAdminService,
	departmentService service_logic.IDepartmentService,
	moderatorService service_logic.IModeratorService,
	logger *log.Logger,
) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc(ADMIN_CREATE_DEPARTMENT, AdminCreateDepartmentHandler(adminService, departmentService, logger))
	router.HandleFunc(ADMIN_GET_DEPARTMENTS, AdminGetDepartmentsHandler(departmentService, moderatorService, logger))
	router.HandleFunc(ADMIN_HIRE_MODERATOR, AdminHireModeratorHandler(departmentService, logger))
	router.HandleFunc(ADMIN_FIRE_MODERATOR, AdminFireModeratorHandler(departmentService, logger))
	return router
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
		_, err = departmentService.CreateDepartment(department)
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
