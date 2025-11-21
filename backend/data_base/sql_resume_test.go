package data_base

import (
	tu "data_base_project/test_database_utility"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func setupResumeTables(db *sql.DB) error {
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
	err = CreateSqlRepetitorTable(db, "repetitors", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating repetitor table: %v", err)
	}
	err = CreateSqlResumeTable(db, "resume", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating resume table: %v", err)
	}
	return nil
}

func TestCreateSqlResumeRepositoryCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupResumeTables(db)
	if err != nil {
		t.Fatalf("Error setting up resume tables: %v", err)
	}
	ResumeRepository := CreateSqlResumeRepository(db, "resume", "sequence")
	if ResumeRepository == nil {
		t.Fatalf("Error creating resume repository: %v", err)
	}
}

func TestInsertResumeIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupResumeTables(db)
	if err != nil {
		t.Fatalf("Error setting up resume tables: %v", err)
	}
	ResumeRepository := CreateSqlResumeRepository(db, "resume", "sequence")
	_, err = ResumeRepository.InsertResume(tu.TestResume)
	if err == nil {
		t.Fatalf("No error inserting resume: %v", err)
	}
}

func TestInsertResumeInSeqIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupResumeTables(db)
	if err != nil {
		t.Fatalf("Error setting up resume tables: %v", err)
	}
	ResumeRepository := CreateSqlResumeRepository(db, "resume", "sequence")
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()
	_, err = ResumeRepository.InsertResumeInSeq(tx, tu.TestResume)
	if err == nil {
		t.Fatalf("No error inserting resume: %v", err)
	}
}

func TestGetResumeIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupResumeTables(db)
	if err != nil {
		t.Fatalf("Error setting up resume tables: %v", err)
	}
	ResumeRepository := CreateSqlResumeRepository(db, "resume", "sequence")
	_, err = ResumeRepository.GetResume(1)
	if err == nil {
		t.Fatalf("No error getting resume: %v", err)
	}
}

func TestUpdateResumeTitleCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupResumeTables(db)
	if err != nil {
		t.Fatalf("Error setting up resume tables: %v", err)
	}
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	tu.TestResume.RepetitorID = repetitorID
	tu.TestResume.RepetitorID = repetitorID
	ResumeRepository := CreateSqlResumeRepository(db, "resume", "sequence")
	resumeID, err := ResumeRepository.InsertResume(tu.TestResume)
	if err != nil {
		t.Fatalf("Error inserting resume: %v", err)
	}
	err = ResumeRepository.UpdateResumeTitle(resumeID, "new title", time.Now())
	if err != nil {
		t.Fatalf("Error updating resume title: %v", err)
	}
	resume, err := ResumeRepository.GetResume(resumeID)
	if err != nil {
		t.Fatalf("Error getting resume: %v", err)
	}
	if resume.Title != "new title" {
		t.Fatalf("Resume title not updated: %v", resume)
	}
}

func TestUpdateResumeDescriptionCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupResumeTables(db)
	if err != nil {
		t.Fatalf("Error setting up resume tables: %v", err)
	}
	ResumeRepository := CreateSqlResumeRepository(db, "resume", "sequence")
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	tu.TestResume.RepetitorID = repetitorID
	resumeID, err := ResumeRepository.InsertResume(tu.TestResume)
	if err != nil {
		t.Fatalf("Error inserting resume: %v", err)
	}
	err = ResumeRepository.UpdateResumeDescription(resumeID, "new description", time.Now())
	if err != nil {
		t.Fatalf("Error updating resume description: %v", err)
	}
	resume, err := ResumeRepository.GetResume(resumeID)
	if err != nil {
		t.Fatalf("Error getting resume: %v", err)
	}
	if resume.Description != "new description" {
		t.Fatalf("Resume description not updated: %v", resume)
	}
}

func TestUpdateResumeDescriptionIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupResumeTables(db)
	if err != nil {
		t.Fatalf("Error setting up resume tables: %v", err)
	}
	ResumeRepository := CreateSqlResumeRepository(db, "resume", "sequence")
	err = ResumeRepository.UpdateResumeDescription(1, "new description", time.Now())
	if err == nil {
		t.Fatalf("No error updating resume description: %v", err)
	}
}

func SetupTestDb(t *testing.T) *sql.DB {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	return db
}

func TestUpdateResumePricesCorrect(t *testing.T) {
	db := SetupTestDb(t)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err := setupResumeTables(db)
	if err != nil {
		t.Fatalf("Error setting up resume tables: %v", err)
	}
	ResumeRepository := CreateSqlResumeRepository(db, "resume", "sequence")
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	tu.TestResume.RepetitorID = repetitorID
	resumeID, err := ResumeRepository.InsertResume(tu.TestResume)
	if err != nil {
		t.Fatalf("Error inserting resume: %v", err)
	}
	newPrices := map[string]int{"Go": 200, "EBPF": 300}
	err = ResumeRepository.UpdateResumePrices(resumeID, newPrices, time.Now())
	if err != nil {
		t.Fatalf("Error updating resume prices: %v", err)
	}
	resume, err := ResumeRepository.GetResume(resumeID)
	if err != nil {
		t.Fatalf("Error getting resume: %v", err)
	}
	for key, value := range resume.Prices {
		if _, ok := newPrices[key]; !ok || value != newPrices[key] {
			t.Fatalf("Resume prices not updated: %v", resume)
		}
	}
}

func TestUpdateResumePricesIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupResumeTables(db)
	if err != nil {
		t.Fatalf("Error setting up resume tables: %v", err)
	}
	ResumeRepository := CreateSqlResumeRepository(db, "resume", "sequence")
	err = ResumeRepository.UpdateResumePrices(1, map[string]int{"Go": 200, "EBPF": 300}, time.Now())
	if err == nil {
		t.Fatalf("No error updating resume prices: %v", err)
	}
}

func TestDeleteResumeCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupResumeTables(db)
	if err != nil {
		t.Fatalf("Error setting up resume tables: %v", err)
	}
	ResumeRepository := CreateSqlResumeRepository(db, "resume", "sequence")
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	tu.TestResume.RepetitorID = repetitorID
	resumeID, err := ResumeRepository.InsertResume(tu.TestResume)
	if err != nil {
		t.Fatalf("Error inserting resume: %v", err)
	}
	err = ResumeRepository.DeleteResume(resumeID)
	if err != nil {
		t.Fatalf("Error deleting resume: %v", err)
	}
	resume, err := ResumeRepository.GetResume(resumeID)
	if err == nil {
		t.Fatalf("No error getting resume after deletion: %v", resume)
	}
}

func TestDeleteResumeIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupResumeTables(db)
	if err != nil {
		t.Fatalf("Error setting up resume tables: %v", err)
	}
	ResumeRepository := CreateSqlResumeRepository(db, "resume", "sequence")
	err = ResumeRepository.DeleteResume(1)
	if err == nil {
		t.Fatalf("No error deleting resume: %v", err)
	}
}
