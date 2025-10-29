package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
)

func TestCreateContractCorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.ClientID != tu.TestServiceContractInitData.ClientID {
		t.Fatalf("client id is not correct: %v", contract.ClientID)
	}
	if contract.Description != tu.TestServiceContractInitData.Description {
		t.Fatalf("description is not correct: %v", contract.Description)
	}
	if contract.Price != tu.TestServiceContractInitData.Price {
		t.Fatalf("price is not correct: %v", contract.Price)
	}
	if contract.Commission != tu.TestServiceContractInitData.Commission {
		t.Fatalf("commission is not correct: %v", contract.Commission)
	}
	if contract.StartDate != tu.TestServiceContractInitData.StartDate {
		t.Fatalf("start date is not correct: %v", contract.StartDate)
	}

	if contract.Status != types.ContractStatusPending {
		t.Fatalf("status is not correct: %v", contract.Status)
	}
	if contract.PaymentStatus != types.PaymentStatusNull {
		t.Fatalf("payment status is not correct: %v", contract.PaymentStatus)
	}
	if contract.ReviewClientID != 0 {
		t.Fatalf("review client id is not correct: %v", contract.ReviewClientID)
	}
	if contract.ReviewRepetitorID != 0 {
		t.Fatalf("review repetitor id is not correct: %v", contract.ReviewRepetitorID)
	}
}

func TestCreateContractCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	contractService := CreateContractService(contractRepository, reviewRepository)
	tu.TestServiceContractInitData.ClientID = result.UserID
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.ClientID != tu.TestServiceContractInitData.ClientID {
		t.Fatalf("client id is not correct: %v", contract.ClientID)
	}
	if contract.Description != tu.TestServiceContractInitData.Description {
		t.Fatalf("description is not correct: %v", contract.Description)
	}
	if contract.Price != tu.TestServiceContractInitData.Price {
		t.Fatalf("price is not correct: %v", contract.Price)
	}
	if contract.Commission != tu.TestServiceContractInitData.Commission {
		t.Fatalf("commission is not correct: %v", contract.Commission)
	}
	if contract.Status != types.ContractStatusPending {
		t.Fatalf("status is not correct: %v", contract.Status)
	}
}

func TestGetContractCorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	contract, err := contractService.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.ClientID != tu.TestServiceContractInitData.ClientID {
		t.Fatalf("client id is not correct: %v", contract.ClientID)
	}
	if contract.Description != tu.TestServiceContractInitData.Description {
		t.Fatalf("description is not correct: %v", contract.Description)
	}
	if contract.Price != tu.TestServiceContractInitData.Price {
		t.Fatalf("price is not correct: %v", contract.Price)
	}
	if contract.Commission != tu.TestServiceContractInitData.Commission {
		t.Fatalf("commission is not correct: %v", contract.Commission)
	}
	if contract.StartDate != tu.TestServiceContractInitData.StartDate {
		t.Fatalf("start date is not correct: %v", contract.StartDate)
	}
	if contract.Status != types.ContractStatusPending {
		t.Fatalf("status is not correct: %v", contract.Status)
	}
	if contract.PaymentStatus != types.PaymentStatusNull {
		t.Fatalf("payment status is not correct: %v", contract.PaymentStatus)
	}
	if contract.ReviewClientID != 0 {
		t.Fatalf("review client id is not correct: %v", contract.ReviewClientID)
	}
	if contract.ReviewRepetitorID != 0 {
		t.Fatalf("review repetitor id is not correct: %v", contract.ReviewRepetitorID)
	}
}

func TestGetContractCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	contractService := CreateContractService(contractRepository, reviewRepository)
	tu.TestServiceContractInitData.ClientID = result.UserID
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	contract, err := contractService.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.ClientID != tu.TestServiceContractInitData.ClientID {
		t.Fatalf("client id is not correct: %v", contract.ClientID)
	}
	if contract.Description != tu.TestServiceContractInitData.Description {
		t.Fatalf("description is not correct: %v", contract.Description)
	}
	if contract.Price != tu.TestServiceContractInitData.Price {
		t.Fatalf("price is not correct: %v", contract.Price)
	}
	if contract.Commission != tu.TestServiceContractInitData.Commission {
		t.Fatalf("commission is not correct: %v", contract.Commission)
	}
	if contract.Status != types.ContractStatusPending {
		t.Fatalf("status is not correct: %v", contract.Status)
	}
}
func TestGetContractIncorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	_, err := contractService.GetContract(1)
	if err == nil {
		t.Fatalf("no error getting not existing contract: %v", err)
	}
}

func TestGetContractIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	contractService := CreateContractService(contractRepository, reviewRepository)
	_, err = contractService.GetContract(1)
	if err == nil {
		t.Fatalf("no error getting not existing contract: %v", err)
	}
}
func TestUpdateContractStatusCorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	err = contractService.UpdateContractStatus(contractID, types.ContractStatusActive)
	if err != nil {
		t.Fatalf("error updating contract status: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.Status != types.ContractStatusActive {
		t.Fatalf("status is not correct: %v", contract.Status)
	}
}

func TestUpdateContractStatusCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	contractService := CreateContractService(contractRepository, reviewRepository)
	tu.TestServiceContractInitData.ClientID = result.UserID
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	err = contractService.UpdateContractStatus(contractID, types.ContractStatusActive)
	if err != nil {
		t.Fatalf("error updating contract status: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.Status != types.ContractStatusActive {
		t.Fatalf("status is not correct: %v", contract.Status)
	}
}
func TestUpdateContractStatusIncorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	err := contractService.UpdateContractStatus(1, types.ContractStatusActive)
	if err == nil {
		t.Fatalf("no error updating not existing contract status: %v", err)
	}
}

func TestUpdateContractStatusIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	contractService := CreateContractService(contractRepository, reviewRepository)
	err = contractService.UpdateContractStatus(1, types.ContractStatusActive)
	if err == nil {
		t.Fatalf("no error updating not existing contract status: %v", err)
	}
}

func TestUpdateContractPaymentStatusCorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	err = contractService.UpdateContractPaymentStatus(contractID, types.PaymentStatusPaid)
	if err != nil {
		t.Fatalf("error updating contract payment status: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.PaymentStatus != types.PaymentStatusPaid {
		t.Fatalf("payment status is not correct: %v", contract.PaymentStatus)
	}
}

func TestUpdateContractPaymentStatusCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	contractService := CreateContractService(contractRepository, reviewRepository)
	tu.TestServiceContractInitData.ClientID = result.UserID
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	err = contractService.UpdateContractPaymentStatus(contractID, types.PaymentStatusPaid)
	if err != nil {
		t.Fatalf("error updating contract payment status: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.PaymentStatus != types.PaymentStatusPaid {
		t.Fatalf("payment status is not correct: %v", contract.PaymentStatus)
	}
}

func TestCreateContractReviewClientCorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	_, err = contractService.CreateContractReviewClient(contractID, tu.TestServiceReview)
	if err != nil {
		t.Fatalf("error creating contract review client: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.ReviewClientID == 0 {
		t.Fatalf("review client id is not correct: %v", contract.ReviewClientID)
	}
}

func TestCreateContractReviewClientCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	contractService := CreateContractService(contractRepository, reviewRepository)
	tu.TestServiceContractInitData.ClientID = result.UserID
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	tu.TestServiceReview.ClientID = result.UserID
	_, err = contractService.CreateContractReviewClient(contractID, tu.TestServiceReview)
	if err != nil {
		t.Fatalf("error creating contract review client: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.ReviewClientID == 0 {
		t.Fatalf("review client id is not correct: %v", contract.ReviewClientID)
	}
}

func TestCreateContractReviewClientIncorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	_, err := contractService.CreateContractReviewClient(1, tu.TestServiceReview)
	if err == nil {
		t.Fatalf("no error creating contract review client: %v", err)
	}
}

func TestCreateContractReviewClientIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	contractService := CreateContractService(contractRepository, reviewRepository)
	_, err = contractService.CreateContractReviewClient(1, tu.TestServiceReview)
	if err == nil {
		t.Fatalf("no error creating contract review client: %v", err)
	}
}

func TestCreateContractReviewRepetitorCorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	_, err = contractService.CreateContractReviewRepetitor(contractID, tu.TestServiceReview)
	if err != nil {
		t.Fatalf("error creating contract review repetitor: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.ReviewClientID != 0 {
		t.Fatalf("review client id is not correct: %v", contract.ReviewClientID)
	}
}

func TestCreateContractReviewRepetitorCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	authRepository := module.AuthRepository
	repetitorRepository := module.RepetitorRepository
	resumeRepository := module.ResumeRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	repetitorService := CreateRepetitorService(repetitorRepository, personalDataRepository, userRepository, reviewRepository, resumeRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	contractService := CreateContractService(contractRepository, reviewRepository)
	tu.TestServiceContractInitData.ClientID = result.UserID
	contractID, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	tu.TestInitRepetitorData.Login = "repetitor"
	tu.TestInitRepetitorData.Password = "repetitor2"
	err = repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("error creating repetitor: %v", err)
	}
	result, err = authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestInitRepetitorData.Login,
		Password: tu.TestInitRepetitorData.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	tu.TestServiceReview.RepetitorID = result.UserID
	_, err = contractService.CreateContractReviewRepetitor(contractID, tu.TestServiceReview)
	if err != nil {
		t.Fatalf("error creating contract review repetitor: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("error getting contract: %v", err)
	}
	if contract.ReviewRepetitorID == 0 {
		t.Fatalf("review client id is not correct: %v", contract.ReviewClientID)
	}
}

func TestCreateContractReviewRepetitorIncorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	_, err := contractService.CreateContractReviewRepetitor(1, tu.TestServiceReview)
	if err == nil {
		t.Fatalf("no error creating contract review repetitor: %v", err)
	}
}

func TestCreateContractReviewRepetitorIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	contractService := CreateContractService(contractRepository, reviewRepository)
	_, err = contractService.CreateContractReviewRepetitor(1, tu.TestServiceReview)
	if err == nil {
		t.Fatalf("no error creating contract review repetitor: %v", err)
	}
}

func TestUpdateContractReviewClientIncorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	err := contractService.UpdateContractReviewClient(1, tu.TestServiceReview)
	if err == nil {
		t.Fatalf("no error updating contract review client: %v", err)
	}
}

func TestUpdateContractReviewClientIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	contractService := CreateContractService(contractRepository, reviewRepository)
	err = contractService.UpdateContractReviewClient(1, tu.TestServiceReview)
	if err == nil {
		t.Fatalf("no error updating contract review client: %v", err)
	}
}

func TestUpdateContractReviewRepetitorIncorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	err := contractService.UpdateContractReviewRepetitor(1, tu.TestServiceReview)
	if err == nil {
		t.Fatalf("no error updating contract review repetitor: %v", err)
	}
}

func TestUpdateContractReviewRepetitorIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	contractService := CreateContractService(contractRepository, reviewRepository)
	err = contractService.UpdateContractReviewRepetitor(1, tu.TestServiceReview)
	if err == nil {
		t.Fatalf("no error updating contract review repetitor: %v", err)
	}
}

func TestGetRepetitorContractListCorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	_, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	_, err = contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	contractList, err := contractService.GetRepetitorContractList(tu.TestServiceReview.RepetitorID, 0, 10, types.ContractStatusPending)
	if err != nil {
		t.Fatalf("error getting contract list: %v", err)
	}
	if len(contractList) != 0 {
		t.Fatalf("contract list is not correct: %v", contractList)
	}
}

func TestGetRepetitorContractListCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractList, err := contractService.GetRepetitorContractList(tu.TestServiceReview.RepetitorID, 0, 10, types.ContractStatusPending)
	if err != nil {
		t.Fatalf("error getting contract list: %v", err)
	}
	if len(contractList) != 0 {
		t.Fatalf("contract list is not correct: %v", contractList)
	}
}

func TestGetClientContractListCorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	_, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	_, err = contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	contractList, err := contractService.GetClientContractList(tu.TestServiceReview.ClientID, 0, 10, types.ContractStatusPending)
	if err != nil {
		t.Fatalf("error getting contract list: %v", err)
	}
	if len(contractList) != 2 {
		t.Fatalf("contract list is not correct: %v", contractList)
	}
}

func TestGetClientContractListCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	contractService := CreateContractService(contractRepository, reviewRepository)
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	tu.TestServiceContractInitData.ClientID = result.UserID
	_, err = contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	_, err = contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	contractList, err := contractService.GetClientContractList(result.UserID, 0, 10, types.ContractStatusPending)
	if err != nil {
		t.Fatalf("error getting contract list: %v", err)
	}
	if len(contractList) != 2 {
		t.Fatalf("contract list is not correct: %v", contractList)
	}
}
func TestGetContractListCorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	_, err := contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	_, err = contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	contractList, err := contractService.GetContractList(0, 10, types.ContractStatusPending)
	if err != nil {
		t.Fatalf("error getting contract list: %v", err)
	}
	if len(contractList) != 2 {
		t.Fatalf("contract list is not correct: %v", contractList)
	}
}

func TestGetContractListCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	tu.TestServiceContractInitData.ClientID = result.UserID
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractList, err := contractService.GetContractList(0, 10, types.ContractStatusPending)
	if err != nil {
		t.Fatalf("error getting contract list: %v", err)
	}
	if len(contractList) != 0 {
		t.Fatalf("contract list is not correct: %v", contractList)
	}
	_, err = contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	_, err = contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	contractList, err = contractService.GetContractList(0, 10, types.ContractStatusPending)
	if err != nil {
		t.Fatalf("error getting contract list: %v", err)
	}
	if len(contractList) != 2 {
		t.Fatalf("contract list is not correct: %v", contractList)
	}
}

func TestGetAllContractsCorrectLondon(t *testing.T) {
	contractRepository := tu.CreateTestContractRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	contractService := CreateContractService(contractRepository, reviewRepository)
	contractList, err := contractService.GetAllContracts(0, 10)
	if err != nil {
		t.Fatalf("error getting contract list: %v", err)
	}
	if len(contractList) != 0 {
		t.Fatalf("contract list is not correct: %v", contractList)
	}
	contractList, err = contractService.GetAllContracts(0, 10)
	if err != nil {
		t.Fatalf("error getting contract list: %v", err)
	}
	if len(contractList) != 0 {
		t.Fatalf("contract list is not correct: %v", contractList)
	}
}

func TestGetAllContractsCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := module.ContractRepository
	reviewRepository := module.ReviewRepository
	contractService := CreateContractService(contractRepository, reviewRepository)
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	tu.TestServiceContractInitData.ClientID = result.UserID
	contractList, err := contractService.GetAllContracts(0, 10)
	if err != nil {
		t.Fatalf("error getting contract list: %v", err)
	}
	if len(contractList) != 0 {
		t.Fatalf("contract list is not correct: %v", contractList)
	}
	_, err = contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	_, err = contractService.CreateContract(tu.TestServiceContractInitData)
	if err != nil {
		t.Fatalf("error creating contract: %v", err)
	}
	contractList, err = contractService.GetAllContracts(0, 10)
	if err != nil {
		t.Fatalf("error getting contract list: %v", err)
	}
	if len(contractList) != 2 {
		t.Fatalf("contract list is not correct: %v", contractList)
	}
}
