package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IUserRepository interface {
	InsertUser(user types.DBUserData) (int64, error)
	InsertUserInSeq(tx *sql.Tx, user types.DBUserData) (int64, error)
	GetUser(id int64) (*types.DBUserData, error)
}

func CreateSqlUserTable(db *sql.DB, userTableName string, personalDataTable string, sequenceName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + userTableName + ` (
		id INTEGER PRIMARY KEY,
		registration_date TIMESTAMP NOT NULL,
		last_login_date TIMESTAMP NOT NULL,
		personal_data_id INTEGER NOT NULL,
		FOREIGN KEY (personal_data_id) REFERENCES ` + personalDataTable + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", userTableName, err)
	}

	return nil
}

type SqlUserRepository struct {
	db           *sql.DB
	userTable    string
	sequenceName string
}

func CreateSqlUserRepository(db *sql.DB, userTable string, sequenceName string) *SqlUserRepository {
	return &SqlUserRepository{
		db:           db,
		userTable:    userTable,
		sequenceName: sequenceName,
	}
}

func (r *SqlUserRepository) InsertUser(user types.DBUserData) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.userTable + ` (id, registration_date, last_login_date, personal_data_id)
	VALUES ($1, $2, $3, $4)
	`
	_, err = r.db.Exec(query, id, user.RegistrationDate, user.LastLoginDate, user.PersonalDataID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlUserRepository) InsertUserInSeq(tx *sql.Tx, user types.DBUserData) (int64, error) {
	var id int64
	err := tx.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.userTable + ` (id, registration_date, last_login_date, personal_data_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	_, err = tx.Exec(query, id, user.RegistrationDate, user.LastLoginDate, user.PersonalDataID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlUserRepository) GetUser(id int64) (*types.DBUserData, error) {
	query := `
	SELECT id, registration_date, last_login_date, personal_data_id FROM ` + r.userTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var user types.DBUserData
	err := row.Scan(&user.ID, &user.RegistrationDate, &user.LastLoginDate, &user.PersonalDataID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
