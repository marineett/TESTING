package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IPersonalDataRepository interface {
	InsertPersonalDataInSeq(tx *sql.Tx, personalData types.DBPersonalData) (int64, error)
	InsertPersonalData(personalData types.DBPersonalData) (int64, error)
	GetPersonalData(id int64) (*types.DBPersonalData, error)
	UpdatePersonalData(id int64, personalData types.DBPersonalData) error
}

func CreateSqlPersonalDataTable(db *sql.DB, personalDataTableName string, sequenceName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + personalDataTableName + ` (
		id INTEGER PRIMARY KEY,
		telephone_number VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		passport_number VARCHAR(255) NOT NULL,
		passport_series VARCHAR(255) NOT NULL,
		passport_date TIMESTAMP NOT NULL,
		passport_issued_by VARCHAR(255) NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		middle_name VARCHAR(255) NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", personalDataTableName, err)
	}

	return nil
}

type SqlPersonalDataRepository struct {
	db                *sql.DB
	personalDataTable string
	sequenceName      string
}

func CreateSqlPersonalDataRepository(db *sql.DB, personalDataTable string, sequenceName string) *SqlPersonalDataRepository {
	return &SqlPersonalDataRepository{
		db:                db,
		personalDataTable: personalDataTable,
		sequenceName:      sequenceName,
	}
}

func (r *SqlPersonalDataRepository) InsertPersonalDataInSeq(tx *sql.Tx, personalData types.DBPersonalData) (int64, error) {
	var id int64
	err := tx.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO ` + r.personalDataTable + ` (
	id, 
	telephone_number, 
	email, 
	passport_number, 
	passport_series, 
	passport_date, 
	passport_issued_by, 
	first_name, 
	last_name, 
	middle_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = tx.Exec(query, id, personalData.TelephoneNumber, personalData.Email, personalData.PassportNumber, personalData.PassportSeries, personalData.PassportDate, personalData.PassportIssuedBy, personalData.FirstName, personalData.LastName, personalData.MiddleName)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SqlPersonalDataRepository) InsertPersonalData(personalData types.DBPersonalData) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('sequence')").Scan(&id)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO ` + r.personalDataTable + ` (
	id, 
	telephone_number, 
	email, 
	passport_number, 
	passport_series, 
	passport_date, 
	passport_issued_by, 
	first_name, 
	last_name, 
	middle_name
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = r.db.Exec(query, id, personalData.TelephoneNumber, personalData.Email, personalData.PassportNumber, personalData.PassportSeries, personalData.PassportDate, personalData.PassportIssuedBy, personalData.FirstName, personalData.LastName, personalData.MiddleName)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SqlPersonalDataRepository) GetPersonalData(id int64) (*types.DBPersonalData, error) {
	query := `
	SELECT id, 
	telephone_number, 
	email, 
	passport_number, 
	passport_series, 
	passport_date, 
	passport_issued_by, 
	first_name, 
	last_name, 
	middle_name  
	FROM ` + r.personalDataTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var personalData types.DBPersonalData
	err := row.Scan(&personalData.ID, &personalData.TelephoneNumber, &personalData.Email, &personalData.PassportNumber, &personalData.PassportSeries, &personalData.PassportDate, &personalData.PassportIssuedBy, &personalData.FirstName, &personalData.LastName, &personalData.MiddleName)
	if err != nil {
		return nil, err
	}
	return &personalData, nil
}

func (r *SqlPersonalDataRepository) UpdatePersonalData(id int64, personalData types.DBPersonalData) error {
	query := `UPDATE ` +
		r.personalDataTable +
		` SET 
	telephone_number = $1, 
	email = $2, 
	passport_number = $3, 
	passport_series = $4, 
	passport_date = $5, 
	passport_issued_by = $6, 
	first_name = $7, 
	last_name = $8, 
	middle_name = $9
	WHERE id = $10`
	result, err := r.db.Exec(
		query,
		personalData.TelephoneNumber,
		personalData.Email,
		personalData.PassportNumber,
		personalData.PassportSeries,
		personalData.PassportDate,
		personalData.PassportIssuedBy,
		personalData.FirstName,
		personalData.LastName,
		personalData.MiddleName,
		id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}
