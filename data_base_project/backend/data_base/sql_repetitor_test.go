package data_base

import (
	tu "data_base_project/test_database_utility"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
)

func setupRepetitorTables(db *sql.DB) error {
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
	err = CreateSqlResumeTable(db, "resume", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating resume table: %v", err)
	}
	err = CreateSqlReviewTable(db, "review", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating review table: %v", err)
	}
	err = CreateSqlRepetitorTable(db, "repetitors", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating repetitor table: %v", err)
	}
	return nil
}

func TestCreateSqlRepetitorRepositoryCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	RepetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	if RepetitorRepository == nil {
		t.Fatalf("Error creating repetitor repository: %v", err)
	}
}

func TestInsertRepetitorCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupRepetitorTables(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	RepetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	RepetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Repetitor: %v", err)
	}
}
func TestGetRepetitorCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupRepetitorTables(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	repetitor, err := repetitorRepository.GetRepetitor(repetitorID)
	if err != nil {
		t.Fatalf("Error getting repetitor: %v, repetitorID: %v", err, repetitorID)
	}
	if repetitor.ID != repetitorID {
		t.Fatalf("repetitor not found: %v", repetitor)
	}
	if repetitor.SummaryRating != tu.TestRepetitor.SummaryRating {
		t.Fatalf("repetitor not found: %v", repetitor)
	}
	if repetitor.ReviewsCount != tu.TestRepetitor.ReviewsCount {
		t.Fatalf("repetitor not found: %v", repetitor)
	}
}

func TestGetRepetitorIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupRepetitorTables(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	RepetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	_, err = RepetitorRepository.GetRepetitor(1)
	if err == nil {
		t.Fatalf("No error getting repetitor: %v", err)
	}
}

func TestUpdateRepetitorPersonalDataCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupRepetitorTables(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	err = repetitorRepository.UpdateRepetitorPersonalData(repetitorID, tu.TestPD)
	if err != nil {
		t.Fatalf("Error updating repetitor personal data: %v", err)
	}
	repetitor, err := repetitorRepository.GetRepetitor(repetitorID)
	if err != nil {
		t.Fatalf("Error getting repetitor: %v", err)
	}
	if repetitor.ID != repetitorID {
		t.Fatalf("repetitor not found: %v", repetitor)
	}
	if repetitor.SummaryRating != tu.TestRepetitor.SummaryRating {
		t.Fatalf("repetitor not found: %v", repetitor)
	}
	if repetitor.ReviewsCount != tu.TestRepetitor.ReviewsCount {
		t.Fatalf("repetitor not found: %v", repetitor)
	}
}

func TestUpdateRepetitorPersonalDataIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupRepetitorTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	err = repetitorRepository.UpdateRepetitorPersonalData(1, tu.TestPD)
	if err == nil {
		t.Fatalf("No error updating Repetitor personal data: %v", err)
	}
}

func TestUpdateRepetitorPasswordCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupRepetitorTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Client: %v", err)
	}
	err = repetitorRepository.UpdateRepetitorPassword(repetitorID, tu.TestAuthData, "test3")
	if err != nil {
		t.Fatalf("Error updating Repetitor password: %v", err)
	}
}

func TestUpdateRepetitorPasswordIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupRepetitorTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	err = repetitorRepository.UpdateRepetitorPassword(1, tu.TestAuthData, "test3")
	if err == nil {
		t.Fatalf("No error updating Repetitor password: %v", err)
	}
}

func TestGetRepetitorsIdsCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupRepetitorTables(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	_, err = repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	_, err = repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	_, err = repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	_, err = repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	repetitorsIds, err := repetitorRepository.GetRepetitorsIds(0, 10)
	if err != nil {
		t.Fatalf("Error getting repetitors ids: %v", err)
	}
	if len(repetitorsIds) != 4 {
		t.Fatalf("Error getting repetitors ids: %v", repetitorsIds)
	}
}
