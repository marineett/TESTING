package data_base

import (
	"data_base_project/types"
	"database/sql"
)

type ILessonRepository interface {
	InsertLesson(lesson types.DBLesson) (int64, error)
	GetLessons(contractID int64, from int64, size int64) ([]types.DBLesson, error)
	GetLesson(lessonID int64) (*types.DBLesson, error)
}

func CreateSqlLessonTable(db *sql.DB, lessonTable string, contractTable string, transactionTable string) error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + lessonTable + ` (
		id INTEGER PRIMARY KEY,
		contract_id INTEGER NOT NULL,
		duration INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (contract_id) REFERENCES ` + contractTable + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

type SqlLessonRepository struct {
	db               *sql.DB
	lessonTable      string
	contractTable    string
	transactionTable string
	sequenceName     string
}

func CreateSqlLessonRepository(db *sql.DB, lessonTable string, contractTable string, transactionTable string, sequenceName string) *SqlLessonRepository {
	return &SqlLessonRepository{
		db:               db,
		lessonTable:      lessonTable,
		contractTable:    contractTable,
		transactionTable: transactionTable,
		sequenceName:     sequenceName,
	}
}

func (r *SqlLessonRepository) InsertLesson(lesson types.DBLesson) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
		INSERT INTO ` + r.lessonTable + ` (id,contract_id, duration, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = r.db.Exec(query, id, lesson.ContractID, lesson.Duration, lesson.CreatedAt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlLessonRepository) GetLessons(contractID int64, from int64, size int64) ([]types.DBLesson, error) {
	query := `
		SELECT id, contract_id, duration, created_at
		FROM ` + r.lessonTable + `
		WHERE contract_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, contractID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lessons := []types.DBLesson{}
	for rows.Next() {
		var lesson types.DBLesson
		err := rows.Scan(&lesson.ID, &lesson.ContractID, &lesson.Duration, &lesson.CreatedAt)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}

func (r *SqlLessonRepository) GetLesson(lessonID int64) (*types.DBLesson, error) {
	query := `
		SELECT id, contract_id, duration, created_at
		FROM ` + r.lessonTable + `
		WHERE id = $1
	`
	var lesson types.DBLesson
	err := r.db.QueryRow(query, lessonID).Scan(&lesson.ID, &lesson.ContractID, &lesson.Duration, &lesson.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}
