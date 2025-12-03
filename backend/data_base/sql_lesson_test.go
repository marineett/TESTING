package data_base

import (
	tu "data_base_project/test_database_utility"
	"database/sql"
	"fmt"
	"testing"
)

func setupLessonTables(db *sql.DB) error {
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
	err = CreateSqlContractTable(db, "contract", "users", "review", "repetitor", "client")
	if err != nil {
		return fmt.Errorf("error creating contract table: %v", err)
	}
	err = CreateSqlLessonTable(db, "lesson", "contract", "transaction")
	if err != nil {
		return fmt.Errorf("error creating lesson table: %v", err)
	}
	err = CreateSqlClientTable(db, "clients", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating client table: %v", err)
	}
	return nil
}

func TestCreateSqlLessonRepositoryCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupLessonTables(db)
	if err != nil {
		t.Fatalf("Error setting up lesson tables: %v", err)
	}
	lessonRepository := CreateSqlLessonRepository(db, "lesson", "contract", "transaction", "sequence")
	if lessonRepository == nil {
		t.Fatalf("Error creating lesson repository: %v", err)
	}
}

func TestInsertLessonCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupLessonTables(db)
	if err != nil {
		t.Fatalf("Error setting up lesson tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contract", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	if clientRepository == nil {
		t.Fatalf("Error creating client repository: %v", err)
	}
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	tu.TestLesson.ContractID = contractID
	lessonRepository := CreateSqlLessonRepository(db, "lesson", "contract", "transaction", "sequence")
	if lessonRepository == nil {
		t.Fatalf("Error creating lesson repository: %v", err)
	}
	_, err = lessonRepository.InsertLesson(tu.TestLesson)
	if err != nil {
		t.Fatalf("Error inserting lesson: %v", err)
	}
}

func TestInsertLessonIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupLessonTables(db)
	if err != nil {
		t.Fatalf("Error setting up lesson tables: %v", err)
	}
	lessonRepository := CreateSqlLessonRepository(db, "lesson", "contract", "transaction", "sequence")
	if lessonRepository == nil {
		t.Fatalf("Error creating lesson repository: %v", err)
	}
	_, err = lessonRepository.InsertLesson(tu.TestLesson)
	if err == nil {
		t.Fatalf("No error inserting lesson: %v", err)
	}
}

func TestGetLessonsCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupLessonTables(db)
	if err != nil {
		t.Fatalf("Error setting up lesson tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contract", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	if clientRepository == nil {
		t.Fatalf("Error creating client repository: %v", err)
	}
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	tu.TestLesson.ContractID = 0
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	tu.TestLesson.ContractID = contractID
	lessonRepository := CreateSqlLessonRepository(db, "lesson", "contract", "transaction", "sequence")
	if lessonRepository == nil {
		t.Fatalf("Error creating lesson repository: %v", err)
	}
	_, err = lessonRepository.InsertLesson(tu.TestLesson)
	if err != nil {
		t.Fatalf("Error inserting lesson: %v", err)
	}
	_, err = lessonRepository.InsertLesson(tu.TestLesson)
	if err != nil {
		t.Fatalf("Error inserting lesson: %v", err)
	}
	lessons, err := lessonRepository.GetLessons(contractID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting lessons: %v", err)
	}
	if len(lessons) != 2 {
		t.Fatalf("Lessons list is not correct: %v", lessons)
	}
}
