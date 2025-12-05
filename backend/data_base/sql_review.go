package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IReviewRepository interface {
	InsertReview(review types.DBReview) (int64, error)
	InsertReviewInSeq(tx *sql.Tx, review types.DBReview) (int64, error)
	UpdateReview(review types.DBReview) error
	GetReview(id int64) (*types.DBReview, error)
	GetReviewsByRepetitorID(repetitorID int64, from int64, size int64) ([]types.DBReview, error)
	GetReviewsByClientID(clientID int64, from int64, size int64) ([]types.DBReview, error)
}

func CreateSqlReviewTable(db *sql.DB, reviewTableName string, userTableName string, sequenceName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + reviewTableName + ` (
		id INTEGER PRIMARY KEY,
        contract_id INTEGER NOT NULL,
		client_id INTEGER NOT NULL,
		repetitor_id INTEGER NOT NULL,
		rating INTEGER NOT NULL,
		comment TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", reviewTableName, err)
	}
	return nil
}

type SqlReviewRepository struct {
	db           *sql.DB
	reviewTable  string
	sequenceName string
}

func CreateSqlReviewRepository(db *sql.DB, reviewTable string, sequenceName string) *SqlReviewRepository {
	return &SqlReviewRepository{
		db:           db,
		reviewTable:  reviewTable,
		sequenceName: sequenceName,
	}
}

func (r *SqlReviewRepository) InsertReview(review types.DBReview) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
    INSERT INTO ` + r.reviewTable + ` (id, contract_id, client_id, repetitor_id, rating, comment, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err = r.db.Exec(query, id, review.ContractID, review.ClientID, review.RepetitorID, review.Rating, review.Comment, review.CreatedAt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlReviewRepository) InsertReviewInSeq(tx *sql.Tx, review types.DBReview) (int64, error) {
	var id int64
	err := tx.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
    INSERT INTO ` + r.reviewTable + ` (id, contract_id, client_id, repetitor_id, rating, comment, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err = tx.Exec(query, id, review.ContractID, review.ClientID, review.RepetitorID, review.Rating, review.Comment, review.CreatedAt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlReviewRepository) GetReview(id int64) (*types.DBReview, error) {
	query := `
	SELECT * FROM ` + r.reviewTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var review types.DBReview
	err := row.Scan(&review.ID, &review.ContractID, &review.ClientID, &review.RepetitorID, &review.Rating, &review.Comment, &review.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *SqlReviewRepository) UpdateReview(review types.DBReview) error {
	query := `
	UPDATE ` + r.reviewTable + ` SET rating = $1, comment = $2 WHERE id = $3
	`
	result, err := r.db.Exec(query, review.Rating, review.Comment, review.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("review not found")
	}
	return nil
}

func (r *SqlReviewRepository) GetReviewsByRepetitorID(repetitorID int64, from int64, size int64) ([]types.DBReview, error) {
	query := `
    SELECT * FROM ` + r.reviewTable + ` WHERE repetitor_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
    `
	rows, err := r.db.Query(query, repetitorID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reviews := []types.DBReview{}
	for rows.Next() {
		var review types.DBReview
		err := rows.Scan(&review.ID, &review.ContractID, &review.ClientID, &review.RepetitorID, &review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *SqlReviewRepository) GetReviewsByClientID(clientID int64, from int64, size int64) ([]types.DBReview, error) {
	query := `
    SELECT * FROM ` + r.reviewTable + ` WHERE client_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
    `
	rows, err := r.db.Query(query, clientID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reviews := []types.DBReview{}
	for rows.Next() {
		var review types.DBReview
		err := rows.Scan(&review.ID, &review.ContractID, &review.ClientID, &review.RepetitorID, &review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reviews, nil
}
