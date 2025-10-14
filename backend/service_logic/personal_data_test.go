package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
)

func TestGetPersonalDataLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	testPD := types.MapperPersonalDataServiceToDB(&tu.TestPD)
	testPDID, err := personalDataRepository.InsertPersonalData(*testPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	personalDataService := CreatePersonalDataService(personalDataRepository)
	personalDataServiceData, err := personalDataService.GetPersonalData(testPDID)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalDataServiceData.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
	if personalDataServiceData.Email != tu.TestPD.Email {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
	if personalDataServiceData.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
	if personalDataServiceData.LastName != tu.TestPD.LastName {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
	if personalDataServiceData.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
	if personalDataServiceData.ServicePassportData.PassportNumber != tu.TestPD.ServicePassportData.PassportNumber {
		t.Fatalf("Personal data not updated: %s", personalDataServiceData.ServicePassportData.PassportNumber)
	}
	if personalDataServiceData.ServicePassportData.PassportSeries != tu.TestPD.ServicePassportData.PassportSeries {
		t.Fatalf("Personal data not updated: %s", personalDataServiceData.ServicePassportData.PassportSeries)
	}
	if personalDataServiceData.ServicePassportData.PassportIssuedBy != tu.TestPD.ServicePassportData.PassportIssuedBy {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
}

func TestGetPersonalDataClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up personal data tables: %v", err)
	}
	personalDataRepository := module.PersonalDataRepository
	testPD := types.MapperPersonalDataServiceToDB(&tu.TestPD)
	personalDataService := CreatePersonalDataService(personalDataRepository)
	personalDataID, err := personalDataRepository.InsertPersonalData(*testPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	personalDataServiceData, err := personalDataService.GetPersonalData(personalDataID)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalDataServiceData.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
	if personalDataServiceData.Email != tu.TestPD.Email {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
	if personalDataServiceData.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
	if personalDataServiceData.LastName != tu.TestPD.LastName {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
	if personalDataServiceData.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
	if personalDataServiceData.ServicePassportData.PassportNumber != tu.TestPD.ServicePassportData.PassportNumber {
		t.Fatalf("Personal data not updated: %s", personalDataServiceData.ServicePassportData.PassportNumber)
	}
	if personalDataServiceData.ServicePassportData.PassportSeries != tu.TestPD.ServicePassportData.PassportSeries {
		t.Fatalf("Personal data not updated: %s", personalDataServiceData.ServicePassportData.PassportSeries)
	}
	if personalDataServiceData.ServicePassportData.PassportIssuedBy != tu.TestPD.ServicePassportData.PassportIssuedBy {
		t.Fatalf("Personal data not updated: %v", personalDataServiceData)
	}
}

func TestGetPersonalDataIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	testPD := types.MapperPersonalDataServiceToDB(&tu.TestPD)
	testPDID, err := personalDataRepository.InsertPersonalData(*testPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	personalDataService := CreatePersonalDataService(personalDataRepository)
	personalDataServiceData, err := personalDataService.GetPersonalData(testPDID + 1)
	if err == nil {
		t.Fatalf("No error getting personal data: %v", personalDataServiceData)
	}
}

func TestGetPersonalDataIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up personal data tables: %v", err)
	}
	personalDataRepository := module.PersonalDataRepository
	personalDataService := CreatePersonalDataService(personalDataRepository)
	personalDataServiceData, err := personalDataService.GetPersonalData(1)
	if err == nil {
		t.Fatalf("No error getting personal data: %v", personalDataServiceData)
	}
}
