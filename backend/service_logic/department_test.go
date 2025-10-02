package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
)

func TestCreateDepartmentCorrectLondon(t *testing.T) {
	departmentRepository := tu.CreateTestDepartmentRepository()
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err := departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	department, err := departmentRepository.GetDepartment(1)
	if err != nil {
		t.Fatalf("Error getting department: %v", err)
	}
	if department.Name != tu.TestInitDepartmentData.Name {
		t.Fatalf("Department not found: %v", department)
	}
	if department.HeadID != tu.TestInitDepartmentData.HeadID {
		t.Fatalf("Department not found: %v", department)
	}
}

func TestCreateDepartmentCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := module.DepartmentRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	moderatorRepository := module.ModeratorRepository
	adminRepository := module.AdminRepository
	authRepository := module.AuthRepository
	adminService := CreateAdminService(adminRepository, userRepository, personalDataRepository)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	tu.TestInitDepartmentData.HeadID = result.UserID
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err = departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
}

func TestGetDepartmentsByHeadIDCorrectLondon(t *testing.T) {
	departmentRepository := tu.CreateTestDepartmentRepository()
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err := departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	departments, err := departmentService.GetDepartmentsByHeadID(tu.TestInitDepartmentData.HeadID)
	if err != nil {
		t.Fatalf("Error getting departments by head id: %v", err)
	}
	if len(departments) != 1 {
		t.Fatalf("Departments not found: %v", departments)
	}
	if departments[0].Name != tu.TestInitDepartmentData.Name {
		t.Fatalf("Department not found: %v", departments)
	}
}

func TestGetDepartmentsByHeadIDCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := module.DepartmentRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	moderatorRepository := module.ModeratorRepository
	adminRepository := module.AdminRepository
	authRepository := module.AuthRepository
	adminService := CreateAdminService(adminRepository, userRepository, personalDataRepository)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	tu.TestInitDepartmentData.HeadID = result.UserID
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err = departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	err = departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	departments, err := departmentService.GetDepartmentsByHeadID(tu.TestInitDepartmentData.HeadID)
	if err != nil {
		t.Fatalf("Error getting departments by head id: %v", err)
	}
	if len(departments) != 2 {
		t.Fatalf("Departments not found: %v", departments)
	}
	if departments[0].HeadID != tu.TestInitDepartmentData.HeadID || departments[1].HeadID != tu.TestInitDepartmentData.HeadID {
		t.Fatalf("Department not found: %v", departments)
	}
}

func TestGetDepartmentCorrectLondon(t *testing.T) {
	departmentRepository := tu.CreateTestDepartmentRepository()
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err := departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	department, err := departmentService.GetDepartment(1)
	if err != nil {
		t.Fatalf("Error getting department: %v", err)
	}
	if department.Name != tu.TestInitDepartmentData.Name {
		t.Fatalf("Department not found: %v", department)
	}
	if department.HeadID != tu.TestInitDepartmentData.HeadID {
		t.Fatalf("Department not found: %v", department)
	}
}

func TestGetDepartmentCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := module.DepartmentRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	moderatorRepository := module.ModeratorRepository
	adminRepository := module.AdminRepository
	authRepository := module.AuthRepository
	adminService := CreateAdminService(adminRepository, userRepository, personalDataRepository)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	tu.TestInitDepartmentData.HeadID = result.UserID
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err = departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	departments, err := departmentService.GetDepartmentsByHeadID(tu.TestInitDepartmentData.HeadID)
	if err != nil {
		t.Fatalf("Error getting departments by head id: %v", err)
	}
	if len(departments) != 1 {
		t.Fatalf("Departments not found: %v", departments)
	}
	department, err := departmentService.GetDepartment(departments[0].ID)
	if err != nil {
		t.Fatalf("Error getting department: %v", err)
	}
	if department.Name != tu.TestInitDepartmentData.Name {
		t.Fatalf("Department not found: %v", department)
	}
	if department.HeadID != tu.TestInitDepartmentData.HeadID {
		t.Fatalf("Department not found: %v", department)
	}
}

func TestAssignAdminToDepartmentCorrectLondon(t *testing.T) {
	departmentRepository := tu.CreateTestDepartmentRepository()
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	initData := tu.TestInitDepartmentData
	initData.HeadID = 0
	err := departmentService.CreateDepartment(initData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	err = departmentService.AssignAdminToDepartment(1, 1)
	if err != nil {
		t.Fatalf("Error assigning admin to department: %v", err)
	}
	department, err := departmentRepository.GetDepartment(1)
	if err != nil {
		t.Fatalf("Error getting department: %v", err)
	}
	if department.HeadID != 1 {
		t.Fatalf("head id not updated: %v", department)
	}
}

func TestAssignAdminToDepartmentCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := module.DepartmentRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	moderatorRepository := module.ModeratorRepository
	adminRepository := module.AdminRepository
	authRepository := module.AuthRepository
	adminService := CreateAdminService(adminRepository, userRepository, personalDataRepository)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	tu.TestInitDepartmentData.HeadID = 0
	err = departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	departments, err := departmentService.GetDepartmentsByHeadID(0)
	if err != nil {
		t.Fatalf("Error getting departments by head id: %v", err)
	}
	if len(departments) != 1 {
		t.Fatalf("Departments not found: %v", departments)
	}
	err = departmentService.AssignAdminToDepartment(result.UserID, departments[0].ID)
	if err != nil {
		t.Fatalf("Error assigning admin to department: %v", err)
	}
	department, err := departmentRepository.GetDepartment(departments[0].ID)
	if err != nil {
		t.Fatalf("Error getting department: %v", err)
	}
	if department.HeadID != result.UserID {
		t.Fatalf("Department not found: %v", department)
	}
}

