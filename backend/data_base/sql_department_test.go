package data_base

import (
	tu "data_base_project/test_database_utility"
	"database/sql"
	"fmt"
	"testing"
)

func setupDepartmentTables(db *sql.DB) error {
	err := CreateSqlSequence(db, "sequence")
	if err != nil {
		return fmt.Errorf("error creating sequence: %v", err)
	}
	err = CreateSqlPersonalDataTable(db, "personal_data", "sequence")
	if err != nil {
		return fmt.Errorf("error creating personal data table: %v", err)
	}
	err = CreateSqlUserTable(db, "users", "personal_data", "sequence")
	if err != nil {
		return fmt.Errorf("error creating user table: %v", err)
	}
	err = CreateSqlDepartmentTable(db, "department", "hire_info", "users")
	if err != nil {
		return fmt.Errorf("error creating department table: %v", err)
	}
	err = CreateSqlModeratorTable(db, "moderators", "users")
	if err != nil {
		return fmt.Errorf("error creating moderator table: %v", err)
	}
	err = CreateSqlAuthTable(db, "auth", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating auth table: %v", err)
	}
	return nil
}

func TestCreateSqlDepartmentTable(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
}

func TestInsertDepartmentCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	_, err = departmentRepository.InsertDepartment(tu.TestDepartment)
	if err != nil {
		t.Fatalf("Error inserting department: %v", err)
	}
}

func TestGetDepartmentCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	departmentID, err := departmentRepository.InsertDepartment(tu.TestDepartment)
	if err != nil {
		t.Fatalf("Error inserting department: %v", err)
	}
	department, err := departmentRepository.GetDepartment(departmentID)
	if err != nil {
		t.Fatalf("Error getting department: %v", err)
	}
	if department.Name != tu.TestDepartment.Name {
		t.Fatalf("Department not found: %v", department)
	}
	if department.HeadID != tu.TestDepartment.HeadID {
		t.Fatalf("Department not found: %v", department)
	}
}

func TestGetDepartmentIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	_, err = departmentRepository.GetDepartment(1)
	if err == nil {
		t.Fatalf("No error getting department: %v", err)
	}
}

func TestGetDepartmentsByHeadIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	tu.TestDepartment.HeadID = 1
	_, err = departmentRepository.InsertDepartment(tu.TestDepartment)
	if err != nil {
		t.Fatalf("Error inserting department: %v", err)
	}
	_, err = departmentRepository.InsertDepartment(tu.TestDepartment)
	if err != nil {
		t.Fatalf("Error inserting department: %v", err)
	}
	tu.TestDepartment.HeadID = 2
	_, err = departmentRepository.InsertDepartment(tu.TestDepartment)
	if err != nil {
		t.Fatalf("Error inserting department: %v", err)
	}
	departments, err := departmentRepository.GetDepartmentsByHeadID(1)
	if err != nil {
		t.Fatalf("Error getting departments by head id: %v", err)
	}
	if len(departments) != 2 {
		t.Fatalf("Departments not found: %v", departments)
	}
}

func TestGetDepartmentIDByNameCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	_, err = departmentRepository.InsertDepartment(tu.TestDepartment)
	if err != nil {
		t.Fatalf("Error inserting department: %v", err)
	}
	_, err = departmentRepository.GetDepartmentIdByName(tu.TestDepartment.Name)
	if err != nil {
		t.Fatalf("Error getting department id by name: %v", err)
	}
}

func TestGetDepartmentIDByNameIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	_, err = departmentRepository.GetDepartmentIdByName("Test Department")
	if err == nil {
		t.Fatalf("No error getting department id by name: %v", err)
	}
}

func TestChangeDepartmentHeadCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	_, err = departmentRepository.InsertDepartment(tu.TestDepartment)
	if err != nil {
		t.Fatalf("Error inserting department: %v", err)
	}
	err = departmentRepository.ChangeDepartmentHead(1, 2)
	if err != nil {
		t.Fatalf("Error changing department head: %v", err)
	}
	department, err := departmentRepository.GetDepartment(1)
	if err != nil {
		t.Fatalf("Error getting department: %v", err)
	}
	if department.HeadID != 2 {
		t.Fatalf("Department head not found: %v", department)
	}
}

func TestChangeDepartmentHeadIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	err = departmentRepository.ChangeDepartmentHead(1, 2)
	if err == nil {
		t.Fatalf("No error changing department head: %v", err)
	}
}

func TestHireInfoInsertCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	moderatorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	if moderatorRepository == nil {
		t.Fatalf("Error creating moderator repository: %v", err)
	}
	moderatorID, err := moderatorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting moderator: %v", err)
	}
	departmentID, err := departmentRepository.InsertDepartment(tu.TestDepartment)
	if err != nil {
		t.Fatalf("Error inserting department: %v", err)
	}
	tu.TestHireInfo.UserID = moderatorID
	tu.TestHireInfo.DepartmentID = departmentID
	err = departmentRepository.HireInfoInsert(tu.TestHireInfo)
	if err != nil {
		t.Fatalf("Error inserting hire info: %v", err)
	}
}

func TestHireInfoDeleteCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	tu.TestHireInfo.UserID = 1
	tu.TestHireInfo.DepartmentID = 1
	err = departmentRepository.HireInfoInsert(tu.TestHireInfo)
	if err != nil {
		t.Fatalf("Error inserting hire info: %v", err)
	}
	err = departmentRepository.HireInfoDelete(1, 1)
	if err != nil {
		t.Fatalf("Error deleting hire info: %v", err)
	}
}

func TestHireInfoDeleteIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	err = departmentRepository.HireInfoDelete(1, 1)
	if err == nil {
		t.Fatalf("No error deleting hire info: %v", err)
	}
}

func TestGetUserDepartmentsIDsCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	tu.TestHireInfo.UserID = 1
	tu.TestHireInfo.DepartmentID = 1
	err = departmentRepository.HireInfoInsert(tu.TestHireInfo)
	if err != nil {
		t.Fatalf("Error inserting hire info: %v", err)
	}
	tu.TestHireInfo.DepartmentID = 2
	err = departmentRepository.HireInfoInsert(tu.TestHireInfo)
	if err != nil {
		t.Fatalf("Error inserting hire info: %v", err)
	}
	tu.TestHireInfo.UserID = 3
	tu.TestHireInfo.DepartmentID = 3
	err = departmentRepository.HireInfoInsert(tu.TestHireInfo)
	if err != nil {
		t.Fatalf("Error inserting hire info: %v", err)
	}
	userDepartmentsIDs, err := departmentRepository.GetUserDepartmentsIDs(1)
	if err != nil {
		t.Fatalf("Error getting user departments ids: %v", err)
	}
	if len(userDepartmentsIDs) != 2 {
		t.Fatalf("User departments ids not found: %v", userDepartmentsIDs)
	}
}

func TestGetDepartmentUsersIDsCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupDepartmentTables(db)
	if err != nil {
		t.Fatalf("Error setting up department tables: %v", err)
	}
	departmentRepository := CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	if departmentRepository == nil {
		t.Fatalf("Error creating department repository: %v", err)
	}
	tu.TestHireInfo.UserID = 1
	tu.TestHireInfo.DepartmentID = 1
	err = departmentRepository.HireInfoInsert(tu.TestHireInfo)
	if err != nil {
		t.Fatalf("Error inserting hire info: %v", err)
	}
	tu.TestHireInfo.UserID = 2
	err = departmentRepository.HireInfoInsert(tu.TestHireInfo)
	if err != nil {
		t.Fatalf("Error inserting hire info: %v", err)
	}
	tu.TestHireInfo.UserID = 3
	tu.TestHireInfo.DepartmentID = 3
	err = departmentRepository.HireInfoInsert(tu.TestHireInfo)
	if err != nil {
		t.Fatalf("Error inserting hire info: %v", err)
	}
	userDepartmentsIDs, err := departmentRepository.GetDepartmentUsersIDs(1)
	if err != nil {
		t.Fatalf("Error getting department users ids: %v", err)
	}
	if len(userDepartmentsIDs) != 2 {
		t.Fatalf("Department users ids not found: %v", userDepartmentsIDs)
	}
}
