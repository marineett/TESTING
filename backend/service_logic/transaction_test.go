package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
)

func TestCreateContractPaymentTransactionCorrectLondon(t *testing.T) {
	transactionRepository := tu.CreateTestTransactionRepository()
	transactionService := CreateTransactionService(transactionRepository)
	transactionID, err := transactionService.CreateContractPaymentTransaction(100, 1, 5)
	if err != nil {
		t.Fatalf("error creating contract payment transaction: %v", err)
	}
	transaction, err := transactionRepository.GetTransaction(transactionID)
	if err != nil {
		t.Fatalf("error getting transaction: %v", err)
	}
	if transaction.Status != types.TransactionStatusPending {
		t.Fatalf("transaction status is not correct: %v", transaction.Status)
	}
}

func TestCreateContractPaymentTransactionCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := module.TransactionRepository
	transactionService := CreateTransactionService(transactionRepository)
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	authRepository := module.AuthRepository
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	transactionID, err := transactionService.CreateContractPaymentTransaction(100, result.UserID, 5)
	if err != nil {
		t.Fatalf("error creating contract payment transaction: %v", err)
	}
	transaction, err := transactionRepository.GetTransaction(transactionID)
	if err != nil {
		t.Fatalf("error getting transaction: %v", err)
	}
	if transaction.Status != types.TransactionStatusPending {
		t.Fatalf("transaction status is not correct: %v", transaction.Status)
	}
}
func TestGetTransactionsListCorrectLondon(t *testing.T) {
	transactionRepository := tu.CreateTestTransactionRepository()
	transactionService := CreateTransactionService(transactionRepository)
	transactions, err := transactionService.GetTransactionsList(1, 0, 10)
	if err != nil {
		t.Fatalf("error getting transactions list: %v", err)
	}
	if len(transactions) != 0 {
		t.Fatalf("transactions list is not correct: %v", transactions)
	}
	transactionRepository.InsertTransaction(tu.TestTransaction)
	transactions, err = transactionService.GetTransactionsList(1, 0, 10)
	if err != nil {
		t.Fatalf("error getting transactions list: %v", err)
	}
	if len(transactions) != 1 {
		t.Fatalf("transactions list is not correct: %v", transactions)
	}
}

func TestGetTransactionsListCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := module.TransactionRepository
	transactionService := CreateTransactionService(transactionRepository)
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	authRepository := module.AuthRepository
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	_, err = transactionService.CreateContractPaymentTransaction(100, result.UserID, 5)
	if err != nil {
		t.Fatalf("error creating contract payment transaction: %v", err)
	}
	transactions, err := transactionService.GetTransactionsList(result.UserID, 0, 10)
	if err != nil {
		t.Fatalf("error getting transactions list: %v", err)
	}
	if len(transactions) != 1 {
		t.Fatalf("transactions list is not correct: %v", transactions)
	}
}

func TestGetTransactionsListIncorrectLondon(t *testing.T) {
	transactionRepository := tu.CreateTestTransactionRepository()
	transactionService := CreateTransactionService(transactionRepository)
	transactions, err := transactionService.GetTransactionsList(1, 0, 10)
	if err != nil {
		t.Fatalf("error getting transactions list: %v", err)
	}
	if len(transactions) != 0 {
		t.Fatalf("transactions list is not correct: %v", transactions)
	}
}

func TestGetTransactionsListIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := module.TransactionRepository
	transactionService := CreateTransactionService(transactionRepository)
	transactions, err := transactionService.GetTransactionsList(1, 0, 10)
	if err != nil {
		t.Fatalf("error getting transactions list: %v", err)
	}
	if len(transactions) != 0 {
		t.Fatalf("transactions list is not correct: %v", transactions)
	}
}

func TestChangeTransactionStatusCorrectLondon(t *testing.T) {
	transactionRepository := tu.CreateTestTransactionRepository()
	transactionService := CreateTransactionService(transactionRepository)
	transactionID, err := transactionService.CreateContractPaymentTransaction(100, 1, 5)
	if err != nil {
		t.Fatalf("error creating contract payment transaction: %v", err)
	}
	transaction, err := transactionRepository.GetTransaction(transactionID)
	if err != nil {
		t.Fatalf("error getting transaction: %v", err)
	}
	if transaction.Status != types.TransactionStatusPending {
		t.Fatalf("transaction status is not correct: %v", transaction.Status)
	}
	err = transactionService.ChangeTransactionStatus(transactionID, types.TransactionStatusPaid)
	if err != nil {
		t.Fatalf("error changing transaction status: %v", err)
	}
	transaction, err = transactionRepository.GetTransaction(transactionID)
	if err != nil {
		t.Fatalf("error getting transaction: %v", err)
	}
	if transaction.Status != types.TransactionStatusPaid {
		t.Fatalf("transaction status is not correct: %v", transaction.Status)
	}
}

func TestChangeTransactionStatusCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := module.TransactionRepository
	transactionService := CreateTransactionService(transactionRepository)
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	authRepository := module.AuthRepository
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	transactionID, err := transactionService.CreateContractPaymentTransaction(100, result.UserID, 5)
	if err != nil {
		t.Fatalf("error creating contract payment transaction: %v", err)
	}
	err = transactionService.ChangeTransactionStatus(transactionID, types.TransactionStatusPaid)
	if err != nil {
		t.Fatalf("error changing transaction status: %v", err)
	}
	transaction, err := transactionRepository.GetTransaction(transactionID)
	if err != nil {
		t.Fatalf("error getting transaction: %v", err)
	}
	if transaction.Status != types.TransactionStatusPaid {
		t.Fatalf("transaction status is not correct: %v", transaction.Status)
	}
}

func TestChangeTransactionStatusIncorrectLondon(t *testing.T) {
	transactionRepository := tu.CreateTestTransactionRepository()
	transactionService := CreateTransactionService(transactionRepository)
	transactionID, err := transactionService.CreateContractPaymentTransaction(100, 1, 5)
	if err != nil {
		t.Fatalf("error creating contract payment transaction: %v", err)
	}
	err = transactionService.ChangeTransactionStatus(transactionID, types.TransactionStatusPaid)
	if err != nil {
		t.Fatalf("error changing transaction status: %v", err)
	}
	transaction, err := transactionRepository.GetTransaction(transactionID)
	if err != nil {
		t.Fatalf("error getting transaction: %v", err)
	}
	if transaction.Status != types.TransactionStatusPaid {
		t.Fatalf("transaction status is not correct: %v", transaction.Status)
	}
	err = transactionService.ChangeTransactionStatus(transactionID, types.TransactionStatusPaid)
	if err != nil {
		t.Fatalf("error changing transaction status: %v", err)
	}
	transaction, err = transactionRepository.GetTransaction(transactionID)
	if err != nil {
		t.Fatalf("error getting transaction: %v", err)
	}
	if transaction.Status != types.TransactionStatusPaid {
		t.Fatalf("transaction status is not correct: %v", transaction.Status)
	}
}

func TestChangeTransactionStatusIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := module.TransactionRepository
	transactionService := CreateTransactionService(transactionRepository)
	err = transactionService.ChangeTransactionStatus(1, types.TransactionStatusPaid)
	if err == nil {
		t.Fatalf("No error changing transaction status: %v", err)
	}
}
func TestGetTransactionCorrectLondon(t *testing.T) {
	transactionRepository := tu.CreateTestTransactionRepository()
	transactionService := CreateTransactionService(transactionRepository)
	transactionID, err := transactionService.CreateContractPaymentTransaction(100, 1, 5)
	if err != nil {
		t.Fatalf("error creating contract payment transaction: %v", err)
	}
	transaction, err := transactionService.GetTransaction(transactionID)
	if err != nil {
		t.Fatalf("error getting transaction: %v", err)
	}
	if transaction.ID != transactionID {
		t.Fatalf("transaction id is not correct: %v", transaction.ID)
	}
	if transaction.Status != types.TransactionStatusPending {
		t.Fatalf("transaction status is not correct: %v", transaction.Status)
	}
}

func TestGetTransactionCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := module.TransactionRepository
	transactionService := CreateTransactionService(transactionRepository)
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	authRepository := module.AuthRepository
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	transactionID, err := transactionService.CreateContractPaymentTransaction(100, result.UserID, 5)
	if err != nil {
		t.Fatalf("error creating contract payment transaction: %v", err)
	}
	transaction, err := transactionService.GetTransaction(transactionID)
	if err != nil {
		t.Fatalf("error getting transaction: %v", err)
	}
	if transaction.ID != transactionID {
		t.Fatalf("transaction id is not correct: %v", transaction.ID)
	}
	if transaction.Status != types.TransactionStatusPending {
		t.Fatalf("transaction status is not correct: %v", transaction.Status)
	}
}

func TestGetTransactionIncorrectLondon(t *testing.T) {
	transactionRepository := tu.CreateTestTransactionRepository()
	transactionService := CreateTransactionService(transactionRepository)
	_, err := transactionService.GetTransaction(1)
	if err == nil {
		t.Fatalf("no error with getting transaction with wrong id: %v", err)
	}
}

func TestGetTransactionIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := module.TransactionRepository
	transactionService := CreateTransactionService(transactionRepository)
	_, err = transactionService.GetTransaction(1)
	if err == nil {
		t.Fatalf("no error with getting transaction with wrong id: %v", err)
	}
}
