package data_base

import (
	tu "data_base_project/test_database_utility"
	"data_base_project/types"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
)

func setupTransactionTables(db *sql.DB) error {
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
	err = CreateSqlTransactionTable(db, "transactions", "users", "pending_contract_payment_transactions", "sequence")
	if err != nil {
		return fmt.Errorf("error creating transaction table: %v", err)
	}
	err = CreateSqlPendingContractPaymentTransactionsTable(db, "pending_contract_payment_transactions", "users", "transactions", "sequence")
	if err != nil {
		return fmt.Errorf("error creating pending contract payment transactions table: %v", err)
	}
	return nil
}

func TestCreateSqlTransactionTable(t *testing.T) {
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
	err = setupTransactionTables(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := CreateSqlTransactionRepository(db, "transactions", "pending_contract_payment_transactions", "sequence")
	if transactionRepository == nil {
		t.Fatalf("Error creating transaction repository: %v", err)
	}
}

func TestInsertTransactionCorrect(t *testing.T) {
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
	err = setupTransactionTables(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := CreateSqlTransactionRepository(db, "transactions", "pending_contract_payment_transactions", "sequence")
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataID, err := personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	tu.TestUser.PersonalDataID = personalDataID
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	userID, err := userRepository.InsertUser(tu.TestUser)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	tu.TestTransaction.UserID = userID
	_, err = transactionRepository.InsertTransaction(tu.TestTransaction)
	if err != nil {
		t.Fatalf("Error inserting transaction: %v", err)
	}
}

func TestInsertTransactionIncorrect(t *testing.T) {
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
	err = setupTransactionTables(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := CreateSqlTransactionRepository(db, "transactions", "pending_contract_payment_transactions", "sequence")
	_, err = transactionRepository.InsertTransaction(tu.TestTransaction)
	if err != nil {
		t.Fatalf("No error inserting transaction: %v", err)
	}
}

func CheckTransaction(
	t *testing.T,
	transaction *types.DBTransaction,
	transactionID int64,
	userID int64,
	amount int64,
	status types.TransactionStatus,
	transactionType types.TransactionType,
) {
	if transaction.ID != transactionID {
		t.Fatalf("Transaction id not correct: %v", transaction.ID)
	}
	if transaction.UserID != userID {
		t.Fatalf("Transaction user id not correct: %v", transaction.UserID)
	}
	if transaction.Amount != amount {
		t.Fatalf("Transaction amount not correct: %v", transaction.Amount)
	}
	if transaction.Status != status {
		t.Fatalf("Transaction status not correct: %v", transaction.Status)
	}
	if transaction.Type != transactionType {
		t.Fatalf("Transaction type not correct: %v", transaction.Type)
	}
}
func TestGetTransactionCorrect(t *testing.T) {
	db := SetupDatabase(t)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err := setupTransactionTables(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := CreateSqlTransactionRepository(db, "transactions", "pending_contract_payment_transactions", "sequence")
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataID, err := personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	tu.TestUser.PersonalDataID = personalDataID
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	userID, err := userRepository.InsertUser(tu.TestUser)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	tu.TestTransaction.UserID = userID
	transactionID, err := transactionRepository.InsertTransaction(tu.TestTransaction)
	if err != nil {
		t.Fatalf("Error inserting transaction: %v", err)
	}
	transaction, err := transactionRepository.GetTransaction(transactionID)
	if err != nil {
		t.Fatalf("Error getting transaction: %v", err)
	}
	CheckTransaction(
		t,
		transaction,
		transactionID,
		userID,
		tu.TestTransaction.Amount,
		tu.TestTransaction.Status,
		tu.TestTransaction.Type,
	)
}

func TestGetTransactionIncorrect(t *testing.T) {
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
	err = setupTransactionTables(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := CreateSqlTransactionRepository(db, "transactions", "pending_contract_payment_transactions", "sequence")
	_, err = transactionRepository.GetTransaction(1)
	if err == nil {
		t.Fatalf("No error getting transaction: %v", err)
	}
}

func CheckTransactionsLength(t *testing.T, transactions []types.DBTransaction, length int) {
	if len(transactions) != length {
		t.Fatalf("Transactions not updated: %v", transactions)
	}
}

func TestGetTransactionsListCorrect(t *testing.T) {
	db := SetupDatabase(t)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err := setupTransactionTables(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := CreateSqlTransactionRepository(db, "transactions", "pending_contract_payment_transactions", "sequence")
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataID, err := personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	tu.TestUser.PersonalDataID = personalDataID
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	userID, err := userRepository.InsertUser(tu.TestUser)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	tu.TestTransaction.UserID = userID
	_, err = transactionRepository.InsertTransaction(tu.TestTransaction)
	if err != nil {
		t.Fatalf("Error inserting transaction: %v", err)
	}
	tu.TestTransaction.UserID = userID
	_, err = transactionRepository.InsertTransaction(tu.TestTransaction)
	if err != nil {
		t.Fatalf("Error inserting transaction: %v", err)
	}
	list, err := transactionRepository.GetTransactionsList(userID, 0, 10)
	if err != nil {
		t.Fatalf("No error getting transaction: %v", err)
	}
	CheckTransactionsLength(t, list, 2)
	CheckTransaction(t, &list[0], list[0].ID, userID, tu.TestTransaction.Amount, tu.TestTransaction.Status, tu.TestTransaction.Type)
	CheckTransaction(t, &list[1], list[1].ID, userID, tu.TestTransaction.Amount, tu.TestTransaction.Status, tu.TestTransaction.Type)
}

func TestInsertPendingContractPaymentTransactionCorrect(t *testing.T) {
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
	err = setupTransactionTables(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := CreateSqlTransactionRepository(db, "transactions", "pending_contract_payment_transactions", "sequence")
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataID, err := personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	tu.TestUser.PersonalDataID = personalDataID
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	userID, err := userRepository.InsertUser(tu.TestUser)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	tu.TestTransaction.UserID = userID
	transactionID, err := transactionRepository.InsertTransaction(tu.TestTransaction)
	if err != nil {
		t.Fatalf("Error inserting transaction: %v", err)
	}
	tu.TestPendingContractPaymentTransaction.UserID = userID
	tu.TestPendingContractPaymentTransaction.TransactionID = transactionID
	tu.TestTransaction.ID = transactionID
	_, err = transactionRepository.InsertPendingContractPaymentTransaction(tu.TestPendingContractPaymentTransaction, tu.TestTransaction)
	if err != nil {
		t.Fatalf("Error inserting pending contract payment transaction: %v", err)
	}
}

func TestInsertPendingContractPaymentTransactionIncorrect(t *testing.T) {
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
	err = setupTransactionTables(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := CreateSqlTransactionRepository(db, "transactions", "pending_contract_payment_transactions", "sequence")
	_, err = transactionRepository.InsertPendingContractPaymentTransaction(tu.TestPendingContractPaymentTransaction, tu.TestTransaction)
	if err == nil {
		t.Fatalf("No error inserting pending contract payment transaction: %v", err)
	}
	_, err = transactionRepository.InsertPendingContractPaymentTransaction(tu.TestPendingContractPaymentTransaction, tu.TestTransaction)
	if err == nil {
		t.Fatalf("No error inserting pending contract payment transaction: %v", err)
	}
}

func CheckPendingContractPaymentTransaction(
	t *testing.T,
	PendingContractPaymentTransaction *types.DBPendingContractPaymentTransaction,
	PendingContractPaymentTransactionID int64,
	UserID int64,
	TransactionID int64,
	Amount int64,
) {
	if PendingContractPaymentTransaction.UserID != UserID {
		t.Fatalf("Pending contract payment transaction user id not updated: %v", PendingContractPaymentTransaction)
	}
	if PendingContractPaymentTransaction.TransactionID != TransactionID {
		t.Fatalf("Pending contract payment transaction transaction id not updated: %v", PendingContractPaymentTransaction)
	}
	if PendingContractPaymentTransaction.Amount != Amount {
		t.Fatalf("Pending contract payment transaction amount not updated: %v", PendingContractPaymentTransaction)
	}
}

func TestGetPendingContractPaymentTransactionCorrect(t *testing.T) {
	db := SetupTestDb(t)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err := setupTransactionTables(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := CreateSqlTransactionRepository(db, "transactions", "pending_contract_payment_transactions", "sequence")
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataID, err := personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	tu.TestUser.PersonalDataID = personalDataID
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	userID, err := userRepository.InsertUser(tu.TestUser)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	tu.TestTransaction.UserID = userID
	transactionID, err := transactionRepository.InsertTransaction(tu.TestTransaction)
	if err != nil {
		t.Fatalf("Error inserting transaction: %v", err)
	}
	tu.TestPendingContractPaymentTransaction.UserID = userID
	tu.TestPendingContractPaymentTransaction.TransactionID = transactionID
	tu.TestTransaction.ID = transactionID
	_, err = transactionRepository.InsertPendingContractPaymentTransaction(tu.TestPendingContractPaymentTransaction, tu.TestTransaction)
	if err != nil {
		t.Fatalf("Error inserting pending contract payment transaction: %v", err)
	}
	transaction, err := transactionRepository.GetPendingContractPaymentTransaction()
	if err != nil {
		t.Fatalf("Error getting pending contract payment transaction: %v", err)
	}
	CheckPendingContractPaymentTransaction(t, transaction, transactionID, userID, transactionID, tu.TestPendingContractPaymentTransaction.Amount)
}

func TestApproveTransactionCorrect(t *testing.T) {
	db := SetupTestDb(t)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err := setupTransactionTables(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := CreateSqlTransactionRepository(db, "transactions", "pending_contract_payment_transactions", "sequence")
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataID, err := personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	tu.TestUser.PersonalDataID = personalDataID
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	userID, err := userRepository.InsertUser(tu.TestUser)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	tu.TestTransaction.UserID = userID
	transactionID, err := transactionRepository.InsertTransaction(tu.TestTransaction)
	if err != nil {
		t.Fatalf("Error inserting transaction: %v", err)
	}
	tu.TestPendingContractPaymentTransaction.UserID = userID
	tu.TestPendingContractPaymentTransaction.TransactionID = transactionID
	tu.TestTransaction.ID = transactionID
	err = transactionRepository.ApproveTransaction(transactionID)
	if err != nil {
		t.Fatalf("Error approving transaction: %v", err)
	}
	transaction, err := transactionRepository.GetPendingContractPaymentTransaction()
	if err != nil {
		t.Fatalf("Error getting pending contract payment transaction: %v", err)
	}
	if transaction != nil {
		t.Fatalf("Transaction not approved: %v", transaction)
	}
}

func TestChangeTransactionStatusCorrect(t *testing.T) {
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
	err = setupTransactionTables(db)
	if err != nil {
		t.Fatalf("Error setting up transaction tables: %v", err)
	}
	transactionRepository := CreateSqlTransactionRepository(db, "transactions", "pending_contract_payment_transactions", "sequence")
	personalDataRepository := CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	personalDataID, err := personalDataRepository.InsertPersonalData(tu.TestPD)
	if err != nil {
		t.Fatalf("Error inserting personal data: %v", err)
	}
	tu.TestUser.PersonalDataID = personalDataID
	userRepository := CreateSqlUserRepository(db, "users", "sequence")
	userID, err := userRepository.InsertUser(tu.TestUser)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	tu.TestTransaction.UserID = userID
	transactionID, err := transactionRepository.InsertTransaction(tu.TestTransaction)
	if err != nil {
		t.Fatalf("Error inserting transaction: %v", err)
	}
	err = transactionRepository.UpdateTransactionStatus(transactionID, types.TransactionStatusPaid)
	if err != nil {
		t.Fatalf("Error changing transaction status: %v", err)
	}
}
