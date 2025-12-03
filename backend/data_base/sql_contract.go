package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
	"time"
)

type IContractRepository interface {
	InsertContract(contract types.DBContract) (int64, error)
	BeginTx() (*sql.Tx, error)
	GetContract(id int64) (*types.DBContract, error)
	GetContractsByRepetitorID(repetitorID int64, from int64, size int64, status types.ContractStatus) ([]types.DBContract, error)
	GetContractsByClientID(clientID int64, from int64, size int64, status types.ContractStatus) ([]types.DBContract, error)
	UpdateContractStatus(id int64, status types.ContractStatus) error
	UpdateContractRepetitorID(contractID int64, repetitorID int64) error
	UpdateContractPaymentStatus(id int64, paymentStatus types.PaymentStatus) error
	UpdateContractReviewClientID(id int64, reviewClientID int64) error
	UpdateContractReviewClientIDInSeq(tx *sql.Tx, id int64, reviewClientID int64) error
	UpdateContractReviewRepetitorID(id int64, reviewRepetitorID int64) error
	GetContractList(from int64, size int64, status types.ContractStatus) ([]types.DBContract, error)
	GetAllContracts(from int64, size int64) ([]types.DBContract, error)
	GetContracts(clientID int64, repetitorID int64, from int64, size int64) ([]types.DBContract, error)
}

