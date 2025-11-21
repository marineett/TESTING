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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupPersonalDataTables(db)
	if err != nil {
		t.Fatalf("Error setting up personal data tables: %v", err)
	}
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	_, err = personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
}

func TestInsertPersonalDataInSeqCorrect(t *testing.T) {
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
	err = setupPersonalDataTables(db)
	if err != nil {
		t.Fatalf("Error setting up personal data tables: %v", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	_, err = personalDataRepository.InsertPersonalDataInSeq(tx, tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
}

func CheckPersonalData(
	t *testing.T,
	personalData *types.DBPersonalData,
	personalDataID int64,
	telephoneNumber string,
	email string,
	firstName string,
	lastName string,
	middleName string,
	passportNumber string,
	passportSeries string,
	passportIssuedBy string,
) {
	if personalData.ID != personalDataID {
		t.Fatalf("Personal data id not updated: %v", personalData)
	}
	if personalData.TelephoneNumber != telephoneNumber {
		t.Fatalf("Personal data telephone number not updated: %v", personalData)
	}
	if personalData.Email != email {
		t.Fatalf("Personal data email not updated: %v", personalData)
	}
	if personalData.FirstName != firstName {
		t.Fatalf("Personal data first name not updated: %v", personalData)
	}
	if personalData.LastName != lastName {
		t.Fatalf("Personal data last name not updated: %v", personalData)
	}
	if personalData.MiddleName != middleName {
		t.Fatalf("Personal data middle name not updated: %v", personalData)
	}
	if personalData.PassportNumber != passportNumber {
		t.Fatalf("Personal data passport number not updated: %v", personalData)
	}
	if personalData.PassportSeries != passportSeries {
		t.Fatalf("Personal data passport series not updated: %v", personalData)
	}
	if personalData.PassportIssuedBy != passportIssuedBy {
		t.Fatalf("Personal data passport issued by not updated: %v", personalData)
	}
}

func TestGetPersonalDataCorrect(t *testing.T) {
	db := SetupDatabase(t)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err := setupPersonalDataTables(db)
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
	CheckPersonalData(t, personalData, personalDataID, tu.TestPD.TelephoneNumber, tu.TestPD.Email, tu.TestPD.FirstName, tu.TestPD.LastName, tu.TestPD.MiddleName, tu.TestPD.PassportNumber, tu.TestPD.PassportSeries, tu.TestPD.PassportIssuedBy)
}

func TestGetPersonalDataIncorrect(t *testing.T) {
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
	db := SetupDatabase(t)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err := setupPersonalDataTables(db)
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
	CheckPersonalData(
		t,
		personalData,
		personalDataID,
		newPersonalData.TelephoneNumber,
		newPersonalData.Email,
		newPersonalData.FirstName,
		newPersonalData.LastName,
		newPersonalData.MiddleName,
		newPersonalData.PassportNumber,
		newPersonalData.PassportSeries,
		newPersonalData.PassportIssuedBy,
	)
}

func TestUpdatePersonalDataIncorrect(t *testing.T) {
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