func TestAssignAdminToDepartmentIncorrectLondon(t *testing.T) {
	departmentRepository := tu.CreateTestDepartmentRepository()
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err := departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	err = departmentService.AssignAdminToDepartment(1, 2)
	if err == nil {
		t.Fatalf("No error assigning admin to department with wrong head id: %v", err)
	}
}

func TestAssignAdminToDepartmentIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := module.DepartmentRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	moderatorRepository := module.ModeratorRepository
	adminRepository := module.AdminRepository
	adminService := CreateAdminService(adminRepository, userRepository, personalDataRepository)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err = departmentService.AssignAdminToDepartment(1, 2)
	if err == nil {
		t.Fatalf("No error assigning admin to department: %v", err)
	}
}

func TestFireAdminFromDepartmentCorrectLondon(t *testing.T) {
	departmentRepository := tu.CreateTestDepartmentRepository()
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err := departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	err = departmentService.FireAdminFromDepartment(tu.TestInitDepartmentData.HeadID, 1)
	if err != nil {
		t.Fatalf("Error firing admin from department: %v", err)
	}
	department, err := departmentRepository.GetDepartment(1)
	if err != nil {
		t.Fatalf("Error getting department: %v", err)
	}
	if department.HeadID != 0 {
		t.Fatalf("wrong head id: %v", department)
	}
}

func TestFireAdminFromDepartmentCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := module.DepartmentRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	moderatorRepository := module.ModeratorRepository
	adminRepository := module.AdminRepository
	authRepository := module.AuthRepository
	adminService := CreateAdminService(adminRepository, userRepository, personalDataRepository)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	tu.TestInitDepartmentData.HeadID = result.UserID
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err = departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	departments, err := departmentService.GetDepartmentsByHeadID(result.UserID)
	if err != nil {
		t.Fatalf("Error getting departments by head id: %v", err)
	}
	if len(departments) != 1 {
		t.Fatalf("Departments not found: %v", departments)
	}
	err = departmentService.FireAdminFromDepartment(result.UserID, departments[0].ID)
	if err != nil {
		t.Fatalf("Error firing admin from department: %v", err)
	}
	department, err := departmentRepository.GetDepartment(departments[0].ID)
	if err != nil {
		t.Fatalf("Error getting department: %v", err)
	}
	if department.HeadID != 0 {
		t.Fatalf("wrong head id: %v", department)
	}
}

func TestFireAdminFromDepartmentIncorrectLondon(t *testing.T) {
	departmentRepository := tu.CreateTestDepartmentRepository()
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err := departmentService.FireAdminFromDepartment(tu.TestInitDepartmentData.HeadID+1, tu.TestInitDepartmentData.ID)
	if err == nil {
		t.Fatalf("No error firing admin from department: %v", err)
	}
}

func TestFireAdminFromDepartmentIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := module.DepartmentRepository
	moderatorRepository := module.ModeratorRepository
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err = departmentService.FireAdminFromDepartment(1, 2)
	if err == nil {
		t.Fatalf("No error firing admin from department: %v", err)
	}
}

func TestFireModeratorFromDepartmentCorrectLondon(t *testing.T) {
	departmentRepository := tu.CreateTestDepartmentRepository()
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	err := departmentService.CreateDepartment(tu.TestInitDepartmentData)
	if err != nil {
		t.Fatalf("Error creating department: %v", err)
	}
	err = departmentService.FireModeratorFromDepartment(tu.TestInitDepartmentData.HeadID, 1)
	if err != nil {
		t.Fatalf("Error firing moderator from department: %v", err)
	}
	department, err := departmentRepository.GetDepartment(1)
	if err != nil {
		t.Fatalf("Error getting department: %v", err)
	}
	if department.HeadID != tu.TestInitDepartmentData.HeadID {
		t.Fatalf("Department not found: %v", department)
	}
}

func TestGetDepartmentUsersIDsIncorrectLondon(t *testing.T) {
	departmentRepository := tu.CreateTestDepartmentRepository()
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	departmentService := CreateDepartmentService(departmentRepository, moderatorRepository)
	_, err := departmentService.GetDepartmentUsersIDs(tu.TestInitDepartmentData.ID + 1)
	if err == nil {
		t.Fatalf("No error getting department users ids: %v", err)
	}
}
