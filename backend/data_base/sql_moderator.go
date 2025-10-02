package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
	"time"
)

type IModeratorRepository interface {
	InsertModerator(moderator types.DBModeratorData, personalData types.DBPersonalData, auth types.DBAuthData) (int64, error)
	GetModerator(moderatorID int64) (*types.DBModeratorData, error)
	UpdateModeratorPersonalData(moderatorID int64, personalData types.DBPersonalData) error
	UpdateModeratorPassword(moderatorID int64, authData types.DBAuthData, newPassword string) error
	UpdateModeratorSalary(moderatorID int64, salary int64) error
	GetModerators() ([]int64, error)
}

func CreateSqlModeratorTable(db *sql.DB, moderatorTableName string, userTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + moderatorTableName + ` (
		id INTEGER PRIMARY KEY,
		salary INTEGER NOT NULL,
		FOREIGN KEY (id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", moderatorTableName, err)
	}
	return nil
}

type SqlModeratorRepository struct {
	db                     *sql.DB
	moderatorTable         string
	userRepository         IUserRepository
	personalDataRepository IPersonalDataRepository
	authRepository         IAuthRepository
	sequenceName           string
}

func CreateSqlModeratorRepository(db *sql.DB, personalDataTable string, userTable string, moderatorTable string, authTable string, sequenceName string) *SqlModeratorRepository {
	return &SqlModeratorRepository{
		db:                     db,
		moderatorTable:         moderatorTable,
		userRepository:         CreateSqlUserRepository(db, userTable, sequenceName),
		personalDataRepository: CreateSqlPersonalDataRepository(db, personalDataTable, sequenceName),
		authRepository:         CreateSqlAuthRepository(db, authTable, sequenceName),
		sequenceName:           sequenceName,
	}
}

func (r *SqlModeratorRepository) InsertModerator(moderator types.DBModeratorData, personalData types.DBPersonalData, auth types.DBAuthData) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	personalData.ID, err = r.personalDataRepository.InsertPersonalDataInSeq(tx, personalData)
	if err != nil {
		return 0, err
	}
	userID, err := r.userRepository.InsertUserInSeq(tx, types.DBUserData{
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
		PersonalDataID:   personalData.ID,
	})
	if err != nil {
		return 0, err
	}
	_, err = r.authRepository.InsertAuthInSeq(tx, types.DBAuthInfo{
		UserID:   userID,
		UserType: types.Moderator,
		Login:    auth.Login,
		Password: auth.Password,
	})
	if err != nil {
		return 0, err
	}

	query := `
	INSERT INTO ` + r.moderatorTable + ` (id, salary)
	VALUES ($1, $2)
	`
	_, err = tx.Exec(query, userID, moderator.Salary)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *SqlModeratorRepository) GetModerator(moderatorID int64) (*types.DBModeratorData, error) {
	query := `
	SELECT * FROM ` + r.moderatorTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, moderatorID)
	var moderator types.DBModeratorData
	err := row.Scan(&moderator.ID, &moderator.Salary)
	if err != nil {
		return nil, err
	}
	return &moderator, nil
}

func (r *SqlModeratorRepository) UpdateModeratorPersonalData(moderator_id int64, personalData types.DBPersonalData) error {
	userData, err := r.userRepository.GetUser(moderator_id)
	if err != nil {
		return err
	}
	personalDataID := userData.PersonalDataID
	return r.personalDataRepository.UpdatePersonalData(personalDataID, personalData)
}

func (r *SqlModeratorRepository) UpdateModeratorPassword(moderator_id int64, authData types.DBAuthData, newPassword string) error {
	return r.authRepository.ChangePassword(moderator_id, authData, newPassword)
}

func (r *SqlModeratorRepository) UpdateModeratorSalary(moderator_id int64, salary int64) error {
	query := `
	UPDATE ` + r.moderatorTable + ` SET salary = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, salary, moderator_id)
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

func (r *SqlModeratorRepository) GetModerators() ([]int64, error) {
	query := `
	SELECT id FROM ` + r.moderatorTable + `
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
