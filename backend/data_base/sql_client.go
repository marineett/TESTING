package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
	"time"
)

type IClientRepository interface {
	InsertClient(client types.DBClientData, personalData types.DBPersonalData, authData types.DBAuthData) (int64, error)
	GetClient(id int64) (*types.DBClientData, error)
	UpdateClientPersonalData(clientID int64, personalData types.DBPersonalData) error
	UpdateClientPassword(clientID int64, authData types.DBAuthData, newPassword string) error
}

func CreateSqlClientTable(db *sql.DB, clientTableName string, userTableName string, sequenceName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + clientTableName + ` (
		id INTEGER PRIMARY KEY,
		summary_rating INTEGER NOT NULL,
		reviews_count INTEGER NOT NULL,
		FOREIGN KEY (id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", clientTableName, err)
	}
	return nil
}

type SqlClientRepository struct {
	db                     *sql.DB
	clientTable            string
	userRepository         IUserRepository
	personalDataRepository IPersonalDataRepository
	authRepository         IAuthRepository
}

func CreateSqlClientRepository(db *sql.DB, personalDataTable string, userTable string, clientTable string, authTable string, sequenceName string) *SqlClientRepository {
	return &SqlClientRepository{
		db:                     db,
		clientTable:            clientTable,
		userRepository:         CreateSqlUserRepository(db, userTable, sequenceName),
		personalDataRepository: CreateSqlPersonalDataRepository(db, personalDataTable, sequenceName),
		authRepository:         CreateSqlAuthRepository(db, authTable, sequenceName),
	}
}

func (r *SqlClientRepository) InsertClient(client types.DBClientData, personalData types.DBPersonalData, authData types.DBAuthData) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	personalDataID, err := r.personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		return 0, err
	}
	userID, err := r.userRepository.InsertUserInSeq(tx, types.DBUserData{
		ID:               client.ID,
		PersonalDataID:   personalDataID,
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
	})
	if err != nil {
		return 0, err
	}
	authInfo := types.DBAuthInfo{
		UserID:   userID,
		UserType: types.Client,
		Login:    authData.Login,
		Password: authData.Password,
	}
	_, err = r.authRepository.InsertAuthInSeq(tx, authInfo)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.clientTable + ` (id, summary_rating, reviews_count)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	row := tx.QueryRow(query, userID, client.SummaryRating, client.ReviewsCount)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlClientRepository) GetClient(id int64) (*types.DBClientData, error) {
	query := `
	SELECT * FROM ` + r.clientTable + ` WHERE id = $1
	`
	summaryRating, reviewsCount := 0, 0
	row := r.db.QueryRow(query, id)
	var client types.DBClientData
	err := row.Scan(&client.ID, &summaryRating, &reviewsCount)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *SqlClientRepository) UpdateClientPersonalData(client_id int64, personalData types.DBPersonalData) error {
	userData, err := r.userRepository.GetUser(client_id)
	if err != nil {
		return err
	}
	return r.personalDataRepository.UpdatePersonalData(userData.PersonalDataID, personalData)
}

func (r *SqlClientRepository) UpdateClientPassword(client_id int64, authData types.DBAuthData, newPassword string) error {
	return r.authRepository.ChangePassword(client_id, authData, newPassword)
}
