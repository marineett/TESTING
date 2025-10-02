package data_base

import (
	tu "data_base_project/test_database_utility"
	"data_base_project/types"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func setupContractTables(db *sql.DB) error {
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
	err = CreateSqlClientTable(db, "clients", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating client table: %v", err)
	}
	err = CreateSqlRepetitorTable(db, "repetitors", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating repetitor table: %v", err)
	}
	err = CreateSqlContractTable(db, "contracts", "users", "reviews", "repetitors", "clients")
	if err != nil {
		return fmt.Errorf("error creating contract table: %v", err)
	}
	err = CreateSqlReviewTable(db, "reviews", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating review table: %v", err)
	}
	return nil
}

func TestCreateSqlContractRepositoryCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
}

func TestInsertContractCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	_, err = contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
}

func TestInsertContractIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	_, err = contractRepository.InsertContract(tu.TestContract)
	if err == nil {
		t.Fatalf("No error inserting contract: %v", err)
	}
}

func TestGetContractCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
	if contract.ID != contractID {
		t.Fatalf("Contract id is not correct: %v", contract.ID)
	}
	if contract.ClientID != clientID {
		t.Fatalf("Contract client id is not correct: %v", contract.ClientID)
	}
	if contract.RepetitorID != 0 {
		t.Fatalf("Contract repetitor id is not correct: %v", contract.RepetitorID)
	}
	if contract.Description != tu.TestContract.Description {
		t.Fatalf("Contract description is not correct: %v", contract.Description)
	}
	if contract.Status != tu.TestContract.Status {
		t.Fatalf("Contract status is not correct: %v", contract.Status)
	}
	if contract.PaymentStatus != tu.TestContract.PaymentStatus {
		t.Fatalf("Contract payment status is not correct: %v", contract.PaymentStatus)
	}
	if contract.ReviewClientID != tu.TestContract.ReviewClientID {
		t.Fatalf("Contract review client id is not correct: %v", contract.ReviewClientID)
	}
	if contract.ReviewRepetitorID != tu.TestContract.ReviewRepetitorID {
		t.Fatalf("Contract review repetitor id is not correct: %v", contract.ReviewRepetitorID)
	}
	if contract.Price != tu.TestContract.Price {
		t.Fatalf("Contract price is not correct: %v", contract.Price)
	}
	if contract.Commission != tu.TestContract.Commission {
		t.Fatalf("Contract commission is not correct: %v", contract.Commission)
	}
	if contract.TransactionID != tu.TestContract.TransactionID {
		t.Fatalf("Contract transaction id is not correct: %v", contract.TransactionID)
	}
}

func TestGetContractIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	_, err = contractRepository.GetContract(1)
	if err == nil {
		t.Fatalf("No error getting contract: %v", err)
	}
}

func TestGetContractsByClientIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	_, err = contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	_, err = contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	contracts, err := contractRepository.GetContractsByClientID(clientID, 0, 10, types.ContractStatusPending)
	if err != nil {
		t.Fatalf("Error getting contracts by client id: %v", err)
	}
	if len(contracts) != 2 {
		t.Fatalf("Number of contracts is not correct: %v", len(contracts))
	}
}

func TestUpdateContractRepetitorIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	err = contractRepository.UpdateContractRepetitorID(contractID, repetitorID)
	if err != nil {
		t.Fatalf("Error updating contract repetitor id: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
	if contract.RepetitorID != repetitorID {
		t.Fatalf("Contract repetitor id is not correct: %v", contract.RepetitorID)
	}
}

func TestUpdateContractRepetitorIDIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	err = contractRepository.UpdateContractRepetitorID(1, 2)
	if err == nil {
		t.Fatalf("No error updating contract repetitor id: %v", err)
	}
}

func TestGetContractsByRepetitorIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
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
	contractRepository.UpdateContractRepetitorID(contractID, repetitorID)
	_, err = contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	contracts, err := contractRepository.GetContractsByRepetitorID(repetitorID, 0, 10, types.ContractStatusActive)
	if err != nil {
		t.Fatalf("Error getting contracts by repetitor id: %v", err)
	}
	if len(contracts) != 1 {
		t.Fatalf("Number of contracts is not correct: %v", len(contracts))
	}
}

func TestGetContractListCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	_, err = contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	contracts, err := contractRepository.GetContractList(0, 10, types.ContractStatusPending)
	if err != nil {
		t.Fatalf("Error getting contract list: %v", err)
	}
	if len(contracts) != 1 {
		t.Fatalf("Number of contracts is not correct: %v", len(contracts))
	}
	contracts, err = contractRepository.GetContractList(0, 10, types.ContractStatusActive)
	if err != nil {
		t.Fatalf("Error getting contract list: %v", err)
	}
	if len(contracts) != 0 {
		t.Fatalf("Number of contracts is not correct: %v", len(contracts))
	}
	contracts, err = contractRepository.GetContractList(0, 10, types.ContractStatusCompleted)
	if err != nil {
		t.Fatalf("Error getting contract list: %v", err)
	}
	if len(contracts) != 0 {
		t.Fatalf("Number of contracts is not correct: %v", len(contracts))
	}
}

func TestGetAllContractsCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	contracts, err := contractRepository.GetAllContracts(0, 10)
	if err != nil {
		t.Fatalf("Error getting all contracts: %v", err)
	}
	if len(contracts) != 0 {
		t.Fatalf("Number of contracts is not correct: %v", len(contracts))
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	_, err = contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	_, err = contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	_, err = contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	contracts, err = contractRepository.GetAllContracts(0, 10)
	if err != nil {
		t.Fatalf("Error getting all contracts: %v", err)
	}
	if len(contracts) != 3 {
		t.Fatalf("Number of contracts is not correct: %v", len(contracts))
	}
}

func TestUpdateContractStatusCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	err = contractRepository.UpdateContractStatus(contractID, types.ContractStatusActive)
	if err != nil {
		t.Fatalf("Error updating contract status: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
	if contract.Status != types.ContractStatusActive {
		t.Fatalf("Contract status is not correct: %v", contract.Status)
	}
}

func TestUpdateContractStatusIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	err = contractRepository.UpdateContractStatus(1, types.ContractStatusActive)
	if err == nil {
		t.Fatalf("No error updating contract status: %v", err)
	}
}

func TestUpdateContractPaymentStatusCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	err = contractRepository.UpdateContractPaymentStatus(contractID, types.PaymentStatusPaid)
	if err != nil {
		t.Fatalf("Error updating contract payment status: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
	if contract.PaymentStatus != types.PaymentStatusPaid {
		t.Fatalf("Contract payment status is not correct: %v", contract.PaymentStatus)
	}
}

func TestUpdateContractPaymentStatusIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	err = contractRepository.UpdateContractPaymentStatus(1, types.PaymentStatusPaid)
	if err == nil {
		t.Fatalf("No error updating contract payment status: %v", err)
	}
}

func TestUpdateContractReviewClientIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	reviewRepository := CreateSqlReviewRepository(db, "reviews", "sequence")
	reviewID, err := reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	err = contractRepository.UpdateContractReviewClientID(contractID, reviewID)
	if err != nil {
		t.Fatalf("Error updating contract review client id: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
	if contract.ReviewClientID != reviewID {
		t.Fatalf("Contract review client id is not correct: %v", contract.ReviewClientID)
	}
}

func TestUpdateContractReviewClientIDIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	err = contractRepository.UpdateContractReviewClientID(1, 1)
	if err == nil {
		t.Fatalf("No error updating contract review client id: %v", err)
	}
	err = contractRepository.UpdateContractReviewClientID(1, 1)
	if err == nil {
		t.Fatalf("No error updating contract review client id: %v", err)
	}
}

func TestUpdateContractReviewRepetitorIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	reviewRepository := CreateSqlReviewRepository(db, "reviews", "sequence")
	reviewID, err := reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	err = contractRepository.UpdateContractReviewRepetitorID(contractID, reviewID)
	if err != nil {
		t.Fatalf("Error updating contract review repetitor id: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
	if contract.ReviewRepetitorID != reviewID {
		t.Fatalf("Contract review repetitor id is not correct: %v", contract.ReviewRepetitorID)
	}
}

func TestUpdateContractReviewRepetitorIDIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	err = contractRepository.UpdateContractReviewRepetitorID(1, 1)
	if err == nil {
		t.Fatalf("No error updating contract review repetitor id: %v", err)
	}
	err = contractRepository.UpdateContractReviewRepetitorID(1, 1)
	if err == nil {
		t.Fatalf("No error updating contract review repetitor id: %v", err)
	}
}

func TestUpdateContractPriceCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	err = contractRepository.UpdateContractPrice(contractID, 1000)
	if err != nil {
		t.Fatalf("Error updating contract price: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
	if contract.Price != 1000 {
		t.Fatalf("Contract price is not correct: %v", contract.Price)
	}
}

func TestUpdateContractPriceIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	err = contractRepository.UpdateContractPrice(1, 1000)
	if err == nil {
		t.Fatalf("No error updating contract price: %v", err)
	}
	err = contractRepository.UpdateContractPrice(1, 1000)
	if err == nil {
		t.Fatalf("No error updating contract price: %v", err)
	}
}

func TestUpdateContractCommissionCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	err = contractRepository.UpdateContractCommission(contractID, 20)
	if err != nil {
		t.Fatalf("Error updating contract price: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
	if contract.Commission != 20 {
		t.Fatalf("Contract price is not correct: %v", contract.Price)
	}
}

func TestUpdateContractCommissionIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	err = contractRepository.UpdateContractCommission(1, 20)
	if err == nil {
		t.Fatalf("No error updating contract commission: %v", err)
	}
	err = contractRepository.UpdateContractCommission(1, 20)
	if err == nil {
		t.Fatalf("No error updating contract commission: %v", err)
	}
	err = contractRepository.UpdateContractCommission(1, 20)
	if err == nil {
		t.Fatalf("No error updating contract commission: %v", err)
	}
}

func TestUpdateContractStartDateCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	err = contractRepository.UpdateContractStartDate(contractID, time.Now())
	if err != nil {
		t.Fatalf("Error updating contract start date: %v", err)
	}
	_, err = contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
}

func TestUpdateContractStartDateIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	err = contractRepository.UpdateContractStartDate(1, time.Now())
	if err == nil {
		t.Fatalf("No error updating contract start date: %v", err)
	}
}

func TestUpdateContractReviewClientIDInSeqCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestContract.ClientID = clientID
	contractID, err := contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	reviewRepository := CreateSqlReviewRepository(db, "reviews", "sequence")
	reviewID, err := reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	defer tx.Rollback()
	err = contractRepository.UpdateContractReviewClientIDInSeq(tx, contractID, reviewID)
	if err != nil {
		t.Fatalf("Error updating contract review client id in seq: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		t.Fatalf("Error committing transaction: %v", err)
	}
	contract, err := contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
	if contract.ReviewClientID != reviewID {
		t.Fatalf("Contract review client id is not correct: %v", contract.ReviewClientID)
	}
}

func TestUpdateContractReviewClientIDInSeqIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	if contractRepository == nil {
		t.Fatalf("Error creating contract repository: %v", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	defer tx.Rollback()
	err = contractRepository.UpdateContractReviewClientIDInSeq(tx, 1, 1)
	if err == nil {
		t.Fatalf("No error updating contract review client id in seq: %v", err)
	}
}