func CreateSqlContractTable(
	db *sql.DB,
	contractTableName string,
	userTableName string,
	reviewTableName string,
	repetitorTableName string,
	clientTableName string,
) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + contractTableName + ` (
		id INTEGER PRIMARY KEY,
		client_id INTEGER NOT NULL,
		repetitor_id INTEGER NOT NULL,
		review_client_id INTEGER NOT NULL,
		review_repetitor_id INTEGER NOT NULL,
		transaction_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		description TEXT NOT NULL,
		status INTEGER NOT NULL,
		payment_status INTEGER NOT NULL,

		price INTEGER NOT NULL,
		commission INTEGER NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,

		FOREIGN KEY (client_id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", contractTableName, err)
	}
	return nil
}

type SqlContractRepository struct {
	db            *sql.DB
	contractTable string
	sequenceName  string
}

func CreateSqlContractRepository(db *sql.DB, contractTable string, sequenceName string) *SqlContractRepository {
	return &SqlContractRepository{
		db:            db,
		contractTable: contractTable,
		sequenceName:  sequenceName,
	}
}

func (r *SqlContractRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *SqlContractRepository) InsertContract(contract types.DBContract) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.contractTable + ` (
	id,
	client_id, 
	repetitor_id, 
	description, 
	status, 
	transaction_id, 
	created_at, 
	payment_status, 
	review_client_id, 
	review_repetitor_id, 
	price, 
	commission, 
	start_date, 
	end_date) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err = r.db.Exec(
		query,
		id,
		contract.ClientID,
		contract.RepetitorID,
		contract.Description,
		contract.Status,
		contract.TransactionID,
		contract.CreatedAt,
		contract.PaymentStatus,
		contract.ReviewClientID,
		contract.ReviewRepetitorID,
		contract.Price,
		contract.Commission,
		contract.StartDate,
		contract.EndDate)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlContractRepository) GetContract(id int64) (*types.DBContract, error) {
	query := `
	SELECT 
	id, 
	client_id, 
	repetitor_id, 
	created_at, 
	description, 
	status, 
	payment_status, 
	review_client_id, 
	review_repetitor_id, 
	price, 
	commission, 
	start_date, 
	end_date, 
	transaction_id
	FROM ` + r.contractTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var contract types.DBContract
	err := row.Scan(
		&contract.ID,
		&contract.ClientID,
		&contract.RepetitorID,
		&contract.CreatedAt,
		&contract.Description,
		&contract.Status,
		&contract.PaymentStatus,
		&contract.ReviewClientID,
		&contract.ReviewRepetitorID,
		&contract.Price,
		&contract.Commission,
		&contract.StartDate,
		&contract.EndDate,
		&contract.TransactionID,
	)
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *SqlContractRepository) GetContractsByRepetitorID(repetitorID int64, from int64, size int64, status types.ContractStatus) ([]types.DBContract, error) {
	query := `
	SELECT 
		id, 
		client_id, 
		repetitor_id, 
		created_at, 
		description, 
		status, 
		payment_status, 
		review_client_id, 
		review_repetitor_id, 
		price, 
		commission, 
		start_date, 
		end_date FROM ` + r.contractTable + ` WHERE 
		repetitor_id = $1 AND status = $2 
		ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`
	rows, err := r.db.Query(query, repetitorID, status, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contracts := []types.DBContract{}
	for rows.Next() {
		var contract types.DBContract
		err := rows.Scan(
			&contract.ID,
			&contract.ClientID,
			&contract.RepetitorID,
			&contract.CreatedAt,
			&contract.Description,
			&contract.Status,
			&contract.PaymentStatus,
			&contract.ReviewClientID,
			&contract.ReviewRepetitorID,
			&contract.Price,
			&contract.Commission,
			&contract.StartDate,
			&contract.EndDate)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *SqlContractRepository) GetContractsByClientID(clientID int64, from int64, size int64, status types.ContractStatus) ([]types.DBContract, error) {
	fmt.Println(r.contractTable)
	query := `
		SELECT id, 
		client_id, 
		repetitor_id, 
		created_at, 
		description, 
		status, 
		payment_status, 
		review_client_id, 
		review_repetitor_id, 
		price, 
		commission, 
		start_date, 
		end_date,
		transaction_id
		FROM ` + r.contractTable + ` WHERE client_id = $1 AND status = $2 
		ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`
	rows, err := r.db.Query(query, clientID, status, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contracts := []types.DBContract{}
	for rows.Next() {
		var contract types.DBContract
		err := rows.Scan(
			&contract.ID,
			&contract.ClientID,
			&contract.RepetitorID,
			&contract.CreatedAt,
			&contract.Description,
			&contract.Status,
			&contract.PaymentStatus,
			&contract.ReviewClientID,
			&contract.ReviewRepetitorID,
			&contract.Price,
			&contract.Commission,
			&contract.StartDate,
			&contract.EndDate,
			&contract.TransactionID,
		)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *SqlContractRepository) GetContractList(from int64, size int64, status types.ContractStatus) ([]types.DBContract, error) {
	query := `
	SELECT 
	id, 
	client_id, 
	repetitor_id,
	created_at, 
	description, 
	status, 
	payment_status, 
	review_client_id, 
	review_repetitor_id, 
	price, 
	commission, 
	start_date, 
	end_date
	FROM ` + r.contractTable + ` 
	WHERE status = $1 
	ORDER BY created_at DESC 
	LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, status, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contracts := []types.DBContract{}
	for rows.Next() {
		var contract types.DBContract
		err := rows.Scan(
			&contract.ID,
			&contract.ClientID,
			&contract.RepetitorID,
			&contract.CreatedAt,
			&contract.Description,
			&contract.Status,
			&contract.PaymentStatus,
			&contract.ReviewClientID,
			&contract.ReviewRepetitorID,
			&contract.Price,
			&contract.Commission,
			&contract.StartDate,
			&contract.EndDate,
		)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *SqlContractRepository) GetAllContracts(from int64, size int64) ([]types.DBContract, error) {
	query := `
	SELECT 
	id, 
	client_id, 
	repetitor_id, 
	created_at, 
	description, 
	status, 
	payment_status, 
	review_client_id, 
	review_repetitor_id, 
	price, 
	commission, 
	start_date, 
	end_date 
	FROM ` + r.contractTable + ` 
	ORDER BY created_at DESC 
	LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contracts := []types.DBContract{}
	for rows.Next() {
		var contract types.DBContract
		err := rows.Scan(
			&contract.ID,
			&contract.ClientID,
			&contract.RepetitorID,
			&contract.CreatedAt,
			&contract.Description,
			&contract.Status,
			&contract.PaymentStatus,
			&contract.ReviewClientID,
			&contract.ReviewRepetitorID,
			&contract.Price,
			&contract.Commission,
			&contract.StartDate,
			&contract.EndDate,
		)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *SqlContractRepository) UpdateContractStatus(id int64, status types.ContractStatus) error {
	query := `
	UPDATE ` + r.contractTable + ` SET status = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *SqlContractRepository) UpdateContractPaymentStatus(id int64, paymentStatus types.PaymentStatus) error {
	query := `
	UPDATE ` + r.contractTable + ` SET payment_status = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, paymentStatus, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *SqlContractRepository) UpdateContractReviewClientID(id int64, reviewClientID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	UPDATE ` + r.contractTable + ` SET review_client_id = $1 WHERE id = $2
	`
	result, err := tx.Exec(query, reviewClientID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}

	return tx.Commit()
}

func (r *SqlContractRepository) UpdateContractReviewRepetitorID(id int64, reviewRepetitorID int64) error {
	query := `
	UPDATE ` + r.contractTable + ` SET review_repetitor_id = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, reviewRepetitorID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *SqlContractRepository) UpdateContractPrice(id int64, price int64) error {
	query := `
	UPDATE ` + r.contractTable + ` SET price = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, price, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *SqlContractRepository) UpdateContractCommission(id int64, commission int64) error {
	query := `
	UPDATE ` + r.contractTable + ` SET commission = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, commission, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *SqlContractRepository) UpdateContractStartDate(id int64, startDate time.Time) error {
	query := `
	UPDATE ` + r.contractTable + ` SET start_date = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, startDate, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *SqlContractRepository) UpdateContractRepetitorID(contractID int64, repetitorID int64) error {

	query := `
	UPDATE ` + r.contractTable + ` SET repetitor_id = $1, status = $2 WHERE id = $3
	`

	result, err := r.db.Exec(query, repetitorID, types.ContractStatusActive, contractID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *SqlContractRepository) UpdateContractReviewClientIDInSeq(tx *sql.Tx, id int64, reviewClientID int64) error {
	query := `
	UPDATE ` + r.contractTable + ` SET review_client_id = $1 WHERE id = $2
	`
	result, err := tx.Exec(query, reviewClientID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

func (r *SqlContractRepository) GetContracts(clientID int64, repetitorID int64, from int64, size int64) ([]types.DBContract, error) {
	query := `
	SELECT id, client_id, repetitor_id, created_at, description, status, payment_status, review_client_id, review_repetitor_id, price, commission, start_date, end_date FROM ` + r.contractTable + ` WHERE client_id = $1 AND repetitor_id = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4
	`
	rows, err := r.db.Query(query, clientID, repetitorID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	contracts := []types.DBContract{}
	for rows.Next() {
		var contract types.DBContract
		err := rows.Scan(
			&contract.ID,
			&contract.ClientID,
			&contract.RepetitorID,
			&contract.CreatedAt,
			&contract.Description,
			&contract.Status,
			&contract.PaymentStatus,
			&contract.ReviewClientID,
			&contract.ReviewRepetitorID,
			&contract.Price,
			&contract.Commission,
			&contract.StartDate,
			&contract.EndDate,
		)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return contracts, nil
}
