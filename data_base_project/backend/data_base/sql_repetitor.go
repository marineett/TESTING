package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
	"time"
)

type IRepetitorRepository interface {
	InsertRepetitor(repetitor types.DBRepetitorData, personalData types.DBPersonalData, auth types.DBAuthData) (int64, error)
	GetRepetitor(repetitorID int64) (*types.DBRepetitorData, error)
	UpdateRepetitorPersonalData(repetitorID int64, personalData types.DBPersonalData) error
	UpdateRepetitorPassword(repetitorID int64, authData types.DBAuthData, newPassword string) error
	GetRepetitorsIds(repetitorsOffset int64, repetitorsLimit int64) ([]int64, error)
}

func CreateSqlRepetitorTable(db *sql.DB, repetitorTableName string, userTableName string, resumeTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + repetitorTableName + ` (
		id INTEGER PRIMARY KEY,
		resume_id INTEGER NOT NULL,
		summary_rating INTEGER NOT NULL,
		reviews_count INTEGER NOT NULL,
		FOREIGN KEY (id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", repetitorTableName, err)
	}
	return nil
}

type SqlRepetitorRepository struct {
	db                     *sql.DB
	userRepository         IUserRepository
	authRepository         IAuthRepository
	personalDataRepository IPersonalDataRepository
	resumeRepository       IResumeRepository
	repetitorTable         string
	reviewTable            string
}

func CreateSqlRepetitorRepository(
	db *sql.DB,
	personalDataTable string,
	userTable string,
	repetitorTable string,
	authTable string,
	resumeTable string,
	reviewTable string,
	sequenceName string,
) *SqlRepetitorRepository {
	return &SqlRepetitorRepository{
		db:                     db,
		userRepository:         CreateSqlUserRepository(db, userTable, sequenceName),
		authRepository:         CreateSqlAuthRepository(db, authTable, sequenceName),
		personalDataRepository: CreateSqlPersonalDataRepository(db, personalDataTable, sequenceName),
		resumeRepository:       CreateSqlResumeRepository(db, resumeTable, sequenceName),
		repetitorTable:         repetitorTable,
		reviewTable:            reviewTable,
	}
}

func (r *SqlRepetitorRepository) InsertRepetitor(repetitor types.DBRepetitorData, personalData types.DBPersonalData, auth types.DBAuthData) (int64, error) {
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
		PersonalDataID:   personalData.ID,
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
	})
	if err != nil {
		return 0, err
	}
	_, err = r.authRepository.InsertAuthInSeq(tx, types.DBAuthInfo{
		UserID:   userID,
		UserType: types.Repetitor,
		Login:    auth.Login,
		Password: auth.Password,
	})
	if err != nil {
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.repetitorTable + ` (id, summary_rating, reviews_count, resume_id)
	VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(query, userID, 0, 0, 0)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *SqlRepetitorRepository) GetRepetitor(id int64) (*types.DBRepetitorData, error) {
	query := `
	SELECT id, resume_id, summary_rating, reviews_count FROM ` + r.repetitorTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var repetitor types.DBRepetitorData
	var reviewsCount int64
	err := row.Scan(&repetitor.ID, &repetitor.ResumeID, &repetitor.SummaryRating, &reviewsCount)
	if err != nil {
		return nil, err
	}
	return &repetitor, nil
}

func (r *SqlRepetitorRepository) UpdateRepetitorPersonalData(repetitor_id int64, personalData types.DBPersonalData) error {
	userData, err := r.userRepository.GetUser(repetitor_id)
	if err != nil {
		return err
	}
	personalDataID := userData.PersonalDataID
	return r.personalDataRepository.UpdatePersonalData(personalDataID, personalData)
}

func (r *SqlRepetitorRepository) UpdateRepetitorPassword(repetitor_id int64, authData types.DBAuthData, newPassword string) error {
	return r.authRepository.ChangePassword(repetitor_id, authData, newPassword)
}

func (r *SqlRepetitorRepository) GetRepetitorsIds(repetitorsOffset int64, repetitorsLimit int64) ([]int64, error) {
	query := `
	SELECT id FROM ` + r.repetitorTable + ` ORDER BY CASE WHEN reviews_count = 0 THEN 0 ELSE summary_rating / reviews_count END DESC OFFSET $1 LIMIT $2
	`
	rows, err := r.db.Query(query, repetitorsOffset, repetitorsLimit)
	if err != nil {
		return nil, err
	}
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
