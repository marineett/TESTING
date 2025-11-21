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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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

func CheckContract(
	t *testing.T,
	contract *types.DBContract,
	contractID int64,
	clientID int64,
	repetitorID int64,
	description string,
	status types.ContractStatus,
	paymentStatus types.PaymentStatus,
	reviewClientID int64,
	reviewRepetitorID int64,
	price int64,
	commission int64,
	transactionID int64,
) {
	if contract.ID != contractID {
		t.Fatalf("Contract id is not correct: %v", contract.ID)
	}
	if contract.ClientID != clientID {
		t.Fatalf("Contract client id is not correct: %v", contract.ClientID)
	}
	if contract.RepetitorID != repetitorID {
		t.Fatalf("Contract repetitor id is not correct: %v", contract.RepetitorID)
	}
	if contract.Description != description {
		t.Fatalf("Contract description is not correct: %v", contract.Description)
	}
	if contract.Status != status {
		t.Fatalf("Contract status is not correct: %v", contract.Status)
	}
}

func TestGetContractCorrect(t *testing.T) {
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
	CheckContract(
		t,
		contract,
		contractID,
		clientID,
		0,
		tu.TestContract.Description,
		tu.TestContract.Status,
		tu.TestContract.PaymentStatus,
		tu.TestContract.ReviewClientID,
		tu.TestContract.ReviewRepetitorID,
		tu.TestContract.Price,
		tu.TestContract.Commission,
		tu.TestContract.TransactionID,
	)
}

func TestGetContractIncorrect(t *testing.T) {
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	_, err = contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
}

func TestUpdateContractRepetitorIDIncorrect(t *testing.T) {
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
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
	err = contractRepository.UpdateContractRepetitorID(contractID, repetitorID)
	if err != nil {
		t.Fatalf("Error updating contract repetitor id: %v", err)
	}
	_, err = contractRepository.InsertContract(tu.TestContract)
	if err != nil {
		t.Fatalf("Error inserting contract: %v", err)
	}
	_, err = contractRepository.GetContractsByRepetitorID(repetitorID, 0, 10, types.ContractStatusActive)
	if err != nil {
		t.Fatalf("Error getting contracts by repetitor id: %v", err)
	}
}

func CheckLengths(t *testing.T, contracts []types.DBContract, length int) {
	if len(contracts) != length {
		t.Fatalf("Number of contracts is not correct: %v", len(contracts))
	}
}

func TestGetContractListCorrect(t *testing.T) {
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
	CheckLengths(t, contracts, 1)
	contracts, err = contractRepository.GetContractList(0, 10, types.ContractStatusActive)
	if err != nil {
		t.Fatalf("Error getting contract list: %v", err)
	}
	CheckLengths(t, contracts, 0)
	contracts, err = contractRepository.GetContractList(0, 10, types.ContractStatusCompleted)
	if err != nil {
		t.Fatalf("Error getting contract list: %v", err)
	}
	CheckLengths(t, contracts, 0)
}

func CheckContractList(t *testing.T, contracts []types.DBContract, length int) {
	if len(contracts) != length {
		t.Fatalf("Number of contracts is not correct: %v", len(contracts))
	}
}

func SetupDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	return db
}

func TestGetAllContractsCorrect(t *testing.T) {
	db := SetupDatabase(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err := setupContractTables(db)
	if err != nil {
		t.Fatalf("Error setting up contract tables: %v", err)
	}
	contractRepository := CreateSqlContractRepository(db, "contracts", "sequence")
	contracts, err := contractRepository.GetAllContracts(0, 10)
	if err != nil {
		t.Fatalf("Error getting all contracts: %v", err)
	}
	CheckContractList(t, contracts, 0)
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
	CheckContractList(t, contracts, 3)
}

func TestUpdateContractStatusCorrect(t *testing.T) {
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	_, err = contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
}

func TestUpdateContractReviewClientIDIncorrect(t *testing.T) {
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	_, err = contractRepository.GetContract(contractID)
	if err != nil {
		t.Fatalf("Error getting contract: %v", err)
	}
}

func TestUpdateContractReviewRepetitorIDIncorrect(t *testing.T) {
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
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

func TestUpdateContractReviewClientIDInSeqIncorrect(t *testing.T) {
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
	defer func() {
		_ = tx.Rollback()
	}()
	err = contractRepository.UpdateContractReviewClientIDInSeq(tx, 1, 1)
	if err == nil {
		t.Fatalf("No error updating contract review client id in seq: %v", err)
	}
}
