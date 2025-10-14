package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
)

func TestCreateClientCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	clientRepository := tu.CreateTestClientRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err := clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	personalData, err := personalDataRepository.GetPersonalData(1)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalData.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Personal data not updated: %v", personalData)
	}
	if personalData.Email != tu.TestPD.Email {
		t.Fatalf("Personal data not updated: %v", personalData)
	}
	authData, err := authRepository.TestGetAuth(1)
	if err != nil {
		t.Fatalf("Error getting auth data: %v", err)
	}
	if authData.Login != tu.TestAuth.Login {
		t.Fatalf("Auth data not updated: %v", authData)
	}
	if authData.Password != tu.TestAuth.Password {
		t.Fatalf("Auth data not updated: %v", authData)
	}
	clientData, err := clientRepository.GetClient(1)
	if err != nil {
		t.Fatalf("Error getting client data: %v", err)
	}
	if clientData.SummaryRating != tu.TestSummaryRating {
		t.Fatalf("Client data not updated: %v", clientData)
	}
}

func TestCreateClientCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up Client tables: %v", err)
	}
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	clientData, err := clientRepository.GetClient(result.UserID)
	if err != nil {
		t.Fatalf("Error getting client data: %v", err)
	}
	if clientData.SummaryRating != tu.TestSummaryRating {
		t.Fatalf("Client data not updated: %v", clientData)
	}
}

func TestGetClientDataCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	clientRepository := tu.CreateTestClientRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err := clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	clientData, err := clientService.GetClientProfile(1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting client data: %v", err)
	}
	if clientData.MeanRating != tu.TestMeanRating {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.MeanRating != tu.TestMeanRating {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.LastName != tu.TestPD.LastName {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.Email != tu.TestPD.Email {
		t.Fatalf("Client data not updated: %v", clientData)
	}
}

func TestGetClientDataCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up Client tables: %v", err)
	}
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing client: %v", err)
	}
	clientData, err := clientService.GetClientProfile(result.UserID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting client data: %v", err)
	}
	if clientData.MeanRating != tu.TestMeanRating {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.MeanRating != tu.TestMeanRating {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.LastName != tu.TestPD.LastName {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Client data not updated: %v", clientData)
	}
	if clientData.Email != tu.TestPD.Email {
		t.Fatalf("Client data not updated: %v", clientData)
	}
}
func TestGetClientDataIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	clientRepository := tu.CreateTestClientRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err := clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	clientData, err := clientService.GetClientProfile(2, 0, 10)
	if err == nil {
		t.Fatalf("No error getting client data: %v", clientData)
	}
}

func TestGetClientDataIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up Client tables: %v", err)
	}
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	clientData, err := clientService.GetClientProfile(1, 0, 10)
	if err == nil {
		t.Fatalf("No error getting client data: %v", clientData)
	}
}
func TestUpdateClientPersonalDataCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	clientRepository := tu.CreateTestClientRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err := clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	newPersonalData := types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	}
	err = clientService.UpdateClientPersonalData(1, newPersonalData)
	if err != nil {
		t.Fatalf("Error updating client personal data: %v", err)
	}
	personalData, err := personalDataRepository.GetPersonalData(1)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalData.FirstName != newPersonalData.FirstName {
		t.Fatalf("Client personal data not updated: %v", personalData)
	}
	if personalData.LastName != newPersonalData.LastName {
		t.Fatalf("Client personal data not updated: %v", personalData)
	}
	if personalData.MiddleName != newPersonalData.MiddleName {
		t.Fatalf("Client personal data not updated: %v", personalData)
	}
	if personalData.TelephoneNumber != newPersonalData.TelephoneNumber {
		t.Fatalf("Client personal data not updated: %v", personalData)
	}
	if personalData.Email != newPersonalData.Email {
		t.Fatalf("Client personal data not updated: %v", personalData)
	}
}

func TestUpdateClientPersonalDataCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up Client tables: %v", err)
	}
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing client: %v", err)
	}
	err = clientService.UpdateClientPersonalData(result.UserID, types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	})
	if err != nil {
		t.Fatalf("Error updating client personal data: %v", err)
	}
	personalData, err := personalDataRepository.GetPersonalData(1)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalData.FirstName != "Petr" {
		t.Fatalf("Client personal data not updated: %v", personalData)
	}
	if personalData.LastName != "Petrov" {
		t.Fatalf("Client personal data not updated: %v", personalData)
	}
	if personalData.MiddleName != "Petrovich" {
		t.Fatalf("Client personal data not updated: %v", personalData)
	}
	if personalData.TelephoneNumber != "88005553536" {
		t.Fatalf("Client personal data not updated: %v", personalData)
	}
	if personalData.Email != "test2@test.com" {
		t.Fatalf("Client personal data not updated: %v", personalData)
	}
}

func TestUpdateClientPersonalDataIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	clientRepository := tu.CreateTestClientRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err := clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	err = clientService.UpdateClientPersonalData(2, types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	})
	if err == nil {
		t.Fatalf("No error updating client personal data: %v", err)
	}
}

func TestUpdateClientPersonalDataIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up Client tables: %v", err)
	}
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	err = clientService.UpdateClientPersonalData(1, types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	})
	if err == nil {
		t.Fatalf("No error updating client personal data: %v", err)
	}
}

func TestUpdateClientPasswordCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	clientRepository := tu.CreateTestClientRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err := clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	newPassword := "test3"
	err = clientService.UpdateClientPassword(1, tu.TestAuth, newPassword)
	if err != nil {
		t.Fatalf("Error updating client password: %v", err)
	}
	authData, err := authRepository.TestGetAuth(1)
	if err != nil {
		t.Fatalf("Error getting auth data: %v", err)
	}
	if authData.Password != newPassword {
		t.Fatalf("Client password not updated: %v", authData)
	}
}

func TestUpdateClientPasswordCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up Client tables: %v", err)
	}
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing client: %v", err)
	}
	err = clientService.UpdateClientPassword(result.UserID, tu.TestAuth, "test3")
	if err != nil {
		t.Fatalf("Error updating client password: %v", err)
	}
	result, err = authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: "test3",
	})
	if err != nil {
		t.Fatalf("Password not updated: %v", err)
	}
}

func TestUpdateClientPasswordIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	clientRepository := tu.CreateTestClientRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err := clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	err = clientService.UpdateClientPassword(2, tu.TestAuth, "test3")
	if err == nil {
		t.Fatalf("No error updating client password: %v", err)
	}
}

func TestUpdateClientPasswordIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up Client tables: %v", err)
	}
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	clientService := CreateClientService(
		clientRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
	)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	err = clientService.UpdateClientPassword(1, tu.TestAuth, "test3")
	if err == nil {
		t.Fatalf("No error updating client password: %v", err)
	}
}
