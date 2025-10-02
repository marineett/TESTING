package data_base

import (
	tu "data_base_project/test_database_utility"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
)

func setupUserTables(db *sql.DB) error {
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
	return nil
}

func TestCreateSqlUserTable(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupUserTables(db)
	if err != nil {
		t.Fatalf("Error setting up user tables: %v", err)
	}
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	if userRepository == nil {
		t.Fatalf("Error creating user repository: %v", err)
	}
}

func TestInsertUserCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupUserTables(db)
	if err != nil {
		t.Fatalf("Error setting up user tables: %v", err)
	}
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	userRepository.InsertUser(tu.TestUser)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
}

func TestInsertUserIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupUserTables(db)
	if err != nil {
		t.Fatalf("Error setting up user tables: %v", err)
	}
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	_, err = userRepository.InsertUser(tu.TestUser)
	if err == nil {
		t.Fatalf("No error inserting user: %v", err)
	}
}

func TestGetUserCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupUserTables(db)
	if err != nil {
		t.Fatalf("Error setting up user tables: %v", err)
	}
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataID, err := personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	tu.TestUser.PersonalDataID = personalDataID
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	userID, err := userRepository.InsertUser(tu.TestUser)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	user, err := userRepository.GetUser(userID)
	if err != nil {
		t.Fatalf("Error getting user: %v", err)
	}
	if user.ID != userID {
		t.Fatalf("User id not updated: %v", user)
	}
	if user.PersonalDataID != tu.TestUser.PersonalDataID {
		t.Fatalf("User personal data id not updated: %v", user)
	}
}

func TestGetUserIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupUserTables(db)
	if err != nil {
		t.Fatalf("Error setting up user tables: %v", err)
	}
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	_, err = userRepository.GetUser(1)
	if err == nil {
		t.Fatalf("No error getting not existing user: %v", err)
	}
}

func TestInsertUserInSeqCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupUserTables(db)
	if err != nil {
		t.Fatalf("Error setting up user tables: %v", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	defer tx.Rollback()
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataID, err := personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	tu.TestUser.PersonalDataID = personalDataID
	_, err = userRepository.InsertUserInSeq(tx, tu.TestUser)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
}

func TestInsertUserInSeqIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupUserTables(db)
	if err != nil {
		t.Fatalf("Error setting up user tables: %v", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	defer tx.Rollback()
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	_, err = userRepository.InsertUserInSeq(tx, tu.TestUser)
	if err == nil {
		t.Fatalf("No error inserting user: %v", err)
	}
}
