package data_base

import (
	tu "data_base_project/test_database_utility"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
)

func setupAdminTables(db *sql.DB) error {
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
	err = CreateSqlAuthTable(db, "auth", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating auth table: %v", err)
	}
	err = CreateSqlAdminTable(db, "admins", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating admin table: %v", err)
	}
	return nil
}

func TestCreateSqlAdminRepositoryCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	if adminRepository == nil {
		t.Fatalf("Error creating admin repository: %v", err)
	}
}

func TestInsertAdminCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAdminTables(db)
	if err != nil {
		t.Fatalf("Error setting up admin tables: %v", err)
	}
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	adminRepository.InsertAdmin(tu.TestAdmin, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting admin: %v", err)
	}
}
func TestGetAdminCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAdminTables(db)
	if err != nil {
		t.Fatalf("Error setting up admin tables: %v", err)
	}
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	adminID, err := adminRepository.InsertAdmin(tu.TestAdmin, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting admin: %v", err)
	}
	admin, err := adminRepository.GetAdmin(adminID)
	if err != nil {
		t.Fatalf("Error getting admin: %v, adminID: %v", err, adminID)
	}
	if admin.ID != adminID {
		t.Fatalf("Admin not found: %v", admin)
	}
	if admin.Salary != tu.TestSalary {
		t.Fatalf("Admin not found: %v", admin)
	}
	if admin.DepartmentID != 0 {
		t.Fatalf("Admin not found: %v", admin)
	}
}

func TestGetAdminIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAdminTables(db)
	if err != nil {
		t.Fatalf("Error setting up admin tables: %v", err)
	}
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	_, err = adminRepository.GetAdmin(1)
	if err == nil {
		t.Fatalf("No error getting admin: %v", err)
	}
}

func TestUpdateAdminPersonalDataCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAdminTables(db)
	if err != nil {
		t.Fatalf("Error setting up admin tables: %v", err)
	}
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	adminID, err := adminRepository.InsertAdmin(tu.TestAdmin, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting admin: %v", err)
	}
	err = adminRepository.UpdateAdminPersonalData(adminID, tu.TestPD)
	if err != nil {
		t.Fatalf("Error updating admin personal data: %v", err)
	}
	admin, err := adminRepository.GetAdmin(adminID)
	if err != nil {
		t.Fatalf("Error getting admin: %v", err)
	}
	if admin.ID != adminID {
		t.Fatalf("Admin not found: %v", admin)
	}
	if admin.Salary != tu.TestSalary {
		t.Fatalf("Admin not found: %v", admin)
	}
	if admin.DepartmentID != 0 {
		t.Fatalf("Admin not found: %v", admin)
	}
}

func TestUpdateAdminPersonalDataIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAdminTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	err = adminRepository.UpdateAdminPersonalData(1, tu.TestPD)
	if err == nil {
		t.Fatalf("No error updating admin personal data: %v", err)
	}
}

func TestUpdateAdminPasswordCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAdminTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	adminID, err := adminRepository.InsertAdmin(tu.TestAdmin, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting admin: %v", err)
	}
	err = adminRepository.UpdateAdminPassword(adminID, tu.TestAuthData, "test3")
	if err != nil {
		t.Fatalf("Error updating admin password: %v", err)
	}
}

func TestUpdateAdminPasswordIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAdminTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	err = adminRepository.UpdateAdminPassword(1, tu.TestAuthData, "test3")
	if err == nil {
		t.Fatalf("No error updating admin password: %v", err)
	}
}

func TestUpdateAdminDepartmentCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAdminTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	adminID, err := adminRepository.InsertAdmin(tu.TestAdmin, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting admin: %v", err)
	}
	err = adminRepository.UpdateAdminDepartment(adminID, 10)
	if err != nil {
		t.Fatalf("Error updating admin department: %v", err)
	}
	admin, err := adminRepository.GetAdmin(adminID)
	if err != nil {
		t.Fatalf("Error getting admin: %v", err)
	}
	if admin.ID != adminID {
		t.Fatalf("Admin not found: %v", admin)
	}
	if admin.DepartmentID != 10 {
		t.Fatalf("Admin not found: %v", admin)
	}
	if admin.Salary != tu.TestSalary {
		t.Fatalf("Admin not found: %v", admin)
	}
}

func TestUpdateAdminDepartmentIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAdminTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	err = adminRepository.UpdateAdminDepartment(1, 10)
	if err == nil {
		t.Fatalf("No error updating admin department: %v", err)
	}
}

func TestUpdateAdminSalaryCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAdminTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	err = CreateSqlUserTable(db, "users", "personal_data", "sequence")
	if err != nil {
		t.Fatalf("Error creating user table: %v", err)
	}
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	adminID, err := adminRepository.InsertAdmin(tu.TestAdmin, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting admin: %v", err)
	}
	err = adminRepository.UpdateAdminSalary(adminID, 100000)
	if err != nil {
		t.Fatalf("Error updating admin salary: %v", err)
	}
	admin, err := adminRepository.GetAdmin(adminID)
	if err != nil {
		t.Fatalf("Error getting admin: %v", err)
	}
	if admin.ID != adminID {
		t.Fatalf("Admin not found: %v", admin)
	}
	if admin.Salary != 100000 {
		t.Fatalf("Admin not found: %v", admin)
	}
	if admin.DepartmentID != 0 {
		t.Fatalf("Admin not found: %v", admin)
	}
}

func TestUpdateAdminSalaryIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAdminTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	adminRepository := CreateSqlAdminRepository(db, "personal_data", "users", "admins", "auth", "sequence")
	err = adminRepository.UpdateAdminSalary(1, 100000)
	if err == nil {
		t.Fatalf("No error updating admin salary: %v", err)
	}
}
