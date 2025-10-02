package data_base

import (
	tu "data_base_project/test_database_utility"
	"data_base_project/types"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func setupPersonalDataTables(db *sql.DB) error {
	err := CreateSqlSequence(db, "sequence")
	if err != nil {
		return fmt.Errorf("error creating sequence: %v", err)
	}
	err = CreateSqlPersonalDataTable(db, "personal_data", "sequence")
	if err != nil {
		return fmt.Errorf("error creating personal data table: %v", err)
	}
	return nil
}

func TestInsertPersonalDataCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupPersonalDataTables(db)
	if err != nil {
		t.Fatalf("Error setting up personal data tables: %v", err)
	}
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
}

func TestInsertPersonalDataInSeqCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupPersonalDataTables(db)
	if err != nil {
		t.Fatalf("Error setting up personal data tables: %v", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	defer tx.Rollback()
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataRepository.InsertPersonalDataInSeq(tx, tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
}

func TestGetPersonalDataCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupPersonalDataTables(db)
	if err != nil {
		t.Fatalf("Error setting up personal data tables: %v", err)
	}
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataID, err := personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	personalData, err := personalDataRepository.GetPersonalData(personalDataID)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalData.ID != personalDataID {
		t.Fatalf("Personal data id not updated: %v", personalData)
	}
	if personalData.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Personal data telephone number not updated: %v", personalData)
	}
	if personalData.Email != tu.TestPD.Email {
		t.Fatalf("Personal data email not updated: %v", personalData)
	}
	if personalData.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Personal data first name not updated: %v", personalData)
	}
	if personalData.LastName != tu.TestPD.LastName {
		t.Fatalf("Personal data last name not updated: %v", personalData)
	}
	if personalData.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Personal data middle name not updated: %v", personalData)
	}
	if personalData.PassportNumber != tu.TestPD.PassportNumber {
		t.Fatalf("Personal data passport number not updated: %v", personalData)
	}
	if personalData.PassportSeries != tu.TestPD.PassportSeries {
		t.Fatalf("Personal data passport series not updated: %v", personalData)
	}
	if personalData.PassportIssuedBy != tu.TestPD.PassportIssuedBy {
		t.Fatalf("Personal data passport issued by not updated: %v", personalData)
	}
}

func TestGetPersonalDataIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupPersonalDataTables(db)
	if err != nil {
		t.Fatalf("Error setting up personal data tables: %v", err)
	}
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	_, err = personalDataRepository.GetPersonalData(1)
	if err == nil {
		t.Fatalf("No error getting personal data: %v", err)
	}
}

func TestUpdatePersonalDataCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupPersonalDataTables(db)
	if err != nil {
		t.Fatalf("Error setting up personal data tables: %v", err)
	}
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataID, err := personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	newPassportData := types.DBPassportData{
		PassportNumber:   "1234567890",
		PassportSeries:   "1234",
		PassportDate:     time.Now(),
		PassportIssuedBy: "Moscow",
	}
	newPersonalData := types.DBPersonalData{
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		DBPassportData:  newPassportData,
	}
	err = personalDataRepository.UpdatePersonalData(personalDataID, newPersonalData)
	if err != nil {
		t.Fatalf("Error updating personal data: %v", err)
	}
	personalData, err := personalDataRepository.GetPersonalData(personalDataID)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalData.ID != personalDataID {
		t.Fatalf("Personal data id not updated: %v", personalData)
	}
	if personalData.TelephoneNumber != newPersonalData.TelephoneNumber {
		t.Fatalf("Personal data telephone number not updated: %v", personalData)
	}
	if personalData.Email != newPersonalData.Email {
		t.Fatalf("Personal data email not updated: %v", personalData)
	}
	if personalData.FirstName != newPersonalData.FirstName {
		t.Fatalf("Personal data first name not updated: %v", personalData)
	}
	if personalData.LastName != newPersonalData.LastName {
		t.Fatalf("Personal data last name not updated: %v", personalData)
	}
	if personalData.MiddleName != newPersonalData.MiddleName {
		t.Fatalf("Personal data middle name not updated: %v", personalData)
	}
	if personalData.PassportNumber != newPersonalData.PassportNumber {
		t.Fatalf("Personal data passport number not updated: %v", personalData)
	}
	if personalData.PassportSeries != newPersonalData.PassportSeries {
		t.Fatalf("Personal data passport series not updated: %v", personalData)
	}
	if personalData.PassportIssuedBy != newPersonalData.PassportIssuedBy {
		t.Fatalf("Personal data passport issued by not updated: %v", personalData)
	}
}

func TestUpdatePersonalDataIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupPersonalDataTables(db)
	if err != nil {
		t.Fatalf("Error setting up personal data tables: %v", err)
	}
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	err = personalDataRepository.UpdatePersonalData(1, tu.TestPD)
	if err == nil {
		t.Fatalf("No error updating personal data: %v", err)
	}
}
