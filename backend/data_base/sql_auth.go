package data_base

import (
	"data_base_project/types"
	"database/sql"
	"errors"
	"fmt"
)

type IAuthRepository interface {
	InsertAuthInSeq(tx *sql.Tx, auth types.DBAuthInfo) (int64, error)
	InsertAuth(auth types.DBAuthInfo) (int64, error)
	ChangePassword(userId int64, authData types.DBAuthData, newPassword string) error
	Authorize(auth_data types.DBAuthData) (types.DBAuthVerdict, error)
	CheckLogin(login string) (bool, error)
}

func CreateSqlAuthTable(db *sql.DB, authTableName string, userTableName string, sequenceName string) error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + authTableName + ` (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		user_type INTEGER NOT NULL,
		login VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", authTableName, err)
	}
	return nil
}

type SqlAuthRepository struct {
	db           *sql.DB
	authTable    string
	sequenceName string
}

func CreateSqlAuthRepository(db *sql.DB, authTable string, sequenceName string) *SqlAuthRepository {
	return &SqlAuthRepository{
		db:           db,
		authTable:    authTable,
		sequenceName: sequenceName,
	}
}

func (r *SqlAuthRepository) InsertAuthInSeq(tx *sql.Tx, auth types.DBAuthInfo) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.authTable + ` (id, user_id, user_type, login, password) VALUES ($1, $2, $3, $4, $5)
	`
	_, err = r.db.Exec(query, id, auth.UserID, auth.UserType, auth.Login, auth.Password)
	if err != nil {
		return 0, err
	}
	return auth.UserID, nil
}

func (r *SqlAuthRepository) InsertAuth(auth types.DBAuthInfo) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.authTable + ` (id, user_id, user_type, login, password) VALUES ($1, $2, $3, $4, $5)
	`
	_, err = r.db.Exec(query, id, auth.UserID, auth.UserType, auth.Login, auth.Password)
	if err != nil {
		return 0, err
	}
	return auth.UserID, nil
}

func (r *SqlAuthRepository) ChangePassword(userId int64, authData types.DBAuthData, newPassword string) error {
	verdict, err := r.Authorize(authData)
	if err != nil {
		return err
	}
	if verdict.UserID != userId {
		return errors.New("invalid user id")
	}
	query := `
	UPDATE ` + r.authTable + ` SET password = $1 WHERE user_id = $2
	`
	_, err = r.db.Exec(query, newPassword, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *SqlAuthRepository) Authorize(auth_data types.DBAuthData) (types.DBAuthVerdict, error) {
	query := `
	SELECT * FROM ` + r.authTable + ` WHERE login = $1
	`
	var auth types.DBAuthInfo
	err := r.db.QueryRow(query, auth_data.Login).Scan(&auth.ID, &auth.UserID, &auth.UserType, &auth.Login, &auth.Password)
	if err != nil {
		return types.DBAuthVerdict{}, err
	}
	if auth.Password != auth_data.Password {
		return types.DBAuthVerdict{}, errors.New("invalid password")
	}
	return types.DBAuthVerdict{
		UserID:   auth.UserID,
		UserType: auth.UserType,
	}, nil
}

func (r *SqlAuthRepository) CheckLogin(login string) (bool, error) {
	query := `
	SELECT COUNT(*) FROM ` + r.authTable + ` WHERE login = $1
	`
	var count int
	err := r.db.QueryRow(query, login).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
