package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type ITransactionRepository interface {
	InsertTransaction(transaction types.DBTransaction) (int64, error)
	UpdateTransactionStatus(transactionId int64, status types.TransactionStatus) error
	GetTransaction(transactionId int64) (*types.DBTransaction, error)
	GetTransactionsList(userId int64, from int64, size int64) ([]types.DBTransaction, error)
	InsertPendingContractPaymentTransaction(
		transactionPendingContractPayment types.DBPendingContractPaymentTransaction,
		transaction types.DBTransaction,
	) (int64, error)
	GetPendingContractPaymentTransaction() (*types.DBPendingContractPaymentTransaction, error)
	ApproveTransaction(transactionId int64) error
}

func CreateSqlTransactionTable(db *sql.DB, transactionTableName string, userTableName string, pendingContractPaymentTransactionsTableName string, sequenceTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + transactionTableName + ` (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		status INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		type INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", transactionTableName, err)
	}

	return nil
}

func CreateSqlPendingContractPaymentTransactionsTable(
	db *sql.DB,
	pendingContractPaymentTransactionsTableName string,
	userTableName string,
	transactionTableName string,
	sequenceTableName string,
) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + pendingContractPaymentTransactionsTableName + ` (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		transaction_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES ` + userTableName + `(id),
		FOREIGN KEY (transaction_id) REFERENCES ` + transactionTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", pendingContractPaymentTransactionsTableName, err)
	}
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating delete trigger: %v", err)
	}

	return nil
}

type SqlTransactionRepository struct {
	db                                      *sql.DB
	transactionTable                        string
	pendingContractPaymentTransactionsTable string
	sequenceName                            string
}

func CreateSqlTransactionRepository(db *sql.DB, transactionTable string, pendingContractPaymentTransactionsTable string, sequenceName string) *SqlTransactionRepository {
	return &SqlTransactionRepository{
		db:                                      db,
		transactionTable:                        transactionTable,
		pendingContractPaymentTransactionsTable: pendingContractPaymentTransactionsTable,
		sequenceName:                            sequenceName,
	}
}

func (r *SqlTransactionRepository) InsertTransaction(transaction types.DBTransaction) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.transactionTable + ` (id, user_id, amount, status, created_at, type)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = r.db.Exec(query, id, transaction.UserID, transaction.Amount, transaction.Status, transaction.CreatedAt, transaction.Type)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlTransactionRepository) UpdateTransactionStatus(transactionId int64, status types.TransactionStatus) error {
	query := `
	UPDATE ` + r.transactionTable + ` SET status = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, status, transactionId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}
	return nil
}

func (r *SqlTransactionRepository) GetTransaction(id int64) (*types.DBTransaction, error) {
	query := `
	SELECT * FROM ` + r.transactionTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var transaction types.DBTransaction
	err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Status, &transaction.CreatedAt, &transaction.Type)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *SqlTransactionRepository) GetTransactionsList(userId int64, from int64, size int64) ([]types.DBTransaction, error) {
	query := `
	SELECT * FROM ` + r.transactionTable + ` WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userId, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []types.DBTransaction
	for rows.Next() {
		var transaction types.DBTransaction
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Status, &transaction.CreatedAt, &transaction.Type)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *SqlTransactionRepository) InsertPendingContractPaymentTransaction(
  transactionPendingContractPayment types.DBPendingContractPaymentTransaction,
  transaction types.DBTransaction,
) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `INSERT INTO ` + r.pendingContractPaymentTransactionsTable + ` (id, user_id, amount, created_at, transaction_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	_, err = r.db.Exec(query, id, transactionPendingContractPayment.UserID, transactionPendingContractPayment.Amount, transactionPendingContractPayment.CreatedAt, transaction.ID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlTransactionRepository) GetPendingContractPaymentTransaction() (*types.DBPendingContractPaymentTransaction, error) {
	query := `
	SELECT * FROM ` + r.pendingContractPaymentTransactionsTable + ` LIMIT 1
	`
	row := r.db.QueryRow(query)
	var transaction types.DBPendingContractPaymentTransaction
	err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.CreatedAt, &transaction.TransactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *SqlTransactionRepository) ApproveTransaction(transactionId int64) error {
	query := `
	DELETE FROM ` + r.pendingContractPaymentTransactionsTable + ` WHERE id = $1
	`
	_, err := r.db.Exec(query, transactionId)
	if err != nil {
		return err
	}
	return nil
}
