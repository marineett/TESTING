package data_base

import (
	tu "data_base_project/test_database_utility"
	"data_base_project/types"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
)

func setupAuthTables(db *sql.DB) error {
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
	return nil
}

func TestCreateSqlAuthRepositoryCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	authRepository := CreateSqlAuthRepository(db, "auth", "sequence")
	if authRepository == nil {
		t.Fatalf("Error creating auth repository: %v", err)
	}
}

func TestInsertAuthCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAuthTables(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := CreateSqlAuthRepository(db, "auth", "sequence")
	authRepository.InsertAuth(tu.TestAuthInfo)
	if err != nil {
		t.Fatalf("Error inserting auth: %v", err)
	}
}

func TestChangePasswordCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAuthTables(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := CreateSqlAuthRepository(db, "auth", "sequence")
	authRepository.InsertAuth(tu.TestAuthInfo)
	if err != nil {
		t.Fatalf("Error inserting auth: %v", err)
	}
	authRepository.ChangePassword(1, tu.TestAuthData, "newpassword")
	if err != nil {
		t.Fatalf("Error changing password: %v", err)
	}
}

func TestChangePasswordIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAuthTables(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := CreateSqlAuthRepository(db, "auth", "sequence")
	authRepository.InsertAuth(tu.TestAuthInfo)
	if err != nil {
		t.Fatalf("Error inserting auth: %v", err)
	}
	wrongAuthData := tu.TestAuthData
	wrongAuthData.Login = "wronglogin"
	err = authRepository.ChangePassword(2, wrongAuthData, "newpassword")
	if err == nil {
		t.Fatalf("No error changing password: %v", err)
	}
}

func TestAuthorizeCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAuthTables(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := CreateSqlAuthRepository(db, "auth", "sequence")
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	_, err = authRepository.InsertAuth(tu.TestAuthInfo)
	if err != nil {
		t.Fatalf("Error inserting auth: %v", err)
	}
	verdict, err := authRepository.Authorize(tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	if verdict.UserID != 1 {
		t.Fatalf("User id not updated: %v", verdict)
	}
	if verdict.UserType != types.Admin {
		t.Fatalf("User type not updated: %v", verdict)
	}
}

func TestAuthorizeIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAuthTables(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := CreateSqlAuthRepository(db, "auth", "sequence")
	authRepository.InsertAuth(tu.TestAuthInfo)
	if err != nil {
		t.Fatalf("Error inserting auth: %v", err)
	}
	wrongAuthData := tu.TestAuthData
	wrongAuthData.Password = "wrongpassword"
	_, err = authRepository.Authorize(wrongAuthData)
	if err == nil {
		t.Fatalf("No error authorizing: %v", err)
	}
}

func TestCheckLoginCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAuthTables(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := CreateSqlAuthRepository(db, "auth", "sequence")
	authRepository.InsertAuth(tu.TestAuthInfo)
	if err != nil {
		t.Fatalf("Error inserting auth: %v", err)
	}
	loginExists, err := authRepository.CheckLogin(tu.TestAuthData.Login)
	if err != nil {
		t.Fatalf("Error checking login: %v", err)
	}
	if !loginExists {
		t.Fatalf("Login not found: %v", loginExists)
	}
}

func TestCheckLoginIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupAuthTables(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := CreateSqlAuthRepository(db, "auth", "sequence")
	authRepository.InsertAuth(tu.TestAuthInfo)
	if err != nil {
		t.Fatalf("Error inserting auth: %v", err)
	}
	loginExists, err := authRepository.CheckLogin("incorrect")
	if err != nil {
		t.Fatalf("Error checking login: %v", err)
	}
	if loginExists {
		t.Fatalf("Login found: %v", loginExists)
	}
}
