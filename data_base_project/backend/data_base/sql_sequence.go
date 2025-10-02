package data_base

import (
	"database/sql"
)

func CreateSqlSequence(db *sql.DB, sequenceName string) error {
	query := `
    CREATE SEQUENCE IF NOT EXISTS ` + sequenceName + ` START 1
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
